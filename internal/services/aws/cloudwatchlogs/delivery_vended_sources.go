package cloudwatchlogs

import (
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The delivery-source family: the source record's cores, view and
// admission rules.

// --- Delivery source ---

// PutDeliverySourceInput holds the parsed parameters of PutDeliverySource.
type PutDeliverySourceInput struct {
	Name                        string
	ResourceArn                 string
	LogType                     string
	DeliverySourceConfiguration map[string]string
	Tags                        map[string]string
	Region                      string
}

func (s *LogsService) putDeliverySourceCore(input *PutDeliverySourceInput) (*DeliverySourceView, error) {
	if input.Name == "" {
		return nil, errValidationMember("name")
	}
	if len(input.Name) > logsstore.DeliveryNameMax || !deliveryNameRe.MatchString(input.Name) {
		return nil, errDeliveryValidation("invalid delivery source name %q: 1-%d characters, word characters and hyphens only",
			input.Name, logsstore.DeliveryNameMax)
	}
	if input.ResourceArn == "" {
		return nil, errValidationMember("resourceArn")
	}
	if input.LogType == "" {
		return nil, errValidationMember("logType")
	}
	if len(input.LogType) > logsstore.DeliveryLogTypeMax || !deliveryLogTypeRe.MatchString(input.LogType) {
		return nil, errDeliveryValidation("invalid logType %q: 1-%d characters, word characters only",
			input.LogType, logsstore.DeliveryLogTypeMax)
	}
	if err := validateDeliveryTags(input.Tags); err != nil {
		return nil, err
	}
	// "Both keys and values must be between 1 and 255 characters in
	// length." (DeliverySourceConfigurationKey/Value length traits).
	for key, value := range input.DeliverySourceConfiguration {
		if len(key) == 0 || len(key) > logsstore.DeliverySourceConfigurationMemberMax ||
			len(value) == 0 || len(value) > logsstore.DeliverySourceConfigurationMemberMax {
			return nil, errDeliveryValidation("deliverySourceConfiguration keys and values must be between 1 and %d characters in length",
				logsstore.DeliverySourceConfigurationMemberMax)
		}
	}

	parsed, err := svcarn.ParseARN(input.ResourceArn)
	if err != nil {
		return nil, errDeliveryValidation("resourceArn %q is not a valid ARN", input.ResourceArn)
	}
	// The source-service vocabulary: the platform's delivery sources are
	// CloudWatch Logs log groups — the one resource whose events the
	// delivery engine can read. AWS's source forms are other services'
	// resources (S3 buckets, VPCs, WAF web ACLs); none of those services
	// vends its logs into the platform's delivery pipeline, so a source
	// ARN of any other service would accept a delivery that can never
	// flow.
	if parsed.Service != "logs" || svcarn.ExtractLogGroupNameFromARN(input.ResourceArn) == "" {
		return nil, errDeliveryValidation("resourceArn %q does not address a deliverable source resource: this platform's delivery sources are CloudWatch Logs log groups", input.ResourceArn)
	}
	// The source ARN's region must address the store this call resolves:
	// the delivery engine reads the source group in the delivery record's
	// region — the region of the CreateDelivery that pairs the source,
	// which finds the source record there and nowhere else. An ARN naming
	// another region would deliver that region's same-named group under a
	// false identity, or admit a group the pipeline can never read.
	if parsed.Region != "" && parsed.Region != input.Region {
		return nil, errDeliveryValidation("resourceArn %q addresses region %s, but delivery sources are read in the call's region %s",
			input.ResourceArn, parsed.Region, input.Region)
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}
	groupName := resolveLogGroupIdentifier(input.ResourceArn)
	if _, err := store.GetLogGroup(groupName); err != nil {
		return nil, NewLogsError("ResourceNotFoundException",
			"The specified resource does not exist: "+input.ResourceArn, 400)
	}

	source := &logsstore.DeliverySource{
		Name:                        input.Name,
		Arn:                         store.ARNBuilder().CloudWatch().DeliverySource(input.Name),
		ResourceArn:                 input.ResourceArn,
		Service:                     parsed.Service,
		LogType:                     input.LogType,
		DeliverySourceConfiguration: input.DeliverySourceConfiguration,
		Tags:                        input.Tags,
	}
	if err := store.PutDeliverySource(source); err != nil {
		return nil, mapStoreError(err)
	}
	return s.deliverySourceStatus(input.Region, source), nil
}

// DeliverySourceView carries one delivery source plus its computed
// status members: ACTIVE while the underlying resource exists, INACTIVE
// with statusReason RESOURCE_DELETED otherwise ("A status reason of
// RESOURCE_DELETED indicates that the resource associated with the
// delivery source has been deleted", member documentation). The status
// is read-side state, so the core computes it and the formatter renders
// it.
type DeliverySourceView struct {
	Source       *logsstore.DeliverySource
	Status       string
	StatusReason string
}

// deliverySourceStatus computes the status members from the underlying
// resource's existence. Core-only: the handler closure never reaches it.
func (s *LogsService) deliverySourceStatus(region string, source *logsstore.DeliverySource) *DeliverySourceView {
	store, err := s.getLogsStoreByRegion(region)
	if err == nil {
		if _, err := store.GetLogGroup(resolveLogGroupIdentifier(source.ResourceArn)); err == nil {
			return &DeliverySourceView{Source: source, Status: "ACTIVE"}
		}
	}
	return &DeliverySourceView{Source: source, Status: "INACTIVE", StatusReason: "RESOURCE_DELETED"}
}

// getDeliverySourceCore resolves one delivery source by name and its
// computed status.
func (s *LogsService) getDeliverySourceCore(region, name string) (*DeliverySourceView, error) {
	source, err := resolveDeliveryFamilyOne(s, region, func(store *logsstore.Store) (*logsstore.DeliverySource, error) {
		return store.GetDeliverySource(name)
	})
	if err != nil {
		return nil, err
	}
	return s.deliverySourceStatus(region, source), nil
}

// describeDeliveryFamilyPage serves one page of a delivery-family
// listing: the bound validates, the family's full listing loads, and the
// name-or-id-ordered page rides the pagination helper — one body behind
// the family's three Describe cores. The page marker rides the scoped
// listing vocabulary: the token carries the family's request identity
// and its documented expiry rather than a bare natural key.
func describeDeliveryFamilyPage[T any, V any](s *LogsService, region, family, nextToken string, limit int32,
	list func(*logsstore.Store) ([]T, error), keyOf func(T) string, toView func(T) V) ([]V, string, error) {
	limit, err := validateListLimitValidation(limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}
	items, err := list(store)
	if err != nil {
		return nil, "", mapStoreError(err)
	}
	scope := listingScope(family)
	result, err := paginateScopedListing(scope, nextToken, items, int(limit), keyOf)
	if err != nil {
		return nil, "", err
	}
	views := make([]V, len(result.Items))
	for i, item := range result.Items {
		views[i] = toView(item)
	}
	return views, result.NextMarker, nil
}

// describeDeliverySourcesCore lists the account's delivery sources
// (name-ordered) one page at a time with their computed statuses.
func (s *LogsService) describeDeliverySourcesCore(region, nextToken string, limit int32) ([]*DeliverySourceView, string, error) {
	return describeDeliveryFamilyPage(s, region, "deliverysources", nextToken, limit,
		func(store *logsstore.Store) ([]*logsstore.DeliverySource, error) { return store.ListDeliverySources() },
		func(src *logsstore.DeliverySource) string { return src.Name },
		func(src *logsstore.DeliverySource) *DeliverySourceView { return s.deliverySourceStatus(region, src) })
}

// describeDeliveryDestinationsCore lists the account's delivery
// destinations (name-ordered) one page at a time.
func (s *LogsService) describeDeliveryDestinationsCore(region, nextToken string, limit int32) ([]*logsstore.DeliveryDestination, string, error) {
	return describeDeliveryFamilyPage(s, region, "deliverydestinations", nextToken, limit,
		func(store *logsstore.Store) ([]*logsstore.DeliveryDestination, error) {
			return store.ListDeliveryDestinations()
		},
		func(dest *logsstore.DeliveryDestination) string { return dest.Name },
		func(dest *logsstore.DeliveryDestination) *logsstore.DeliveryDestination { return dest })
}

// resolveDeliveryFamilyOne loads one delivery-family record through the
// region's store — one body behind the family's Get cores.
func resolveDeliveryFamilyOne[T any](s *LogsService, region string, get func(*logsstore.Store) (T, error)) (T, error) {
	var zero T
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return zero, err
	}
	one, err := get(store)
	if err != nil {
		return zero, mapStoreError(err)
	}
	return one, nil
}

