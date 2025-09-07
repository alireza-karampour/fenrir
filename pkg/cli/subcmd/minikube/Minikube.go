package minikube

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/alireza-karampour/fenrir/pkg/cli"
	"github.com/alireza-karampour/fenrir/pkg/task"
	. "github.com/alireza-karampour/fenrir/pkg/utils"
)

var singleton *Cmd

const MINIKUBE_CHECKSUM string = "cddeab5ab86ab98e4900afac9d62384dae0941498dfbe712ae0c8868250bc3d7"
const MINIKUBE_CLI_VER string = "v1.36.0"
const MINIKUBE_DL_URL_FORMAT string = "https://github.com/kubernetes/minikube/releases/download/%s/minikube-linux-amd64"
const MINIKUBE_BIN_DEST string = "bin"
const MINIKUBE_EXE_NAME string = "minikube"
const MINIKUBE_IMAGES_DIR string = "images"

var MK_MINIKUBE_DL_URL string = fmt.Sprintf(MINIKUBE_DL_URL_FORMAT, MINIKUBE_CLI_VER)

type Cmd struct {
	*cli.Downloadable
}

func New() *Cmd {
	if singleton != nil {
		return singleton
	}
	c := &cli.Downloadable{
		Name:     MINIKUBE_EXE_NAME,
		Url:      MK_MINIKUBE_DL_URL,
		Dest:     MINIKUBE_BIN_DEST,
		Checksum: MINIKUBE_CHECKSUM,
	}
	mk := &Cmd{
		Downloadable: c,
	}
	singleton = mk
	return mk
}

func getInstance() *Cmd {
	fn := sync.OnceValue(func() *Cmd {
		return New()
	})
	return fn()
}

func Init() error {
	c := getInstance()
	err := c.Download(cli.DL_VERBOSE)
	if err != nil {
		return err
	}
	err = Start()
	if err != nil {
		return err
	}
	err = Stop()
	if err != nil {
		return err
	}
	err = enableMetallbAddon()
	if err != nil {
		return err
	}
	err = configureMetallb()
	if err != nil {
		return err
	}

	err = Start()
	if err != nil {
		return err
	}

	err = LoadAll(MINIKUBE_IMAGES_DIR)
	if err != nil {
		return err
	}
	return nil
}

func Stop() error {
	res, err := singleton.Run("stop", nil)
	msg, _ := io.ReadAll(res)
	if err != nil {
		PrintlnErr("failed to stop cluster")
		return errors.Join(err, fmt.Errorf("%s", string(msg)))
	}

	return nil
}

func Delete() error {
	res, err := singleton.Run("delete", nil)
	msg, _ := io.ReadAll(res)
	if err != nil {
		PrintlnErr("failed to delete cluster")
		return errors.Join(err, fmt.Errorf("%s", string(msg)))
	}
	PrintlnOk("successfully deleted cluster")
	return nil
}

func Start() (err error) {
	defer func() {
		if err != nil {
			PrintlnErr("minikube failed")
		} else {
			PrintlnOk("minikube started successfully")
		}
	}()
	Println("starting minikube")

	err = task.NewTask(fmt.Sprintf("%s %s", path.Join(MINIKUBE_BIN_DEST, MINIKUBE_EXE_NAME), "start --driver docker --listen-address 0.0.0.0")).SetOut(os.Stdout).Run()
	if err != nil {
		return
	}

	return
}

func LoadAll(root string) (err error) {
	Println(fmt.Sprintf("loading images from %s", MINIKUBE_IMAGES_DIR))
	defer func() {
		if err != nil {
			PrintlnErr("minikube failed to load images")
			return
		} else {
			PrintlnOk("minikube loaded images")
		}
	}()
	err = os.MkdirAll(root, 0777)
	if err != nil {
		return
	}
	err = fs.WalkDir(os.DirFS(root), ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if fileName := strings.Split(info.Name(), "."); len(fileName) > 1 && fileName[1] == "tar" {
				return LoadImage(p)
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func LoadImage(tar string) error {
	c := getInstance()
	Println(fmt.Sprintf("loading image %s", tar))
	res, err := c.Run(fmt.Sprintf("image load %s", tar), nil)
	msg, _ := io.ReadAll(res)
	if err != nil {
		return errors.Join(err, fmt.Errorf("%s", string(msg)))
	}
	PrintlnOk(fmt.Sprintf("loaded image %s", tar))
	return nil
}

func enableMetallbAddon() (err error) {
	defer func() {
		if err != nil {
			PrintlnErr("failed to enable metallb addon")
			return

		} else {
			PrintlnOk("enabled metallb addon")
			return
		}
	}()
	Println("enabling metallb")

	tsk := task.NewTask(fmt.Sprintf("%s %s", path.Join(MINIKUBE_BIN_DEST, MINIKUBE_EXE_NAME), "addons enable metallb"))
	tsk.SetOut(os.Stdout)
	tsk.SetErr(os.Stderr)
	err = tsk.Run()
	if err != nil {
		return
	}

	return
}

func configureMetallb() (err error) {
	defer func() {
		if err != nil {
			PrintlnErr("failed to configure metallb")
			return
		} else {
			PrintlnOk("configured metallb")
			return
		}
	}()
	Println("configuring metallb")

	tsk := task.NewTask(fmt.Sprintf("%s %s", path.Join(MINIKUBE_BIN_DEST, MINIKUBE_EXE_NAME), "addons configure metallb"))
	tsk.SetIn(os.Stdin).
		SetErr(os.Stderr).
		SetOut(os.Stdout)
	err = tsk.Run()
	return
}
