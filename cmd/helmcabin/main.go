package main

import (
	"github.com/Nick-Triller/helm-cabin/internal/server"
)

func main() {
	server.NewServer(server.SettingsFromCli()).Start()
}
