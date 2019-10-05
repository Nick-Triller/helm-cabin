package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strconv"
)

// instance references the server struct for which the router handler was created
var instance *server

// router is the main handler that bundles all other handlers.
func router(s *server) http.Handler {
	instance = s
	r := mux.NewRouter()
	// API routes
	r.HandleFunc("/api/releasesCache", releasesHandler)
	// Serve frontend
	frontendDir := http.Dir(*s.settings.FrontendPath)
	// Serve index.html if nothing matched
	r.PathPrefix("/").Handler(fileServerWithCustom404(frontendDir))
	return r
}

// fileServerWithCustom404 serves frontend files.
// If no file was found, it always serves the single page application.
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

// notFoundHandler always serves the single page application
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(*instance.settings.FrontendPath, "index.html"))
}

// releaseHandler
func releasesHandler(w http.ResponseWriter, r *http.Request) {
	filterStatuses := r.URL.Query()["status"]
	resources := make([]*releaseResource, 0)
	for _, r := range instance.getCachedReleases().GetReleases() {
		// Filter according to status
		releaseStatus := strconv.FormatInt(int64(r.Info.Status.Code), 10)
		if contains(filterStatuses, releaseStatus) {
			resources = append(resources, releaseToResource(r))
		}
	}
	jsonData, _ := json.MarshalIndent(resources, "", "  ")
	_, _ = w.Write(jsonData)
}
