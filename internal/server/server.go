package server

import (
	"github.com/Nick-Triller/helm-cabin/internal/helm2"
	"github.com/Nick-Triller/helm-cabin/internal/helm3"
	"github.com/Nick-Triller/helm-cabin/internal/resources"
	"github.com/Nick-Triller/helm-cabin/internal/settings"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/version"
)

// Server is the main application struct
type Server struct {
	releasesChan       chan []resources.ReleaseResource
	releasesCache      []resources.ReleaseResource
	releasesCacheMutex sync.RWMutex
	settings           *settings.Settings
}

// NewServer creates a server struct
func NewServer(settings *settings.Settings) *Server {
	return &Server{
		settings: settings,
		releasesChan: make(chan []resources.ReleaseResource),
	}
}

// getCachedReleases returns a list of cached releasesCache
func (s *Server) getCachedReleases() []resources.ReleaseResource {
	s.releasesCacheMutex.RLock()
	releases := s.releasesCache
	defer s.releasesCacheMutex.RUnlock()
	return releases
}

// Start is the main entrypoint that bootstraps the application
func (s *Server) Start() {
	helmVersion := 3

	switch helmVersion {
	case 2:
		log.Infof("helm client version: %s", version.GetVersion())
		go helm2.PollReleases(s.releasesChan, s.settings)
	case 3:
		go helm3.PollReleases(s.releasesChan, s.settings)
	default:
		log.Fatalf("Unknown helm version: %d", helmVersion)
	}
	log.Infof("Running in helm %d mode", helmVersion)

	go cacheReleases(s)

	log.Infof("Starting server on %s", *s.settings.ListenAddress)
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
