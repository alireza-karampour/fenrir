package coredns

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"codeberg.org/bit101/go-ansi"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/kubectl"
	"github.com/alireza-karampour/fenrir/pkg/task"
	"github.com/alireza-karampour/fenrir/pkg/utils"
)

const (
	IMAGE_PATCH_FORMAT string = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
   - ./coredns.yaml
patches:
   - target:
        name: coredns
        namespace: kube-system
        kind: Deployment
     patch: |-
        - op: replace
          path: /spec/template/spec/containers/0/image
          value: %s`
)

type Cmd struct {
	kc *kubectl.Cmd
}

func New() *Cmd {
	return &Cmd{
		kc: kubectl.New(),
	}
}

func (c *Cmd) Export() error {
	err := os.MkdirAll("kustomize", 0777)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("%s %s", path.Join(kubectl.KUBECTL_BIN_DEST, kubectl.KUBECTL_EXE_NAME), "get deployment/coredns -n kube-system -o yaml")
	tExport := task.NewTask(cmd)
	file, err := os.OpenFile(path.Join("kustomize", "coredns.yaml"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	tExport.SetOut(file)
	tExport.SetErr(os.Stderr)
	err = tExport.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cmd) ChangeImage(image string) error {
	err := c.Export()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join("kustomize", "kustomization.yaml"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, IMAGE_PATCH_FORMAT, image)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	res, err := c.kc.Run(fmt.Sprintf("apply -k %s", path.Join("kustomize")), nil)
	msg, _ := io.ReadAll(res)
	if err != nil {
		utils.PrintErr("kubectl error")
		return errors.Join(err, fmt.Errorf("%s", string(msg)))
	}
	utils.PrintOk("kubectl result")
	ansi.NewLine()
	fmt.Println(string(msg))
	utils.Print("removing pods to force new image")
	ansi.NewLine()
	err = c.RemovePods("k8s-app", "kube-dns")
	if err != nil {
		return err
	}
	return nil
}

func (c *Cmd) RemovePods(key string, val string) error {
	res, err := c.kc.Run(fmt.Sprintf("delete pod -A -l %s=%s", key, val), nil)
	msg, _ := io.ReadAll(res)
	if err != nil {
		utils.PrintErr("kubectl error")
		return errors.Join(err, fmt.Errorf("%s", string(msg)))
	}

	utils.PrintOk("kubectl result")
	ansi.NewLine()
	fmt.Println(string(msg))
	return nil
}
