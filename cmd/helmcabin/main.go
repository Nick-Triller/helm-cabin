package main

import (
	"github.com/Nick-Triller/helm-cabin/internal/server"
	"github.com/Nick-Triller/helm-cabin/internal/settings"
)

func main() {
	server.NewServer(settings.FromCli()).Start()
}
