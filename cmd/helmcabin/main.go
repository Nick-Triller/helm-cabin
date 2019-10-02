package main

import (
	"github.com/Nick-Triller/helmcabin/internal/server"
)

func main() {
	server.NewServer().Start()
}
