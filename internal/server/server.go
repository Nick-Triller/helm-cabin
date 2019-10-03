package server

import (
	"encoding/json"
	"fmt"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	releasesChan        chan *services.ListReleasesResponse
	releases            *services.ListReleasesResponse
	tillerReachableChan chan bool
	tillerReachable     bool
	settings            *Settings
}

// NewServer creates a server struct
var instance *server

func NewServer(settings *Settings) *server {
	instance = &server{settings: settings}
	return instance
}

// Start is the main entrypoint that bootstraps the application
func (s *server) Start() {
	log.Printf("helm client version: %s\n", version.GetVersion())

	s.releasesChan = make(chan *services.ListReleasesResponse)
	s.tillerReachableChan = make(chan bool)
	go PollReleases(s.releasesChan, s.tillerReachableChan, s.settings)
	go watchChannels(s)

	http.HandleFunc("/releases", releasesEndpoint)

	log.Printf("Starting server on port %d\n ", *s.settings.ListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *s.settings.ListenPort), nil))
}

func releasesEndpoint(w http.ResponseWriter, r *http.Request) {
	filterStatuses := r.URL.Query()["status"]
	resources := make([]*releaseResource, 0)
	for _, r := range instance.releases.GetReleases() {
		// Filter according to status
		releaseStatus := strconv.FormatInt(int64(r.Info.Status.Code), 10)
		if contains(filterStatuses, releaseStatus) {
			resources = append(resources, releaseToResource(r))
		}
	}
	jsonData, _ := json.MarshalIndent(resources, "", "  ")
	_, _ = w.Write(jsonData)
}

func watchChannels(s *server) {
	for {
		select {
		case s.releases = <-s.releasesChan:
		case s.tillerReachable = <-s.tillerReachableChan:
		}
	}
}
