package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

// instance references the server struct for which the router handler was created
var instance *Server

// router is the main handler that bundles all other handlers.
func router(s *Server) http.Handler {
	instance = s
	r := mux.NewRouter()
	// API routes
	r.HandleFunc("/api/releases", releasesHandler)
	r.HandleFunc("/api/releases/{name}/versions/{version}", releaseHandler)
	r.HandleFunc("/api/releases/{name}/revisions/", revisionsHandler)
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

// releasesHandler
func releasesHandler(w http.ResponseWriter, r *http.Request) {
	//qparam := r.URL.Query()
	//statusCodes := qparam["statusCodes"]
	//limit := atoiOrDefault(qparam.Get("limit"), 9999999)
	//offset := qparam.Get("offset")
	//sortBy := atoiOrDefault(qparam.Get("sortBy"), int(services.ListSort_LAST_RELEASED))
	//filter := qparam.Get("filter")
	//sortOrder := atoiOrDefault(qparam.Get("sortOrder"), int(services.ListSort_DESC))
	//namespace := qparam.Get("namespace") // All namespace if empty

	//listReleasesRequest := &services.ListReleasesRequest{
	//	Limit:       int64(limit),
	//	Offset:      offset,
	//	SortBy:      services.ListSort_SortBy(sortBy),
	//	Filter:      filter,
	//	SortOrder:   services.ListSort_SortOrder(sortOrder),
	//	StatusCodes: toStatusCodes(statusCodes),
	//	Namespace:   namespace,
	//}
	// resp, err := getReleasesList(listReleasesRequest)

	resp := instance.getCachedReleases()
	// if err != nil {
	//	w.WriteHeader(500)
	//	return
	// }

	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	_, _ = w.Write(jsonData)
	return
}

// releaseHandler
func releaseHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	releaseName := routeVars["name"]
	releaseVersion := routeVars["version"]
	releases := instance.getCachedReleases()
	release := findRelease(releases, releaseName, releaseVersion)

	if release == nil {
		w.WriteHeader(404)
		return
	}

	jsonData, _ := json.MarshalIndent(release, "", "  ")
	_, _ = w.Write(jsonData)
}

// revisionsHandler
func revisionsHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	releaseName := routeVars["name"]
	revisions := findRevisions(instance.getCachedReleases(), releaseName)
	versions := make([]int32, 0)
	for _, revision := range revisions {
		versions = append(versions, revision.Version)
	}
	jsonData, _ := json.MarshalIndent(versions, "", "  ")
	_, _ = w.Write(jsonData)
}
