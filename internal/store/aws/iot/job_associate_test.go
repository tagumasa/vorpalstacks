package iot

import (
	"testing"

	"vorpalstacks/internal/store/aws/common"
)

const (
	thingA = "arn:aws:iot:us-east-1:000000000000:thing/thing-a"
	thingB = "arn:aws:iot:us-east-1:000000000000:thing/thing-b"
)

// Associating additional targets with a job must materialise a QUEUED
// execution record for each genuinely new thing, exactly as CreateJob
// does, so associated targets show up in job execution listings instead
// of being silently merged into the job's target list.
func TestAssociateJobTargetsCreatesExecutions(t *testing.T) {
	st := newIotStore(t)

	if _, err := st.CreateJob(&Job{JobID: "job-1", Targets: []string{thingA}}); err != nil {
		t.Fatal(err)
	}

	job, err := st.AssociateJobTargets("job-1", []string{thingB, thingA}, "added later")
	if err != nil {
		t.Fatal(err)
	}

	if len(job.Targets) != 2 {
		t.Fatalf("expected the merged target list to hold 2 entries, got %v", job.Targets)
	}
	if job.Description != "added later" {
		t.Fatalf("comment not recorded: %q", job.Description)
	}

	executions, err := st.ListGeneric("jobExecution/job-1/")
	if err != nil {
		t.Fatal(err)
	}
	byThing := map[string]bool{}
	for _, rec := range executions {
		name, _ := rec["thingName"].(string)
		status, _ := rec["status"].(string)
		if status != "QUEUED" {
			t.Fatalf("execution for %s has status %q, want QUEUED", name, status)
		}
		byThing[name] = true
	}
	if !byThing["thing-a"] || !byThing["thing-b"] {
		t.Fatalf("execution records missing after associate: %v", byThing)
	}
	if len(executions) != 2 {
		t.Fatalf("duplicate execution records created: %d", len(executions))
	}
}

// Switching a template's default version must leave exactly one version
// flagged as the default: the newly designated one, with every sibling's
// flag cleared.
func TestSetDefaultProvisioningTemplateVersionLeavesSingleFlag(t *testing.T) {
	st := newIotStore(t)

	if _, err := st.CreateProvisioningTemplate(&ProvisioningTemplate{
		TemplateName: "tpl",
		TemplateBody: "{}",
		Enabled:      true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateProvisioningTemplateVersion("tpl", &ProvisioningTemplateVersion{
		VersionID:        "1",
		TemplateBody:     "{}",
		IsDefaultVersion: true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateProvisioningTemplateVersion("tpl", &ProvisioningTemplateVersion{
		VersionID:        "2",
		TemplateBody:     "{\"v\":2}",
		IsDefaultVersion: false,
	}); err != nil {
		t.Fatal(err)
	}

	if err := st.SetDefaultProvisioningTemplateVersion("tpl", 2); err != nil {
		t.Fatal(err)
	}

	versions, err := st.ListProvisioningTemplateVersions("tpl", common.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defaults := 0
	for _, v := range versions {
		if !v.IsDefaultVersion {
			continue
		}
		defaults++
		if v.VersionID != "2" {
			t.Fatalf("version %s flagged as default, want version 2", v.VersionID)
		}
	}
	if defaults != 1 {
		t.Fatalf("expected exactly one default version, got %d", defaults)
	}

	template, err := st.GetProvisioningTemplate("tpl")
	if err != nil {
		t.Fatal(err)
	}
	if template.DefaultVersionID != 2 {
		t.Fatalf("template DefaultVersionID = %d, want 2", template.DefaultVersionID)
	}
}

// Repointing the default to a version that does not exist must fail before
// any mutation: the previous default's flag and the template's
// DefaultVersionID stay untouched instead of leaving the template with no
// default at all.
func TestSetDefaultProvisioningTemplateVersionRejectsUnknownVersion(t *testing.T) {
	st := newIotStore(t)

	if _, err := st.CreateProvisioningTemplate(&ProvisioningTemplate{
		TemplateName: "tpl",
		TemplateBody: "{}",
		Enabled:      true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateProvisioningTemplateVersion("tpl", &ProvisioningTemplateVersion{
		VersionID:        "1",
		TemplateBody:     "{}",
		IsDefaultVersion: true,
	}); err != nil {
		t.Fatal(err)
	}

	if err := st.SetDefaultProvisioningTemplateVersion("tpl", 99); err != ErrTemplateVersionNotFound {
		t.Fatalf("unknown versionID error = %v, want ErrTemplateVersionNotFound", err)
	}

	versions, err := st.ListProvisioningTemplateVersions("tpl", common.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 1 || !versions[0].IsDefaultVersion {
		t.Fatalf("existing default flag was disturbed: %+v", versions)
	}
	template, err := st.GetProvisioningTemplate("tpl")
	if err != nil {
		t.Fatal(err)
	}
	if template.DefaultVersionID == 99 {
		t.Fatal("template DefaultVersionID was repointed to a non-existent version")
	}
}

// The same existence contract applies to UpdateProvisioningTemplate's
// defaultVersionId, and a valid repoint moves the versions' flags so the
// stored default set never disagrees with the template record.
func TestUpdateProvisioningTemplateValidatesDefaultVersionID(t *testing.T) {
	st := newIotStore(t)

	if _, err := st.CreateProvisioningTemplate(&ProvisioningTemplate{
		TemplateName: "tpl",
		TemplateBody: "{}",
		Enabled:      true,
	}); err != nil {
		t.Fatal(err)
	}
	for _, vid := range []string{"1", "2"} {
		isDefault := vid == "1"
		if _, err := st.CreateProvisioningTemplateVersion("tpl", &ProvisioningTemplateVersion{
			VersionID:        vid,
			TemplateBody:     "{}",
			IsDefaultVersion: isDefault,
		}); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := st.UpdateProvisioningTemplate("tpl", ProvisioningTemplateUpdateOpts{DefaultVersionID: 99}); err != ErrTemplateVersionNotFound {
		t.Fatalf("unknown defaultVersionId error = %v, want ErrTemplateVersionNotFound", err)
	}
	template, err := st.GetProvisioningTemplate("tpl")
	if err != nil {
		t.Fatal(err)
	}
	if template.DefaultVersionID == 99 {
		t.Fatal("template DefaultVersionID was repointed to a non-existent version")
	}

	if _, err := st.UpdateProvisioningTemplate("tpl", ProvisioningTemplateUpdateOpts{DefaultVersionID: 2}); err != nil {
		t.Fatal(err)
	}
	versions, err := st.ListProvisioningTemplateVersions("tpl", common.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defaults := 0
	for _, v := range versions {
		if v.IsDefaultVersion {
			defaults++
			if v.VersionID != "2" {
				t.Fatalf("version %s flagged as default after update, want version 2", v.VersionID)
			}
		}
	}
	if defaults != 1 {
		t.Fatalf("expected exactly one default version after update, got %d", defaults)
	}
}
