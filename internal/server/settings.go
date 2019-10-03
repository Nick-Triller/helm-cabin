package server

import "flag"

type Settings struct {
	TillerAddress *string
	ListenPort *int
}

func SettingsFromCli() *Settings {
	defaultTillerAddress := "tiller-deploy.svc.kube-system.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	listenPort := flag.Int("port", 8080, "Server listen port")
	flag.Parse()
	return &Settings{
		TillerAddress: tillerAddress,
		ListenPort:    listenPort,
	}
}