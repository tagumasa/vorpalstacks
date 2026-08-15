package cognitoidentityprovider

import (
	"encoding/json"
	"strconv"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// GetLogDeliveryConfiguration retrieves the log delivery configuration for a user pool.
func (s *CognitoStore) GetLogDeliveryConfiguration(userPoolID string) (*LogDeliveryConfiguration, error) {
	var cfg LogDeliveryConfiguration
	if err := s.BaseStore.Get(logDeliveryKey(userPoolID), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveLogDeliveryConfiguration stores the log delivery configuration for a user pool.
func (s *CognitoStore) SaveLogDeliveryConfiguration(cfg *LogDeliveryConfiguration) error {
	return s.BaseStore.Put(logDeliveryKey(cfg.UserPoolID), cfg)
}

// GetRiskConfiguration retrieves the risk configuration for a user pool/client.
func (s *CognitoStore) GetRiskConfiguration(userPoolID, clientID string) (*RiskConfiguration, error) {
	var cfg RiskConfiguration
	if err := s.BaseStore.Get(riskConfigKey(userPoolID, clientID), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveRiskConfiguration stores the risk configuration for a user pool/client.
func (s *CognitoStore) SaveRiskConfiguration(cfg *RiskConfiguration) error {
	cfg.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(riskConfigKey(cfg.UserPoolID, cfg.ClientID), cfg)
}

// GetUICustomization retrieves UI customisation for a user pool/client.
func (s *CognitoStore) GetUICustomization(userPoolID, clientID string) (*UICustomization, error) {
	var ui UICustomization
	if err := s.BaseStore.Get(uiCustomizationKey(userPoolID, clientID), &ui); err != nil {
		return nil, err
	}
	return &ui, nil
}

// SaveUICustomization stores UI customisation for a user pool/client.
// CSSVersion is incremented only when the CSS content actually changes,
// matching AWS behaviour for cache-busting on the hosted UI.
func (s *CognitoStore) SaveUICustomization(ui *UICustomization) error {
	key := uiCustomizationKey(ui.UserPoolID, ui.ClientID)
	now := time.Now().UTC()

	var prev UICustomization
	hasPrev := s.BaseStore.Get(key, &prev) == nil

	if !hasPrev {
		ui.CreationDate = now
		ui.CSSVersion = "20200101"
	} else {
		if ui.CreationDate.IsZero() {
			ui.CreationDate = prev.CreationDate
		}
		if ui.CSS != prev.CSS {
			if prev.CSSVersion == "" {
				ui.CSSVersion = "20200101"
			} else if v, err := strconv.Atoi(prev.CSSVersion); err == nil {
				ui.CSSVersion = strconv.Itoa(v + 1)
			} else {
				ui.CSSVersion = "20200101"
			}
		} else {
			ui.CSSVersion = prev.CSSVersion
		}
	}

	ui.LastModifiedDate = now
	return s.BaseStore.Put(key, ui)
}

// CreateUserImportJob stores a new user import job.
func (s *CognitoStore) CreateUserImportJob(job *UserImportJob) error {
	return s.userImportJobsStore.Put(userImportJobKey(job.UserPoolID, job.JobID), job)
}

// GetUserImportJob retrieves an import job by ID.
func (s *CognitoStore) GetUserImportJob(userPoolID, jobID string) (*UserImportJob, error) {
	var job UserImportJob
	if err := s.userImportJobsStore.Get(userImportJobKey(userPoolID, jobID), &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateUserImportJob updates an import job.
func (s *CognitoStore) UpdateUserImportJob(job *UserImportJob) error {
	return s.userImportJobsStore.Put(userImportJobKey(job.UserPoolID, job.JobID), job)
}

// ListUserImportJobs lists import jobs for a user pool.
func (s *CognitoStore) ListUserImportJobsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserImportJob], error) {
	opts.Prefix = userImportJobPrefix(userPoolID)
	return common.List[UserImportJob](s.userImportJobsStore, opts, nil)
}

// ===================== Managed Login Branding =====================

func (s *CognitoStore) SaveManagedLoginBranding(b *ManagedLoginBranding) error {
	if b.CreationDate.IsZero() {
		b.CreationDate = time.Now().UTC()
	}
	b.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(managedLoginBrandingKey(b.UserPoolID, b.ManagedLoginBrandingId), b)
}

func (s *CognitoStore) GetManagedLoginBranding(userPoolID, brandingID string) (*ManagedLoginBranding, error) {
	var b ManagedLoginBranding
	if err := s.BaseStore.Get(managedLoginBrandingKey(userPoolID, brandingID), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *CognitoStore) GetManagedLoginBrandingByClient(userPoolID, clientID string) (*ManagedLoginBranding, error) {
	var found *ManagedLoginBranding
	prefix := managedLoginBrandingPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var b ManagedLoginBranding
		if err := json.Unmarshal(value, &b); err != nil {
			return err
		}
		if b.ClientID == clientID {
			found = &b
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrNotFound
	}
	return found, nil
}

func (s *CognitoStore) DeleteManagedLoginBranding(userPoolID, brandingID string) error {
	return s.BaseStore.Delete(managedLoginBrandingKey(userPoolID, brandingID))
}

func (s *CognitoStore) ListManagedLoginBrandings(userPoolID string) ([]*ManagedLoginBranding, error) {
	var brandings []*ManagedLoginBranding
	prefix := managedLoginBrandingPrefix(userPoolID)
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var b ManagedLoginBranding
		if err := json.Unmarshal(value, &b); err != nil {
			return err
		}
		brandings = append(brandings, &b)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return brandings, nil
}

// ===================== Terms =====================

func (s *CognitoStore) SaveTerms(t *Terms) error {
	if t.CreationDate.IsZero() {
		t.CreationDate = time.Now().UTC()
	}
	t.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(termsKey(t.UserPoolID, t.TermsID), t)
}

func (s *CognitoStore) GetTerms(userPoolID, termsID string) (*Terms, error) {
	var t Terms
	if err := s.BaseStore.Get(termsKey(userPoolID, termsID), &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *CognitoStore) DeleteTerms(userPoolID, termsID string) error {
	return s.BaseStore.Delete(termsKey(userPoolID, termsID))
}

func (s *CognitoStore) ListTermsPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[Terms], error) {
	opts.Prefix = termsPrefix(userPoolID)
	return common.List[Terms](s.BaseStore, opts, nil)
}

// ===================== User Pool Replicas =====================

func (s *CognitoStore) SaveUserPoolReplica(r *UserPoolReplica) error {
	return s.BaseStore.Put(userPoolReplicaKey(r.UserPoolID, r.RegionName), r)
}

func (s *CognitoStore) GetUserPoolReplica(userPoolID, regionName string) (*UserPoolReplica, error) {
	var r UserPoolReplica
	if err := s.BaseStore.Get(userPoolReplicaKey(userPoolID, regionName), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CognitoStore) DeleteUserPoolReplica(userPoolID, regionName string) error {
	return s.BaseStore.Delete(userPoolReplicaKey(userPoolID, regionName))
}

func (s *CognitoStore) ListUserPoolReplicasPaginated(userPoolID string, opts common.ListOptions) (*common.ListResult[UserPoolReplica], error) {
	opts.Prefix = userPoolReplicaPrefix(userPoolID)
	return common.List[UserPoolReplica](s.BaseStore, opts, nil)
}
