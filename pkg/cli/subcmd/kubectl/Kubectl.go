package kubectl

import (
	"fmt"

	"github.com/alireza-karampour/fenrir/pkg/cli"
)

const (
	KUBECTL_DL_URL_FORMAT string = "https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl"
	KUBECTL_VERSION       string = "v1.33.0"
	KUBECTL_CHECKSUM      string = "9efe8d3facb23e1618cba36fb1c4e15ac9dc3ed5a2c2e18109e4a66b2bac12dc"
	KUBECTL_BIN_DEST      string = "bin"
	KUBECTL_EXE_NAME      string = "kubectl"
)

var (
	KUBECTL_DL_URL string = fmt.Sprintf(KUBECTL_DL_URL_FORMAT, KUBECTL_VERSION)
)

type Cmd struct {
	*cli.Downloadable
}

func New() *Cmd {
	d := &cli.Downloadable{
		Name:     KUBECTL_EXE_NAME,
		Url:      KUBECTL_DL_URL,
		Dest:     KUBECTL_BIN_DEST,
		Checksum: KUBECTL_CHECKSUM,
	}
	return &Cmd{Downloadable: d}
}

func (k *Cmd) Init() (err error) {
	err = k.Download(cli.DL_VERBOSE)
	return
}
