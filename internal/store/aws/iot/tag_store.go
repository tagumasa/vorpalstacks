package iot

func (s *IotStore) ListTags(resourceARN string) (map[string]string, error) {
	return s.TagStore.List(resourceARN)
}

func (s *IotStore) TagResource(resourceARN string, tags map[string]string) error {
	return s.TagStore.Tag(resourceARN, tags)
}

func (s *IotStore) UntagResource(resourceARN string, tagKeys []string) error {
	return s.TagStore.Untag(resourceARN, tagKeys)
}

func (s *IotStore) DeleteAllTags(resourceARN string) error {
	return s.TagStore.Delete(resourceARN)
}
