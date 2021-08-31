package main

import (
	"fmt"
	"somewhat-emulate-skopeo/main/functions"
	"github.com/containers/storage/pkg/reexec"
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