// getDeliveryDestinationCore resolves one delivery destination by name.
func (s *LogsService) getDeliveryDestinationCore(region, name string) (*logsstore.DeliveryDestination, error) {
	return resolveDeliveryFamilyOne(s, region, func(store *logsstore.Store) (*logsstore.DeliveryDestination, error) {
		return store.GetDeliveryDestination(name)
	})
}

// describeDeliveriesCore lists the account's deliveries (id-ordered)
// one page at a time.
func (s *LogsService) describeDeliveriesCore(region, nextToken string, limit int32) ([]*logsstore.Delivery, string, error) {
	return describeDeliveryFamilyPage(s, region, "deliveries", nextToken, limit,
		func(store *logsstore.Store) ([]*logsstore.Delivery, error) { return store.ListDeliveries() },
		func(d *logsstore.Delivery) string { return d.Id },
		func(d *logsstore.Delivery) *logsstore.Delivery { return d })
}

// getDeliveryCore resolves one delivery by id.
func (s *LogsService) getDeliveryCore(region, id string) (*logsstore.Delivery, error) {
	return resolveDeliveryFamilyOne(s, region, func(store *logsstore.Store) (*logsstore.Delivery, error) {
		return store.GetDelivery(id)
	})
}

// refuseDeliveryFamilyIfReferenced guards a family delete: a record
// current deliveries still reference refuses with the family's conflict
// error ("You can't delete a delivery source if any current deliveries
// are associated with it", operation documentation — the destination
// twin carries the same rule).
func refuseDeliveryFamilyIfReferenced(deliveries []*logsstore.Delivery, referenced func(*logsstore.Delivery) bool, label, name string) error {
	for _, d := range deliveries {
		if referenced(d) {
			return errDeliveryConflict(label + " " + name + " is referenced by delivery " + d.Id)
		}
	}
	return nil
}

// deleteDeliverySourceCore refuses a source that deliveries still
// reference ("You can't delete a delivery source if any current
// deliveries are associated with it", operation documentation).
func (s *LogsService) deleteDeliverySourceCore(region, name string) error {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	if _, err := store.GetDeliverySource(name); err != nil {
		return mapStoreError(err)
	}
	deliveries, err := store.ListDeliveries()
	if err != nil {
		return mapStoreError(err)
	}
	if err := refuseDeliveryFamilyIfReferenced(deliveries,
		func(d *logsstore.Delivery) bool { return d.DeliverySourceName == name },
		"Delivery source", name); err != nil {
		return err
	}
	return mapStoreError(store.DeleteDeliverySource(name))
}
