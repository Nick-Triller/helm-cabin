package server

import (
	"encoding/json"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/version"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	releasesChan chan *services.ListReleasesResponse
	releases *services.ListReleasesResponse
	tillerAddress *string
	tillerReachableChan chan bool
	tillerReachable bool
}

// NewServer creates a server struct
var serverInstance *server
func NewServer(tillerAddress *string) *server {
	serverInstance = &server{tillerAddress:tillerAddress}
	return serverInstance
}


// Start is the main entrypoint that bootstraps the application
func (s *server) Start() {
	log.Printf("helm client version: %s\n", version.GetVersion())

	s.releasesChan = make(chan *services.ListReleasesResponse)
	s.tillerReachableChan = make(chan bool)
	go PollReleases(s.releasesChan, s.tillerReachableChan, s.tillerAddress)
	go watchChannels(s)

	http.HandleFunc("/releases", releasesEndpoint)

	addr := ":" + strconv.Itoa(8080)
	log.Println("Starting server with listen address " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func releasesEndpoint(w http.ResponseWriter, r *http.Request) {
	filterStatuses := r.URL.Query()["status"]
	resources  := make([]*releaseResource, 0)
	for _, r := range serverInstance.releases.GetReleases() {
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
		case s.releases = <- s.releasesChan:
		case s.tillerReachable = <- s.tillerReachableChan:
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
