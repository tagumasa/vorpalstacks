package s3

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"vorpalstacks/internal/core/logs"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// WebsiteServer serves static website content from S3 buckets
// configured with a website configuration.
type WebsiteServer struct {
	s3Store   S3StoreProvider
	accountID string
	region    string
}

// S3StoreProvider abstracts access to S3 bucket and object stores.
type S3StoreProvider interface {
	Buckets(region string) s3store.BucketStoreInterface
	Objects(region string) s3store.ObjectStoreInterface
}

// NewWebsiteServer creates a new WebsiteServer backed by the shared S3 store.
func NewWebsiteServer(s3Store S3StoreProvider, accountID, region string) *WebsiteServer {
	return &WebsiteServer{
		s3Store:   s3Store,
		accountID: accountID,
		region:    region,
	}
}

// HandleRequest serves static website content from an S3 bucket with website configuration enabled.
func (s *WebsiteServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	bucket := s.extractBucket(r.Host)
	if bucket == "" {
		http.Error(w, "NoSuchBucket: could not determine bucket from Host", http.StatusNotFound)
		return
	}

	var bucketStore s3store.BucketStoreInterface
	var objectStore s3store.ObjectStoreInterface
	bucketStore = s.s3Store.Buckets(s.region)
	objectStore = s.s3Store.Objects(s.region)
	if bucketStore == nil || objectStore == nil {
		http.Error(w, "InternalError: storage unavailable", http.StatusInternalServerError)
		return
	}

	b, err := bucketStore.Get(bucket)
	if err != nil || b == nil {
		s.writeWebsiteError(w, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist")
		return
	}

	wc := b.WebsiteConfiguration
	if wc == nil {
		s.writeWebsiteError(w, http.StatusNotFound, "NoSuchWebsiteConfiguration", "The specified bucket does not have a website configuration")
		return
	}

	if wc.RedirectAllRequestsTo != nil && wc.RedirectAllRequestsTo.HostName != "" {
		target := wc.RedirectAllRequestsTo.HostName
		if wc.RedirectAllRequestsTo.Protocol != "" {
			target = wc.RedirectAllRequestsTo.Protocol + "://" + target
		} else {
			target = "http://" + target
		}
		http.Redirect(w, r, target+r.URL.Path, http.StatusFound)
		return
	}

	requestKey := strings.TrimPrefix(r.URL.Path, "/")

	if requestKey == "" || requestKey == "/" {
		if wc.IndexDocument != "" {
			requestKey = wc.IndexDocument
		} else {
			requestKey = "index.html"
		}
	}

	if strings.HasSuffix(requestKey, "/") {
		if wc.IndexDocument != "" {
			requestKey = requestKey + wc.IndexDocument
		} else {
			requestKey = requestKey + "index.html"
		}
	}

	reader, obj, err := objectStore.Get(r.Context(), bucket, requestKey)
	if err == nil && reader != nil {
		defer reader.Close()

		ct := "application/octet-stream"
		if obj != nil && obj.ContentType != "" {
			ct = obj.ContentType
		}
		if obj != nil && obj.ETag != "" {
			w.Header().Set("ETag", obj.ETag)
		}
		if obj != nil && !obj.LastModified.IsZero() {
			w.Header().Set("Last-Modified", obj.LastModified.UTC().Format(http.TimeFormat))
		}

		w.Header().Set("Content-Type", ct)
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, reader); err != nil {
			logs.Warn("s3 website copy error", logs.String("bucket", bucket), logs.String("key", requestKey), logs.Err(err))
		}
		return
	}

	for _, rule := range wc.RoutingRules {
		if routingRuleMatches(&rule, requestKey, 404) {
			redirect := rule.Redirect
			if redirect.ReplaceKeyWith != nil {
				loc := "/" + *redirect.ReplaceKeyWith
				if redirect.HostName != nil {
					loc = "http://" + *redirect.HostName + loc
				}
				code := http.StatusFound
				if redirect.HTTPRedirectCode != nil {
					_, _ = fmt.Sscanf(*redirect.HTTPRedirectCode, "%d", &code)
				}
				http.Redirect(w, r, loc, code)
				return
			}
			if redirect.ReplaceKeyPrefixWith != nil {
				loc := "/" + *redirect.ReplaceKeyPrefixWith + strings.TrimPrefix(requestKey, *rule.Condition.KeyPrefixEquals)
				if redirect.HostName != nil {
					loc = "http://" + *redirect.HostName + loc
				}
				code := http.StatusFound
				if redirect.HTTPRedirectCode != nil {
					_, _ = fmt.Sscanf(*redirect.HTTPRedirectCode, "%d", &code)
				}
				http.Redirect(w, r, loc, code)
				return
			}
			if redirect.HostName != nil {
				proto := "http"
				if redirect.Protocol != nil {
					proto = *redirect.Protocol
				}
				loc := proto + "://" + *redirect.HostName + "/" + requestKey
				code := http.StatusFound
				if redirect.HTTPRedirectCode != nil {
					_, _ = fmt.Sscanf(*redirect.HTTPRedirectCode, "%d", &code)
				}
				http.Redirect(w, r, loc, code)
				return
			}
		}
	}

	if wc.ErrorDocument != "" {
		errReader, errObj, errGet := objectStore.Get(r.Context(), bucket, wc.ErrorDocument)
		if errGet == nil && errReader != nil {
			defer errReader.Close()
			ct := "text/html"
			if errObj != nil && errObj.ContentType != "" {
				ct = errObj.ContentType
			}
			w.Header().Set("Content-Type", ct)
			w.WriteHeader(http.StatusNotFound)
			if _, err := io.Copy(w, errReader); err != nil {
				logs.Warn("s3 website error document copy error", logs.String("bucket", bucket), logs.Err(err))
			}
			return
		}
	}

	s.writeWebsiteError(w, http.StatusNotFound, "NoSuchKey", fmt.Sprintf("The specified key %s does not exist", requestKey))
}

func (s *WebsiteServer) extractBucket(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	for i, part := range parts {
		if (strings.HasPrefix(part, "s3-website") || strings.HasPrefix(part, "s3-website-")) && i > 0 {
			return parts[0]
		}
	}
	if len(parts) >= 2 && parts[len(parts)-1] == "localhost" {
		if len(parts) >= 3 {
			candidate := strings.Join(parts[:len(parts)-2], ".")
			if candidate != "" {
				return candidate
			}
		}
		return parts[0]
	}
	return ""
}

func (s *WebsiteServer) writeWebsiteError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>%s</title></head>
<body>
<h1>%s</h1>
<p>%s</p>
</body>
</html>`, code, code, message)
}

func routingRuleMatches(rr *s3store.RoutingRule, key string, httpCode int) bool {
	if rr.Condition == nil {
		return false
	}
	if rr.Condition.HTTPErrorCodeReturnedEquals != nil {
		var code int
		_, _ = fmt.Sscanf(*rr.Condition.HTTPErrorCodeReturnedEquals, "%d", &code)
		if code != httpCode {
			return false
		}
	}
	if rr.Condition.KeyPrefixEquals != nil {
		if !strings.HasPrefix(key, *rr.Condition.KeyPrefixEquals) {
			return false
		}
	}
	return true
}
