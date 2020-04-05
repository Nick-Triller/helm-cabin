package settings

import (
	"flag"

	"k8s.io/klog"
)

// Settings holds the application settings
type Settings struct {
	TillerAddress *string
	ListenAddress *string
	FrontendPath  *string
}

// FromCli reads the CLI options and constructs a Settings
// object from them.
func FromCli() *Settings {
	defaultTillerAddress := "tiller-deploy.kube-system.svc.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	listenAddress := flag.String("listenAddress", ":8080", "Server listen address")
	frontendPath := flag.String("frontendPath", "web/dist", "Path to frontend files")

	klog.InitFlags(nil)
	flag.Parse()

	return &Settings{
		TillerAddress: tillerAddress,
		ListenAddress: listenAddress,
		FrontendPath:  frontendPath,
	}
}
