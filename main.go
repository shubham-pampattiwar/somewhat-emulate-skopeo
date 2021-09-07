package main

import (
	"github.com/containers/storage/pkg/reexec"
	"somewhat-emulate-skopeo/main/functions"
)

func main() {
	if reexec.Init() {
		return
	}
	// initializing storage
	functions.InitDefaultStoreOptions()
	functions.Show()
	functions.ImagePull("docker://alpine:latest")
	//functions.InitDefaultStoreOptions()
	//functions.ClearStuff()
	functions.Show()
}
