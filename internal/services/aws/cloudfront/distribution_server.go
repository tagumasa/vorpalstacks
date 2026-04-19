package cloudfront

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	cfstore "vorpalstacks/internal/store/aws/cloudfront"
)

// DistributionServer serves requests for CloudFront distributions by proxying
// to the configured origin and applying cache behaviour headers.
type DistributionServer struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	region         string
	client         *http.Client
	distributionMu sync.RWMutex
	distribution   *cfstore.DistributionStore
}

// NewDistributionServer creates a new DistributionServer with the given storage, account, and region.
func NewDistributionServer(storageManager *storage.RegionStorageManager, accountID, region string) *DistributionServer {
	return &DistributionServer{
		storageManager: storageManager,
		accountID:      accountID,
		region:         region,
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// SetDistributionStore injects a DistributionStore instance, typically the same
// one used by CloudFrontService, so that both components share a single store.
func (s *DistributionServer) SetDistributionStore(store *cfstore.DistributionStore) {
	s.distributionMu.Lock()
	s.distribution = store
	s.distributionMu.Unlock()
}

// HandleRequest proxies an incoming request to the matching CloudFront distribution's origin.
func (s *DistributionServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	distributionID := s.extractDistributionID(r.Host)
	if distributionID == "" {
		http.Error(w, `{"message":"Distribution not found"}`, http.StatusNotFound)
		return
	}

	store := s.getDistributionStore()
	dist, err := store.Get(distributionID)
	if err != nil || dist == nil {
		http.Error(w, `{"message":"Distribution not found"}`, http.StatusNotFound)
		return
	}

	if !dist.Enabled {
		http.Error(w, `{"message":"Distribution is disabled"}`, http.StatusForbidden)
		return
	}

	config := dist.DistributionConfig
	if config == nil {
		http.Error(w, `{"message":"Distribution has no configuration"}`, http.StatusInternalServerError)
		return
	}

	requestPath := r.URL.Path
	if requestPath == "" || requestPath == "/" {
		if config.DefaultRootObject != "" {
			requestPath = "/" + config.DefaultRootObject
		}
	}

	behavior := s.matchCacheBehavior(config, requestPath)
	if behavior == nil {
		http.Error(w, `{"message":"No matching cache behavior"}`, http.StatusNotFound)
		return
	}

	origin := s.resolveOrigin(config.Origins, behavior.TargetOriginId)
	if origin == nil {
		http.Error(w, `{"message":"Origin not found"}`, http.StatusNotFound)
		return
	}

	targetURL := s.buildOriginURL(origin, requestPath, r.URL.RawQuery)

	originReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, `{"message":"Failed to create origin request"}`, http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			originReq.Header.Add(key, value)
		}
	}
	originReq.Header.Set("X-Forwarded-Host", r.Host)
	originReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	resp, err := s.client.Do(originReq)
	if err != nil {
		logs.Error("CloudFront origin request failed", logs.String("url", targetURL), logs.Err(err))
		http.Error(w, `{"message":"Origin request failed"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("X-Cache", "Miss from cloudfront")
	w.Header().Set("Via", "1.1 "+distributionID+".cloudfront.net (CloudFront)")
	w.Header().Del("Transfer-Encoding")

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		logs.Warn("cloudfront distribution proxy copy error", logs.String("distributionId", distributionID), logs.Err(err))
	}
}

func (s *DistributionServer) extractDistributionID(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) >= 1 && parts[0] != "" && parts[0] != "localhost" {
		return parts[0]
	}
	return ""
}

func (s *DistributionServer) matchCacheBehavior(config *cfstore.DistributionConfig, path string) *cfstore.CacheBehavior {
	if config.CacheBehaviors != nil {
		for _, behavior := range config.CacheBehaviors.Items {
			if s.pathMatches(behavior.PathPattern, path) {
				return behavior
			}
		}
	}
	return config.DefaultCacheBehavior
}

func (s *DistributionServer) pathMatches(pattern, path string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return path == pattern
}

func (s *DistributionServer) resolveOrigin(origins cfstore.Origins, targetOriginID string) *cfstore.Origin {
	for _, origin := range origins.Items {
		if origin.ID == targetOriginID {
			return origin
		}
	}
	return nil
}

func (s *DistributionServer) buildOriginURL(origin *cfstore.Origin, path, query string) string {
	scheme := "http"
	if origin.CustomOriginConfig != nil && origin.CustomOriginConfig.OriginProtocolPolicy == "https-only" {
		scheme = "https"
	}

	host := origin.DomainName
	if strings.Contains(host, "s3-website") || strings.Contains(host, ".s3-website.") {
		host = fmt.Sprintf("localhost:%d", config.GetInt("ports.s3_website"))
	} else if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
	} else if strings.Contains(host, "lambda-url") {
		host = fmt.Sprintf("localhost:%d", config.GetInt("ports.lambda_url"))
	} else if strings.Contains(host, "s3.") && !strings.Contains(host, "website") {
		host = fmt.Sprintf("localhost:%d", config.ServerPort())
	}

	originPath := strings.TrimSuffix(origin.OriginPath, "/")
	url := fmt.Sprintf("%s://%s%s%s", scheme, host, originPath, path)
	if query != "" {
		url += "?" + query
	}
	return url
}

func (s *DistributionServer) getDistributionStore() *cfstore.DistributionStore {
	s.distributionMu.RLock()
	if s.distribution != nil {
		cached := s.distribution
		s.distributionMu.RUnlock()
		return cached
	}
	s.distributionMu.RUnlock()

	s.distributionMu.Lock()
	defer s.distributionMu.Unlock()
	if s.distribution != nil {
		return s.distribution
	}
	store, err := s.storageManager.GetStorage(s.region)
	if err != nil {
		return cfstore.NewDistributionStore(nil, s.accountID)
	}
	s.distribution = cfstore.NewDistributionStore(store, s.accountID)
	return s.distribution
}
