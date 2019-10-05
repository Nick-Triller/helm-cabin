package server

import (
	"fmt"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	"log"
	"net/http"
	"sync"
)

type server struct {
	releasesChan        chan *services.ListReleasesResponse
	releasesCache       *services.ListReleasesResponse
	releasesCacheMutex  sync.RWMutex
	tillerReachableChan chan bool
	settings            *Settings
}

// NewServer creates a server struct
func NewServer(settings *Settings) *server {
	return &server{settings: settings}
}

// getReleases returns a list of cached releasesCache
func (s *server) getReleases() *services.ListReleasesResponse {
	s.releasesCacheMutex.RLock()
	releases := s.releasesCache
	defer s.releasesCacheMutex.RUnlock()
	return releases
}

// Start is the main entrypoint that bootstraps the application
func (s *server) Start() {
	log.Printf("helm client version: %s\n", version.GetVersion())

	s.releasesChan = make(chan *services.ListReleasesResponse)
	s.tillerReachableChan = make(chan bool)

	// Drain unused channel
	go func() { for { <- s.tillerReachableChan }}()
	go PollReleases(s.releasesChan, s.tillerReachableChan, s.settings)
	go cacheReleases(s)

	log.Printf("Starting server on port %d ", *s.settings.ListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *s.settings.ListenPort), router(s)))
}

// cacheReleases caches polled releasesCache
func cacheReleases(s *server) {
	for {
		releases := <-s.releasesChan
		s.releasesCacheMutex.Lock()
		s.releasesCache = releases
		s.releasesCacheMutex.Unlock()
	}
}
