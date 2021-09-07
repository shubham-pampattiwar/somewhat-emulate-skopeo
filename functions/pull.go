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
	archieve "github.com/containers/storage/pkg/archive"
	"github.com/sirupsen/logrus"
	"io"
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

	fmt.Printf("\n Pulled Image Src refs: %+v \n", sourceImageRef)

	manifest, err := copy.Image(
		context.TODO(),
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
	// fetch image
	images, _ := _defaultStore.Images()
	image, _ := _defaultStore.Image(images[0].ID)
	// create container
	container, _ := _defaultStore.CreateContainer("", nil, image.ID, "", "", nil)

	// mount container
	mountPoint, _ := _defaultStore.Mount(container.ID, "")
	fmt.Printf("Container mounted at %v \n", mountPoint)

	// convert the mounted container image to tarball
	exportedImageTar, _ := archieve.Tar(mountPoint, archieve.Uncompressed)
	if err != nil {
		logrus.WithError(err).Fatal("error creating tar")
	}
	fmt.Printf("exported tar: %v \n", exportedImageTar)
	//homeDir, _ := os.UserHomeDir()
	outputFile, err := os.Create( "/home/shubham/foo-tar-new-8")
	if err != nil {
		logrus.WithError(err).Fatal("unable to create tar output file")
	}

	_, err = io.Copy(outputFile, exportedImageTar)

	if err != nil {
		logrus.WithError(err).Fatal("unable to create tarball")
	}


	// Importing an image from tarball... Lets go
	// convert the tarball to image and compare images
	fmt.Printf("converting tarball to image from tar: %v \n", outputFile.Name())
	//reader, err := os.Open(outputFile.Name())
	//if err != nil {
	//	logrus.WithError(err).Fatal("unable to open the tar ball file")
	//}
	//err = Untar("/home/shubham/foo-tar-new-8", "/home/shubham/test-untar/")

	//err = archieve.Untar(reader, "/home/shubham/foo-tar-new-8", nil)
	//
	//if err != nil {
	//	logrus.WithError(err).Fatal("unable to untar the tarball")
	//}

	src, err := alltransports.ParseImageName("tarball://" + outputFile.Name())
	//src, err := tarball.Transport.ParseReference("/home/shubham/foo-tar-new-8")
	if err != nil {
		logrus.WithError(err).Fatal("tarball name parsing error")
	}

	if src.DockerReference() != nil {
		destinationImageName = src.DockerReference().Name()
	}



	dest, err := storage.Transport.ParseStoreReference(defaultStore(), destinationImageName)
	if err != nil {
		logrus.WithError(err).Fatal("error parsing dest image [tar 2 image]")
	}

	policy2, err := signature.DefaultPolicy(systemContext)
	if err != nil {
		logrus.WithError(err).Fatal("tarball stuff policy 2 error")
	}

	pc, err := signature.NewPolicyContext(policy2)
	if err != nil {
		logrus.WithError(err).Fatal("policy 2 context creation error")
	}

	fmt.Printf("\n img to be created from tar src ref: %+v \n", src)
	//defer pc.Destroy()
	manifest2, err := copy.Image(context.TODO(), pc, dest, src, &copy.Options{
		ReportWriter: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	if err != nil {
		logrus.WithError(err).Fatal("unable to convert tar to image")
	}

	fmt.Printf("tar converted Image manifest %v \n", string(manifest2))
	//_ = _defaultStore.Wipe()

}

func defaultStore() store.Store {
	options, err := store.DefaultStoreOptions(false, 0)
	fmt.Printf("\n filesystem path is %v \n", options.GraphRoot)
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

func InitDefaultStoreOptions() {
	options, err := store.DefaultStoreOptions(false, 0)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create default image store options")
	}
	//options.RunRoot = "/run/containers/storage"
	//options.GraphRoot = "/var/lib/containers/storage"
	//options.GraphDriverName = "overlay"
	//_storeOptions = options

	if _defaultStore == nil {
		gotStorage, err := store.GetStore(options)
		if err != nil {
			logrus.WithError(err).Fatal("Could not create image store")
		}
		_defaultStore = gotStorage
	}

}


func Show() {
	imagesNew, _ := _defaultStore.Images()
	//spew.Dump(imagesNew)
	for _, img := range imagesNew {
		fmt.Printf("\n" + img.NamesHistory[0] + "\n")
	}

	if len(imagesNew) == 0 {
		fmt.Printf("\n No images available to show \n")
	}

	ctrs, _ := _defaultStore.Containers()
	for _, c := range ctrs {
		fmt.Printf("\n" + c.ID + "\n")
	}
	if len(ctrs) == 0 {
		fmt.Printf("\n No containers available to show \n")
	}
}

func ClearStuff() {
	_ = _defaultStore.Wipe()
}