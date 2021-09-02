package main

import (
	"fmt"
	"github.com/containers/storage/pkg/reexec"
	"somewhat-emulate-skopeo/main/functions"
)

func main() {
	fmt.Printf("Hello World")
	if reexec.Init() {
		return
	}
	// initializing storage
	//functions.InitDefaultStoreOptions()
	functions.ImagePull("docker://alpine:latest")
}
