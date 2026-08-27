package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
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
