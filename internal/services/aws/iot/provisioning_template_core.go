package iot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Provisioning Template Core (fleet provisioning templates and their
// versions). Templates use the typed ProvisioningTemplate store; the
// template body must be well-formed JSON before it is stored.
// ---------------------------------------------------------------------------

// CreateProvisioningTemplateInput carries the parsed
// CreateProvisioningTemplate request. EnabledProvided distinguishes an
// explicitly supplied enabled flag from an omitted one (default true).
type CreateProvisioningTemplateInput struct {
	TemplateName        string
	Description         string
	TemplateBody        string
	Enabled             bool
	EnabledProvided     bool
	ProvisioningRoleARN string
	Type                string
	Tags                map[string]string
}

// templateTypes is the TemplateType enum member set.
var templateTypes = map[string]bool{"FLEET_PROVISIONING": true, "JITP": true}

// createProvisioningTemplateCore validates and persists a provisioning
// template and returns the created record. The templateBody and
// provisioningRoleArn members are required by the model; type is optional
// with a FLEET_PROVISIONING default and validated against the enum.
func (s *IoTService) createProvisioningTemplateCore(store iotstore.IotStoreInterface, in CreateProvisioningTemplateInput) (*iotstore.ProvisioningTemplate, error) {
	if in.TemplateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.TemplateBody == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.ProvisioningRoleARN == "" {
		return nil, iotstore.ErrMissingParam
	}
	if in.Type != "" && !templateTypes[in.Type] {
		return nil, iotstore.ErrInvalidRequest
	}
	// Validate template body is well-formed JSON before storing.
	var v interface{}
	if err := json.Unmarshal([]byte(in.TemplateBody), &v); err != nil {
		return nil, fmt.Errorf("invalid templateBody: %w", err)
	}
	enabled := true
	if in.EnabledProvided {
		enabled = in.Enabled
	}
	now := time.Now().UTC()
	tmpl := &iotstore.ProvisioningTemplate{
		TemplateName:        in.TemplateName,
		Description:         in.Description,
		Enabled:             enabled,
		ProvisioningRoleARN: in.ProvisioningRoleARN,
		TemplateBody:        in.TemplateBody,
		Type:                in.Type,
		CreationDate:        now,
		LastModifiedDate:    now,
	}
	created, err := store.CreateProvisioningTemplate(tmpl)
	if err != nil {
		return nil, err
	}
	// Create-time tags live in the ARN-keyed tag store so
	// ListTagsForResource observes them; the delete path already clears
	// them with DeleteAllTags.
	if len(in.Tags) > 0 {
		if err := store.TagResource(created.TemplateARN, in.Tags); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// describeProvisioningTemplateCore loads a provisioning template.
func (s *IoTService) describeProvisioningTemplateCore(store iotstore.IotStoreInterface, name string) (*iotstore.ProvisioningTemplate, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetProvisioningTemplate(name)
}

// UpdateProvisioningTemplateInput carries the parsed
// UpdateProvisioningTemplate request. The pointer members follow the store
// partial-update convention (nil = no change).
type UpdateProvisioningTemplateInput struct {
	TemplateName              string
	Description               string
	RoleARN                   string
	Enabled                   *bool
	TemplateBody              string
	DefaultVersionID          int64
	PreProvisioningHook       string
	RemovePreProvisioningHook bool
}

// updateProvisioningTemplateCore applies a partial update to a provisioning
// template.
func (s *IoTService) updateProvisioningTemplateCore(store iotstore.IotStoreInterface, in UpdateProvisioningTemplateInput) error {
	if in.TemplateName == "" {
		return iotstore.ErrMissingParam
	}
	opts := iotstore.ProvisioningTemplateUpdateOpts{
		Description: in.Description,
		RoleARN:     in.RoleARN,
	}
	opts.Enabled = in.Enabled
	if in.TemplateBody != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(in.TemplateBody), &v); err != nil {
			return fmt.Errorf("invalid templateBody: %w", err)
		}
		opts.TemplateBody = &in.TemplateBody
	}
	opts.DefaultVersionID = in.DefaultVersionID
	opts.PreProvisioningHook = in.PreProvisioningHook
	opts.RemovePreProvisioningHook = in.RemovePreProvisioningHook
	_, err := store.UpdateProvisioningTemplate(in.TemplateName, opts)
	return err
}

// deleteProvisioningTemplateCore removes a provisioning template and its
// tags.
func (s *IoTService) deleteProvisioningTemplateCore(store iotstore.IotStoreInterface, name string) error {
	if name == "" {
		return iotstore.ErrMissingParam
	}
	arn := iotstore.BuildProvisioningTemplateARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return store.DeleteProvisioningTemplate(name)
}

// listProvisioningTemplatesCore lists provisioning templates.
func (s *IoTService) listProvisioningTemplatesCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions) ([]*iotstore.ProvisioningTemplate, string, error) {
	templates, err := store.ListProvisioningTemplates(opts)
	if err != nil {
		return nil, "", err
	}
	return templates.Items, templates.NextMarker, nil
}

// ProvisioningTemplateVersionItem is one ListProvisioningTemplateVersions
// entry: the model's ProvisioningTemplateVersionSummary member set.
type ProvisioningTemplateVersionItem struct {
	VersionID        int32
	IsDefaultVersion bool
	CreationDate     int64
}

