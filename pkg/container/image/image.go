package image

import (
	_ "github.com/containers/buildah"
	_ "github.com/containers/image/storage"
	"github.com/containers/storage"
	_ "github.com/sirupsen/logrus"
)

func init() {
	storeOpt, err := storage.DefaultStoreOptions()
	if err != nil {
		panic(err)
	}

}
