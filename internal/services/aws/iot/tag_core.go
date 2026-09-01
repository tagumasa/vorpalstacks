package iot

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Tag Core. The tag trio delegates wire parsing and validation to the
// shared tag handler framework; these Cores carry the persistence leg so
// the store is reached only from Core functions.
// ---------------------------------------------------------------------------

// tagResourceCore applies a tag set to an IoT resource key.
func (s *IoTService) tagResourceCore(store iotstore.TagOps, resourceKey string, tags map[string]string) error {
	return store.TagResource(resourceKey, tags)
}

// untagResourceCore removes tag keys from an IoT resource key.
func (s *IoTService) untagResourceCore(store iotstore.TagOps, resourceKey string, tagKeys []string) error {
	return store.UntagResource(resourceKey, tagKeys)
}

// listTagsCore reads the tag set of an IoT resource key.
func (s *IoTService) listTagsCore(store iotstore.TagOps, resourceKey string) (map[string]string, error) {
	return store.ListTags(resourceKey)
}

// validateTagResourceCore resolves the IoT resource behind a tag resource
// ARN so the tag trio rejects ARNs addressing no existing resource with the
// ResourceNotFoundException (404) the service model attaches to every tag
// operation. Every ARN kind the platform hosts routes to its owning
// lookup; an ARN of any other kind addresses no resource this service can
// own and fails the same way.
func (s *IoTService) validateTagResourceCore(store iotstore.IotStoreInterface, resourceArn string) error {
	notFound := func() error {
		return awserrors.NewAWSError("ResourceNotFoundException", "The specified resource does not exist.", 404)
	}
	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	kind, name, ok := strings.Cut(resource, "/")
	if !ok || name == "" {
		return notFound()
	}
	var err error
	switch kind {
	case "thing":
		_, err = store.GetThing(name)
	case "thingtype":
		_, err = store.GetThingType(name)
	case "thinggroup":
		_, err = store.GetThingGroup(name)
	case "billinggroup":
		_, err = store.GetBillingGroup(name)
	case "cert":
		_, err = store.GetCertificate(name)
	case "cacert":
		var exists bool
		exists, err = store.GetGenericExists("caCert/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrCertificateNotFound
		}
	case "policy":
		_, err = store.GetPolicy(name)
	case "rule":
		_, err = store.GetRule(name)
	case "rolealias":
		_, err = store.GetRoleAlias(name)
	case "job":
		_, err = store.GetJob(name)
	case "authorizer":
		_, err = store.GetAuthorizer(name)
	case "provisioningtemplate":
		_, err = store.GetProvisioningTemplate(name)
	case "securityprofile":
		_, err = store.GetSecurityProfile(name)
	case "domainconfiguration":
		_, err = store.GetDomainConfiguration(name)
	case "custommetric":
		var exists bool
		exists, err = store.GetGenericExists("customMetric/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrCustomMetricNotFound
		}
	case "dimension":
		var exists bool
		exists, err = store.GetGenericExists("dimension/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrDimensionNotFound
		}
	case "mitigationaction":
		var exists bool
		exists, err = store.GetGenericExists("mitigationAction/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrMitigationActionNotFound
		}
	case "scheduledaudit":
		var exists bool
		exists, err = store.GetGenericExists("scheduledAudit/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrScheduledAuditNotFound
		}
	case "fleetmetric":
		var exists bool
		exists, err = store.GetGenericExists("fleetMetric/"+name, &map[string]interface{}{})
		if err == nil && !exists {
			return iotstore.ErrFleetMetricNotFound
		}
	default:
		return notFound()
	}
	return err
}

// iotTagConfig builds a TagHandlerConfig backed by the IoT store Cores.
// It lives in the Core layer so the ValidateResource/TagFunc/UntagFunc/
// ListFunc store calls stay out of every handler closure.
func (s *IoTService) iotTagConfig(store iotstore.IotStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(_ context.Context, resourceKey string) error {
			return s.validateTagResourceCore(store, resourceKey)
		},
		TagFunc: func(_ context.Context, resourceKey string, tags []tagutil.Tag) error {
			return s.tagResourceCore(store, resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return s.untagResourceCore(store, resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			tagsMap, err := s.listTagsCore(store, resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(tagsMap), nil
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
			tagList := make([]map[string]interface{}, 0, len(tags))
			for _, t := range tags {
				tagList = append(tagList, map[string]interface{}{
					"Key":   t.Key,
					"Value": t.Value,
				})
			}
			return map[string]interface{}{"tags": tagList}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		MapError: func(err error) error {
			if _, ok := err.(*tagutil.MissingResourceError); ok {
				return awserrors.NewAWSError("ResourceNotFoundException", "The specified resource does not exist.", 404)
			}
			return err
		},
	}
}
