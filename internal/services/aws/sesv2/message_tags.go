package sesv2

import "vorpalstacks/internal/common/request"

// MessageTag represents a tag with a name and value for message routing
// (SendEmail EmailTags / DefaultEmailTags).
type MessageTag struct {
	Name  string
	Value string
}

// ParseMessageTags parses message tags from request parameters.
func ParseMessageTags(params map[string]interface{}, listKey string) []MessageTag {
	tagsIf, ok := params[listKey]
	if !ok {
		return nil
	}
	tagList, ok := tagsIf.([]interface{})
	if !ok {
		return nil
	}
	return ParseMessageTagsFromList(tagList)
}

// ParseMessageTagsFromList parses message tags from a list of interface values.
func ParseMessageTagsFromList(tagList []interface{}) []MessageTag {
	var result []MessageTag
	for _, t := range tagList {
		tagMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, MessageTag{
			Name:  request.GetStringParam(tagMap, "Name"),
			Value: request.GetStringParam(tagMap, "Value"),
		})
	}
	return result
}
