package server

import (
	"flag"

	"k8s.io/klog"
)

// Settings holds the application settings
type Settings struct {
	TillerAddress *string
	ListenPort    *int
	FrontendPath  *string
}

// SettingsFromCli reads the CLI options and constructs a Settings
// object from them.
func SettingsFromCli() *Settings {
	defaultTillerAddress := "tiller-deploy.kube-system.svc.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	listenPort := flag.Int("port", 8080, "Server listen port")
	frontendPath := flag.String("frontendPath", "web/dist", "Path to frontend files")

	klog.InitFlags(nil)
	flag.Parse()

	return &Settings{
		TillerAddress: tillerAddress,
		ListenPort:    listenPort,
		FrontendPath:  frontendPath,
	}
}
