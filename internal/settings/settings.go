package settings

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

// Settings holds the application settings
type Settings struct {
	TillerAddress *string
	ListenAddress *string
	FrontendPath  *string
	HelmVersion   int
}

// FromCli reads the CLI options and constructs a Settings
// object from them.
func FromCli() *Settings {
	log.SetFormatter(&log.TextFormatter{})

	defaultTillerAddress := "tiller-deploy.kube-system.svc.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	listenAddress := flag.String("listenAddress", ":8080", "Server listen address")
	frontendPath := flag.String("frontendPath", "web/dist", "Path to frontend files")
	helmVersion := flag.Int("helmVersion", 3, "Show releases of helm 2 or 3")

	flag.Parse()

	return &Settings{
		TillerAddress: tillerAddress,
		ListenAddress: listenAddress,
		FrontendPath:  frontendPath,
		HelmVersion: *helmVersion,
	}
}
