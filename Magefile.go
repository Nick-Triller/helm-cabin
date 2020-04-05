//+build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// BuildDockerImage builds all artifacts in a docker multi stage build
// and results in a Docker Image containing the artifacts
func BuildDockerImage() error {
	return sh.RunV("docker", "build", "-f", "build/Dockerfile", "-t", "nicktriller/helm-cabin:latest", ".")
}

// BuildBackend locally runs the backend with `go run`.
// This target is supposed to be run in development only and assumes tiller is
// reachable at 127.0.0.1:44134
func RunServer() error {
	helmVersion := "3"
	if value, ok := os.LookupEnv("CABIN_HELM_VERSION"); ok {
		helmVersion = value
	}
	return sh.RunV("go", "run", "cmd/helmcabin/main.go", "--tillerAddress", "127.0.0.1:44134",
		"--listenAddress", "localhost:8080", "--helmVersion", helmVersion)
}

// BuildServerAll locally builds the artifacts for all supported platforms
func BuildServerAll() error {
	if err := BuildServerWindows(); err != nil {
		return err
	}
	return BuildServerLinux()
}

// BuildWindows locally builds the artifacts for windows
func BuildServerWindows() error {
	return buildLocal("windows", "amd64")
}

// BuildWindows locally builds the artifacts for linux
func BuildServerLinux() error {
	return buildLocal("linux", "amd64")
}

// Lint lints the go code
func Lint() error {
	return sh.RunV("golint", "./...")
}

func buildLocal(goos string, goarch string) error {
	mg.Deps(Lint)
	_ = os.Setenv("GOOS", goos)
	_ = os.Setenv("GOARCH", goarch)
	out := "bin/helmcabin"
	if goos == "windows" {
		out += ".exe"
	}
	return sh.RunV("go", "build", "-o", out, "./cmd/helmcabin")
}
