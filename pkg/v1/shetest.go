package v1

import (
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/helm"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/kubectl"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/minikube"
)

func init() {
	kc := kubectl.New()
	err := kc.Init()
	if err != nil {
		panic(err)
	}

	helm := helm.New()
	err = helm.Init()
	if err != nil {
		panic(err)
	}

	err = minikube.Init()
	if err != nil {
		panic(err)
	}
}
