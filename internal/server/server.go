package server

import (
	settings2 "github.com/Nick-Triller/helm-cabin/internal/settings"
	"net/http"
	"sync"

	"github.com/Nick-Triller/helm-cabin/internal/helm2"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	log "k8s.io/klog"
)

// Server is the main application struct
type Server struct {
	releasesChan       chan *services.ListReleasesResponse
	releasesCache      *services.ListReleasesResponse
	releasesCacheMutex sync.RWMutex
	settings           *settings2.Settings
}

// NewServer creates a server struct
func NewServer(settings *settings2.Settings) *Server {
	return &Server{settings: settings}
}

// getCachedReleases returns a list of cached releasesCache
func (s *Server) getCachedReleases() *services.ListReleasesResponse {
	s.releasesCacheMutex.RLock()
	releases := s.releasesCache
	defer s.releasesCacheMutex.RUnlock()
	return releases
}

// Start is the main entrypoint that bootstraps the application
func (s *Server) Start() {
	log.Infof("helm client version: %s\n", version.GetVersion())
	s.releasesChan = make(chan *services.ListReleasesResponse)

	go helm2.PollReleases(s.releasesChan, s.settings)
	go cacheReleases(s)

	log.Infof("Starting server on %s ", *s.settings.ListenAddress)
	log.Fatal(http.ListenAndServe(*s.settings.ListenAddress, router(s)))
}

// cacheReleases caches polled releasesCache
func cacheReleases(s *Server) {
	for {
		releases := <-s.releasesChan
		s.releasesCacheMutex.Lock()
		s.releasesCache = releases
		s.releasesCacheMutex.Unlock()
	}
}
