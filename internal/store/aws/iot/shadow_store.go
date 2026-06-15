package iot

import (
		"strings"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func shadowKey(thingName, shadowName string) string {
	if shadowName == "" || shadowName == "classic" {
		return thingName + "/$current"
	}
	return thingName + "/" + shadowName
}

func (s *IotStore) GetShadow(thingName, shadowName string) (*ShadowDocument, error) {
	pb := &pb.ShadowDocument{}
	if err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrShadowNotFound
		}
		return nil, err
	}
	return ProtoToShadow(pb), nil
}

func (s *IotStore) PutShadow(thingName, shadowName string, doc *ShadowDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

// PutShadowWithVersion atomically reads the current shadow, checks the version,
// and writes the new document if the version matches.  Returns ErrVersionConflict
// if clientVersion > 0 and does not match the stored version.
func (s *IotStore) PutShadowWithVersion(thingName, shadowName string, doc *ShadowDocument, clientVersion int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pb := &pb.ShadowDocument{}
	err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb)
	if err == nil {
		stored := ProtoToShadow(pb)
		if clientVersion > 0 && stored.VersionNumber != clientVersion {
			return ErrVersionConflict
		}
	} else if !common.IsNotFound(err) {
		return err
	}
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

func (s *IotStore) DeleteShadow(thingName, shadowName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shadowsBase.Delete(shadowKey(thingName, shadowName))
}

// ListShadowNames returns the distinct shadow names for a given thing by
// scanning the shadows bucket and extracting the shadow name component.
func (s *IotStore) ListShadowNames(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prefix := thingName + "/"
	var names []string
	seen := make(map[string]bool)
	err := s.shadowsBase.ScanPrefix(prefix, func(key string, _ []byte) error {
		parts := strings.SplitN(key[len(prefix):], "/", 2)
		name := parts[0]
		if name == "$current" {
			return nil
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
		return nil
	})
	return names, err
}
