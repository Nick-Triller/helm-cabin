package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strconv"
)

func router(s *server) http.Handler {
	r := mux.NewRouter()
	// API routes
	r.HandleFunc("/api/releases", releasesHandler)
	// Serve frontend
	frontendDir := http.Dir(*s.settings.FrontendPath)
	// Serve index.html if nothing matched
	r.PathPrefix("/").Handler(fileServerWithCustom404(frontendDir))
	return r
}

func fileServerWithCustom404(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			notFoundHandler(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(*instance.settings.FrontendPath, "index.html"))
}

func releasesHandler(w http.ResponseWriter, r *http.Request) {
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
