package request

import (
	"net/http"
	"strings"
)

func extractSESv2Operation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !strings.HasPrefix(path, "/v2/email/") {
		return ""
	}

	path = strings.TrimPrefix(path, "/v2/email/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		return ""
	}

	resource := parts[0]

	switch resource {
	case "configuration-sets":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListConfigurationSets"
			case http.MethodPost:
				return "CreateConfigurationSet"
			}
		} else if len(parts) >= 2 {
			switch method {
			case http.MethodGet:
				return "GetConfigurationSet"
			case http.MethodDelete:
				return "DeleteConfigurationSet"
			}
			if len(parts) >= 3 && parts[2] == "event-destinations" {
				switch method {
				case http.MethodGet:
					return "GetConfigurationSetEventDestinations"
				}
				if len(parts) >= 4 {
					switch method {
					case http.MethodPut:
						return "UpdateConfigurationSetEventDestination"
					case http.MethodDelete:
						return "DeleteConfigurationSetEventDestination"
					}
				} else {
					switch method {
					case http.MethodPost:
						return "CreateConfigurationSetEventDestination"
					}
				}
			}
		}

	case "identities":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListEmailIdentities"
			case http.MethodPost:
				return "CreateEmailIdentity"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 {
				switch parts[2] {
				case "dkim":
					switch method {
					case http.MethodPut:
						if len(parts) >= 4 && parts[3] == "signing" {
							return "PutEmailIdentityDkimSigningAttributes"
						}
						return "PutEmailIdentityDkimAttributes"
					}
				case "policies":
					if len(parts) >= 4 {
						switch method {
						case http.MethodGet:
							return "GetEmailIdentityPolicies"
						case http.MethodPut:
							return "UpdateEmailIdentityPolicy"
						case http.MethodDelete:
							return "DeleteEmailIdentityPolicy"
						}
					} else {
						switch method {
						case http.MethodGet:
							return "GetEmailIdentityPolicies"
						case http.MethodPost:
							return "CreateEmailIdentityPolicy"
						}
					}
				case "feedback":
					if method == http.MethodPut {
						return "PutEmailIdentityFeedbackAttributes"
					}
				case "mail-from":
					if method == http.MethodPut {
						return "PutEmailIdentityMailFromAttributes"
					}
				}
			}
			switch method {
			case http.MethodGet:
				return "GetEmailIdentity"
			case http.MethodDelete:
				return "DeleteEmailIdentity"
			}
		}

	case "templates":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListEmailTemplates"
			case http.MethodPost:
				return "CreateEmailTemplate"
			}
		} else if len(parts) >= 2 {
			switch method {
			case http.MethodGet:
				return "GetEmailTemplate"
			case http.MethodPut:
				return "UpdateEmailTemplate"
			case http.MethodDelete:
				return "DeleteEmailTemplate"
			}
		}

	case "outbound-emails":
		if method == http.MethodPost {
			return "SendEmail"
		}

	case "account":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "GetAccount"
			case http.MethodPut:
				return "PutAccountAttributes"
			}
		} else if len(parts) >= 2 {
			switch parts[1] {
			case "sending":
				if method == http.MethodPut {
					return "PutAccountSendingAttributes"
				}
			case "suppression":
				if method == http.MethodPut {
					return "PutAccountSuppressionAttributes"
				}
			case "details":
				if method == http.MethodPut {
					return "PutAccountDetails"
				}
			}
		}

	case "tags":
		switch method {
		case http.MethodPost:
			return "TagResource"
		case http.MethodDelete:
			return "UntagResource"
		case http.MethodGet:
			return "ListTagsForResource"
		}

	case "contact-lists":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListContactLists"
			case http.MethodPost:
				return "CreateContactList"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "contacts" {
				if len(parts) >= 4 {
					switch method {
					case http.MethodGet:
						return "GetContact"
					case http.MethodPut:
						return "UpdateContact"
					case http.MethodDelete:
						return "DeleteContact"
					}
				} else {
					switch method {
					case http.MethodGet:
						return "ListContacts"
					case http.MethodPost:
						return "CreateContact"
					}
				}
			} else {
				switch method {
				case http.MethodGet:
					return "GetContactList"
				case http.MethodPut:
					return "UpdateContactList"
				case http.MethodDelete:
					return "DeleteContactList"
				}
			}
		}

	case "suppression":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListSuppressedDestinations"
			case http.MethodPost:
				return "PutSuppressedDestination"
			}
		} else if len(parts) >= 2 {
			switch method {
			case http.MethodGet:
				return "GetSuppressedDestination"
			case http.MethodDelete:
				return "DeleteSuppressedDestination"
			}
		}

	case "dedicated-ip-pools":
		if len(parts) == 1 {
			switch method {
			case http.MethodGet:
				return "ListDedicatedIpPools"
			case http.MethodPost:
				return "CreateDedicatedIpPool"
			}
		} else if len(parts) >= 2 {
			switch method {
			case http.MethodGet:
				return "GetDedicatedIpPool"
			case http.MethodDelete:
				return "DeleteDedicatedIpPool"
			}
		}
	}

	return ""
}

func extractSESv2PathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/v2/email/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[0] == "" {
		return
	}

	resource := parts[0]

	switch resource {
	case "configuration-sets":
		if len(parts) >= 2 && parts[1] != "" {
			params["ConfigurationSetName"] = parts[1]
		}
	case "identities":
		if len(parts) >= 2 && parts[1] != "" {
			params["EmailIdentity"] = parts[1]
			if len(parts) >= 4 && parts[2] == "policies" {
				params["PolicyName"] = parts[3]
			}
		}
	case "templates":
		if len(parts) >= 2 && parts[1] != "" {
			params["TemplateName"] = parts[1]
		}
	case "contact-lists":
		if len(parts) >= 2 && parts[1] != "" {
			params["ContactListName"] = parts[1]
			if len(parts) >= 4 && parts[2] == "contacts" {
				params["EmailAddress"] = parts[3]
			}
		}
	case "suppression":
		if len(parts) >= 2 && parts[1] != "" {
			params["EmailAddress"] = parts[1]
		}
	case "dedicated-ip-pools":
		if len(parts) >= 2 && parts[1] != "" {
			params["PoolName"] = parts[1]
		}
	case "tags":
		if len(parts) >= 2 && parts[1] != "" {
			params["ResourceArn"] = parts[1]
		}
	}
}
