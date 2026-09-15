package sqs

import (
	"context"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Tag-operation Core — the queue-tag wiring for the shared tag machinery.
// Validation, resolution and persistence closures live here so the handlers
// in tag_operations.go stay thin adapters.
// ---------------------------------------------------------------------------

func sqsMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return ErrMissingParameter
	case *tagutil.MissingTagsError:
		return ErrMissingParameter
	case *tagutil.MissingTagKeysError:
		return ErrMissingParameter
	}
	return convertStoreError(err)
}

func sqsTagConfig(store sqsstore.SQSStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.SQSConfig,
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			_, err := store.GetQueue(resourceKey)
			if err != nil {
				return convertStoreError(err)
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			return store.TagQueue(resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.UntagQueue(resourceKey, tagKeys)
		},
		// The TagKeys member carries xmlName "TagKey" + xmlFlattened: the
		// query wire spells the indexed list TagKey.N, which the generic
		// fallback (TagKeys.member.N / TagKeys) never tries. Without this
		// override a query-wire UntagQueue parses zero keys and the
		// len(tagKeys) > 0 guard turns the request into a silent no-op.
		// The stem read rejects a gapped list; this hook cannot carry the
		// error, and a nil return makes the framework's RequireTagKeys rule
		// reject the request instead — never a silent truncation.
		ParseTagKeys: func(params map[string]interface{}) []string {
			keys := tagutil.GetTagKeys(params, tagutil.SQSConfig)
			if len(keys) > 0 {
				return keys
			}
			stemKeys, err := getQueryListByStem(params, "TagKey")
			if err != nil {
				return nil
			}
			return stemKeys
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			m, err := store.ListQueueTags(resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []tagutil.Tag, rawResourceKey string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToMap(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: sqsMapError,
	}
}