// listProvisioningTemplateVersionsCore lists a template's versions in
// numeric version order. The store iterates keys lexicographically ("1",
// "10", "11", "2"); version IDs are numeric and must be listed in numeric
// order. Version IDs that fail numeric parsing are skipped (warned).
func (s *IoTService) listProvisioningTemplateVersionsCore(store iotstore.IotStoreInterface, name string, opts storecommon.ListOptions) ([]ProvisioningTemplateVersionItem, error) {
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	versions, err := store.ListProvisioningTemplateVersions(name, opts)
	if err != nil {
		return nil, err
	}
	items := make([]ProvisioningTemplateVersionItem, 0, len(versions))
	for _, v := range versions {
		var vidInt int32
		if _, scanErr := fmt.Sscanf(v.VersionID, "%d", &vidInt); scanErr != nil {
			slog.Warn("provisioning template version ID parse failed", "versionID", v.VersionID, "error", scanErr)
			continue
		}
		items = append(items, ProvisioningTemplateVersionItem{
			VersionID:        vidInt,
			IsDefaultVersion: v.IsDefaultVersion,
			CreationDate:     v.CreationDate.UTC().Unix(),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].VersionID < items[j].VersionID
	})
	return items, nil
}

// CreateProvisioningTemplateVersionInput carries the parsed
// CreateProvisioningTemplateVersion request.
type CreateProvisioningTemplateVersionInput struct {
	TemplateName string
	TemplateBody string
	SetAsDefault bool
	ListOpts     storecommon.ListOptions
}

// CreateProvisioningTemplateVersionResult is the transport-agnostic result
// of CreateProvisioningTemplateVersion.
type CreateProvisioningTemplateVersionResult struct {
	TemplateArn      string
	TemplateName     string
	VersionID        int32
	IsDefaultVersion bool
}

// createProvisioningTemplateVersionCore validates and persists a new
// template version, optionally repointing the template's default version.
// An unknown template yields ErrTemplateNotFound.
func (s *IoTService) createProvisioningTemplateVersionCore(store iotstore.IotStoreInterface, in CreateProvisioningTemplateVersionInput) (*CreateProvisioningTemplateVersionResult, error) {
	if in.TemplateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if tmpl, err := store.GetProvisioningTemplate(in.TemplateName); err != nil {
		return nil, err
	} else if tmpl == nil || tmpl.TemplateName == "" {
		return nil, iotstore.ErrTemplateNotFound
	}
	// Determine the next version ID by scanning existing versions.
	existing, err := store.ListProvisioningTemplateVersions(in.TemplateName, in.ListOpts)
	if err != nil {
		return nil, err
	}
	maxVersion := 0
	for _, v := range existing {
		var n int
		fmt.Sscanf(v.VersionID, "%d", &n)
		if n > maxVersion {
			maxVersion = n
		}
	}
	versionID := fmt.Sprintf("%d", maxVersion+1)
	v := &iotstore.ProvisioningTemplateVersion{
		VersionID:        versionID,
		TemplateBody:     in.TemplateBody,
		IsDefaultVersion: in.SetAsDefault,
		CreationDate:     time.Now().UTC(),
	}
	if _, err := store.CreateProvisioningTemplateVersion(in.TemplateName, v); err != nil {
		return nil, err
	}
	if v.IsDefaultVersion {
		// setAsDefault must repoint the template's default version and
		// clear the sibling versions' flags in one step; otherwise the new
		// version claims to be the default while the template still
		// resolves to the previous one, and two versions would report
		// themselves as the default at once.
		if err := store.SetDefaultProvisioningTemplateVersion(in.TemplateName, int64(maxVersion+1)); err != nil {
			return nil, err
		}
	}
	return &CreateProvisioningTemplateVersionResult{
		TemplateArn:      iotstore.BuildProvisioningTemplateARN(store.GetAccountID(), store.GetRegion(), in.TemplateName),
		TemplateName:     in.TemplateName,
		VersionID:        int32(maxVersion + 1),
		IsDefaultVersion: v.IsDefaultVersion,
	}, nil
}

// deleteProvisioningTemplateVersionCore removes one template version.
func (s *IoTService) deleteProvisioningTemplateVersionCore(store iotstore.IotStoreInterface, templateName, versionID string) error {
	if templateName == "" || versionID == "" {
		return iotstore.ErrMissingParam
	}
	if _, err := store.GetProvisioningTemplateVersion(templateName, versionID); err != nil {
		return err
	}
	return store.DeleteProvisioningTemplateVersion(templateName, versionID)
}

// DescribeProvisioningTemplateVersionResult is the transport-agnostic
// result of DescribeProvisioningTemplateVersion.
type DescribeProvisioningTemplateVersionResult struct {
	TemplateName     string
	VersionID        int32
	TemplateBody     string
	IsDefaultVersion bool
	CreationDate     int64
}

// describeProvisioningTemplateVersionCore loads one template version.
func (s *IoTService) describeProvisioningTemplateVersionCore(store iotstore.IotStoreInterface, templateName, versionID string) (*DescribeProvisioningTemplateVersionResult, error) {
	if templateName == "" || versionID == "" {
		return nil, iotstore.ErrMissingParam
	}
	v, err := store.GetProvisioningTemplateVersion(templateName, versionID)
	if err != nil {
		return nil, err
	}
	var versionIDInt int32
	fmt.Sscanf(v.VersionID, "%d", &versionIDInt)
	return &DescribeProvisioningTemplateVersionResult{
		TemplateName:     templateName,
		VersionID:        versionIDInt,
		TemplateBody:     v.TemplateBody,
		IsDefaultVersion: v.IsDefaultVersion,
		CreationDate:     v.CreationDate.UTC().Unix(),
	}, nil
}
