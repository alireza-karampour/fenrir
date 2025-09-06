package main

import (
	"github.com/alireza-karampour/fenrir/cmd"
	_ "github.com/alireza-karampour/fenrir/cmd/test"
	"github.com/containers/buildah"
)

func main() {
	if buildah.InitReexec() {
		return
	}
	cmd.Execute()
}
