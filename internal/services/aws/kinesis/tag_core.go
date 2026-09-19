package kinesis

import (
	"context"
	"sort"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// resolveTagResourceKey resolves the tag-able resource a request addresses
// — a stream by name, a stream or consumer by ARN — to its tag-store key.
// An ARN that is neither a stream nor a registered consumer falls through
// unresolved; validateTagResource reports it as not-found instead of the
// tag writes silently keying junk. The shared definition: the HTTP tag
// framework's config closure and the ARN-addressed tag cores both route
// through it.
func resolveTagResourceKey(store *kinesisstore.KinesisStore, streamName, streamARN, resourceARN, rawKey string) string {
	if streamName != "" {
		return streamName
	}
	arn := streamARN
	if arn == "" {
		arn = resourceARN
	}
	if arn == "" {
		return rawKey
	}
	if stream, err := store.GetStreamByARN(arn); err == nil {
		return stream.StreamName
	}
	if _, err := store.GetStreamConsumer(arn); err == nil {
		return arn
	}
	return rawKey
}

// validateTagResource reports whether a resolved tag-resource key addresses
// an existing stream or registered consumer.
func validateTagResource(store *kinesisstore.KinesisStore, resourceKey string) error {
	if _, err := store.GetStream(resourceKey); err == nil {
		return nil
	}
	if _, err := store.GetStreamConsumer(resourceKey); err == nil {
		return nil
	}
	return kinesisstore.ErrStreamNotFound
}

// tagResourceCore applies one tag set to the resource an ARN addresses —
// the ARN-addressed half of the tag surface, shared by the HTTP tag
// framework and the admin console plane. The ARN's shape validates here,
// in the shared core, so both planes answer an absent or malformed ARN
// with the same identity instead of drifting between them.
func (s *KinesisService) tagResourceCore(reqCtx *request.RequestContext, resourceARN string, tagList []tags.Tag) error {
	if !validateResourceARN(resourceARN) {
		return ErrInvalidArgument
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	key := resolveTagResourceKey(store, "", "", resourceARN, resourceARN)
	if err := validateTagResource(store, key); err != nil {
		return s.mapStoreError(err)
	}
	if err := tags.ValidateTags(tagList); err != nil {
		return ErrInvalidArgument
	}
	if err := store.Tag(key, tags.ToMap(tagList)); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// untagResourceCore removes tag keys from the resource an ARN addresses.
// Its store failure carries the same mapping the tag core applies, so both
// planes answer a storage failure with the family's mapped identity.
func (s *KinesisService) untagResourceCore(reqCtx *request.RequestContext, resourceARN string, tagKeys []string) error {
	if !validateResourceARN(resourceARN) {
		return ErrInvalidArgument
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	key := resolveTagResourceKey(store, "", "", resourceARN, resourceARN)
	if err := validateTagResource(store, key); err != nil {
		return s.mapStoreError(err)
	}
	if !validateTagKeysList(tagKeys) {
		return ErrInvalidArgument
	}
	if err := store.Untag(key, tagKeys); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// listResourceTagsCore lists the tags of the resource an ARN addresses.
func (s *KinesisService) listResourceTagsCore(reqCtx *request.RequestContext, resourceARN string) ([]tags.Tag, error) {
	if !validateResourceARN(resourceARN) {
		return nil, ErrInvalidArgument
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := resolveTagResourceKey(store, "", "", resourceARN, resourceARN)
	if err := validateTagResource(store, key); err != nil {
		return nil, s.mapStoreError(err)
	}
	tagList, err := store.ListAsSlice(key)
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return tagList, nil
}

// kinesisTagConfig builds the shared tag-handler configuration. The two
// list shapes differ: ListTagsForStream's output carries HasMoreTags
// alongside Tags, while ListTagsForResource's output is Tags only — the
// emitPagination flag selects the caller's shape.
//
// The tag-able resource space is bipartite: streams, whose tag entries key
// by stream name, and registered consumers, whose entries key by consumer
// ARN (the register-time write's key — the model's Tag shape documents
// "Metadata assigned to the stream or consumer"). Every write path applies
// the same shared tag-set validation: the handler framework here, and the
// explicit CreateStream / RegisterStreamConsumer Core checks.
func (s *KinesisService) kinesisTagConfig(store *kinesisstore.KinesisStore, req *request.ParsedRequest, emitPagination bool) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.KinesisStreamConfig,
		ResourceKey: func(rawKey string) string {
			return resolveTagResourceKey(store,
				request.GetParamLowerFirst(req.Parameters, "StreamName"),
				request.GetParamLowerFirst(req.Parameters, "StreamARN"),
				request.GetParamLowerFirst(req.Parameters, "ResourceARN"),
				rawKey)
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			return validateTagResource(store, resourceKey)
		},
		ParseTags: func(params map[string]interface{}) []tags.Tag {
			return tags.ParseTags(params, "Tags")
		},
		ValidateTagsFunc: func(tagList []tags.Tag) error {
			if err := tags.ValidateTags(tagList); err != nil {
				return ErrInvalidArgument
			}
			return nil
		},
		ValidateTagKeysFunc: func(tagKeys []string) error {
			if !validateTagKeysList(tagKeys) {
				return ErrInvalidArgument
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []tags.Tag) error {
			if err := store.Tag(resourceKey, tags.ToMap(tagList)); err != nil {
				return s.mapStoreError(err)
			}
			return nil
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			if err := store.Untag(resourceKey, tagKeys); err != nil {
				return s.mapStoreError(err)
			}
			return nil
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tags.Tag, error) {
			tagList, err := store.ListAsSlice(resourceKey)
			if err != nil {
				return nil, s.mapStoreError(err)
			}
			return tagList, nil
		},
		FormatResponse: func(tagList []tags.Tag, _ string) (interface{}, error) {
			// ListTagsForResource models ResourceARN alone — it carries no
			// Limit or ExclusiveStartTagKey member, so the whole
			// (quota-bounded) set is its response. The pagination members
			// belong to ListTagsForStream's shape alone.
			if !emitPagination {
				return map[string]interface{}{
					"Tags": tags.ToResponse(tagList),
				}, nil
			}
			// Implement pagination with ExclusiveStartTagKey and Limit.
			// The member's range trait bounds it 1-50; a value outside the
			// window is rejected, and an absent one serves one full page.
			// The member targets the TagKey shape (length 1-128): an
			// explicitly present member that is empty or past the shape's
			// ceiling is out of range, not an absent filter.
			startKey := request.GetStringParam(req.Parameters, "ExclusiveStartTagKey")
			if request.HasParam(req.Parameters, "ExclusiveStartTagKey") &&
				(startKey == "" || utf8.RuneCountInString(startKey) > tags.MaxTagKeyLength) {
				return nil, ErrInvalidArgument
			}
			limit, hasLimit, err := strictIntParam(req.Parameters, "Limit")
			if err != nil {
				return nil, err
			}
			if hasLimit && (limit < 1 || limit > kinesisstore.MaxTagPageResults) {
				return nil, ErrInvalidArgument
			}
			if !hasLimit {
				limit = kinesisstore.MaxTagPageResults
			}

			// Sort first so that ExclusiveStartTagKey filtering is lexical
			sort.Slice(tagList, func(i, j int) bool {
				return tagList[i].Key < tagList[j].Key
			})

			// Apply ExclusiveStartTagKey: keep only keys > startKey
			if startKey != "" {
				filtered := tagList[:0]
				for _, t := range tagList {
					if t.Key > startKey {
						filtered = append(filtered, t)
					}
				}
				tagList = filtered
			}

			hasMore := false
			if int32(len(tagList)) > int32(limit) {
				hasMore = true
				tagList = tagList[:limit]
			}

			return map[string]interface{}{
				"Tags":        tags.ToResponse(tagList),
				"HasMoreTags": hasMore,
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			switch err.(type) {
			case *tags.MissingResourceError:
				return ErrInvalidArgument
			case *tags.MissingTagsError:
				return ErrInvalidArgument
			case *tags.MissingTagKeysError:
				return ErrInvalidArgument
			}
			return s.mapStoreError(err)
		},
	}
}
