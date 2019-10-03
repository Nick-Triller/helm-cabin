package main

import (
	"flag"
	"github.com/Nick-Triller/helmcabin/internal/server"
)

func main() {
	defaultTillerAddress := "tiller-deploy.svc.kube-system.cluster.local:44134"
	tillerAddress := flag.String("tillerAddress", defaultTillerAddress, "Tiller address")
	flag.Parse()
	server.NewServer(tillerAddress).Start()
}
