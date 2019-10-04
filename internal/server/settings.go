package server

import "flag"

type Settings struct {
	TillerAddress *string
	ListenPort    *int
	FrontendPath  *string
}

func SettingsFromCli() *Settings {
	defaultTillerAddress := "tiller-deploy.kube-system.svc.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	listenPort := flag.Int("port", 8080, "Server listen port")
	frontendPath := flag.String("frontendPath", "web/dist", "Path to frontend files")
	flag.Parse()
	return &Settings{
		TillerAddress: tillerAddress,
		ListenPort:    listenPort,
		FrontendPath:  frontendPath,
	}
}
