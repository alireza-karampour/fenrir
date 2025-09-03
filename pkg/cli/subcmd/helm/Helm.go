package helm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"codeberg.org/bit101/go-ansi"
	"github.com/alireza-karampour/fenrir/pkg/cli"
	"github.com/alireza-karampour/fenrir/pkg/utils"
)

const (
	HELM_DL_URL_FORMAT   string = "https://get.helm.sh/helm-%s-linux-amd64.tar.gz"
	HELM_VERSION         string = "v3.18.6"
	HELM_BIN_DEST        string = "bin"
	HELM_EXE_NAME        string = "helm"
	HELM_CHECKSUM        string = "c153fd9c1173f39aefe8e9aa9f00fd3daf6b40c8ea01e94a0d2f2c1787fc60e0"
	HELM_TAR_NAME_FORMAT string = "helm-%s-linux-amd64.tar.gz"
	HELM_TAR_TARGET_FILE string = "linux-amd64/helm"
	HELM_TAR_DEST        string = "tars"
	HELM_CHART_DIR       string = "charts"
)

var (
	HELM_DL_URL   string = fmt.Sprintf(HELM_DL_URL_FORMAT, HELM_VERSION)
	HELM_TAR_NAME string = fmt.Sprintf(HELM_TAR_NAME_FORMAT, HELM_VERSION)
)

type Cmd struct {
	*cli.Downloadable
}

func New() *Cmd {
	return &Cmd{
		&cli.Downloadable{
			Name:     HELM_EXE_NAME,
			Url:      "",
			Dest:     HELM_BIN_DEST,
			Checksum: HELM_CHECKSUM,
			TarFile: &cli.Tar{
				Name:       HELM_TAR_NAME,
				URL:        HELM_DL_URL,
				Dest:       HELM_TAR_DEST,
				TargetFile: HELM_TAR_TARGET_FILE,
			},
		},
	}
}

func (c *Cmd) Init() error {
	err := c.Download(cli.DL_TAR | cli.DL_GZIP | cli.DL_VERBOSE)
	if err != nil {
		return err
	}
	err = os.MkdirAll(HELM_CHART_DIR, 0777)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(HELM_CHART_DIR)
	if err != nil {
		return err
	}
	for _, v := range entries {
		if v.IsDir() {
			chartPath := path.Join(HELM_CHART_DIR, v.Name())
			res, err := c.Run(fmt.Sprintf("install %s %s", v.Name(), chartPath), nil)
			msg, _ := io.ReadAll(res)
			if err != nil {
				utils.PrintErr(fmt.Sprintf("failed to install chart for %s", v.Name()))
				ansi.NewLine()
				return errors.Join(err, fmt.Errorf("%s", string(msg)))
			}
			_, err = os.Stdout.Write(msg)
			if err != nil {
				return err
			}
			ansi.NewLine()
			utils.PrintOk(fmt.Sprintf("chart for %s installed successfully", v.Name()))
			ansi.NewLine()
		}
	}
	return nil
}
