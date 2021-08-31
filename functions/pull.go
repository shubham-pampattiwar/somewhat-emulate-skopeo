package functions

import (
	"context"
	"fmt"
	"github.com/containers/image/copy"
	"github.com/containers/image/signature"
	"github.com/containers/image/storage"
	"github.com/containers/image/transports/alltransports"
	"github.com/containers/image/types"
	store "github.com/containers/storage"
	"github.com/sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"os"
)
var _defaultStore store.Store
//var _storeOptions store.StoreOptions

func ImagePull(ImageName string) {
	fmt.Printf("Pulling image %v \n", ImageName)

	sourceImageRef, err := alltransports.ParseImageName(ImageName)
	if err != nil {
		logrus.WithError(err).Fatal("image name parsing error")
	}

	systemContext := &types.SystemContext{}
	policy, err := signature.DefaultPolicy(systemContext)
	if err != nil {
		logrus.WithError(err).Fatal("policy creation error")
	}

	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		logrus.WithError(err).Fatal("policy context creation error")
	}

	destinationImageName := ImageName
	if sourceImageRef.DockerReference() != nil {
		destinationImageName = sourceImageRef.DockerReference().Name()
	}

	destinationImageRef, err := storage.Transport.ParseStoreReference(defaultStore(), destinationImageName)
	if err != nil {
		logrus.WithError(err).Fatal("Could not parse local image reference")
	}

	manifest, err := copy.Image(
		context.Background(),
		policyContext,
		destinationImageRef,
		sourceImageRef,
		&copy.Options{
			ReportWriter: os.Stdout,
		},
		)
	if err != nil {
		logrus.WithError(err).Fatal("unable to pull image")
	}

	fmt.Printf("Pulled Image manifest %v \n", string(manifest))
	images, _ := _defaultStore.Images()
	image, _ := _defaultStore.Image(images[0].ID)
	spew.Dump(image)
}

func defaultStore() store.Store {
	options, err := store.DefaultStoreOptions(false, 0)
	fmt.Printf("\n filesystem path is %v \n,", options.GraphRoot)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create default image store options")
	}
	//options.RunRoot = "/run/containers/storage"
	//options.GraphRoot = "/var/lib/containers/storage"
	//options.GraphDriverName = "overlay"

	if _defaultStore == nil {
		gotStorage, err := store.GetStore(options)
		if err != nil {
			logrus.WithError(err).Fatal("Could not create image store")
		}
		_defaultStore = gotStorage
	}

	return _defaultStore
}

//func InitDefaultStoreOptions() {
//	options, err := store.DefaultStoreOptions(false, 0)
//	if err != nil {
//		logrus.WithError(err).Fatal("Could not create default image store options")
//	}
//	options.RunRoot = "/run/containers/storage"
//	options.GraphRoot = "/var/lib/containers/storage"
//	options.GraphDriverName = "overlay"
//	_storeOptions = options
//
//}