package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"codeberg.org/bit101/go-ansi"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/helm"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/kubectl"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/minikube"
	"github.com/alireza-karampour/fenrir/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
)

const (
	BOOTSTRAP_FILENAME_FORMAT string = "%s_suite_test.go"
)

var (
	bootstrap *bool
)

//go:embed templates/*
var templates embed.FS

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "fenrir",
	Short: "a cli for setting up e2e test env for Fenrir",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		if *bootstrap {
			return Bootstarp(".")
		}

		err := minikube.Init()
		if err != nil {
			return err
		}

		kc := kubectl.New()
		err = kc.Init()
		if err != nil {
			return err
		}

		helm := helm.New()
		err = helm.Init()
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	bootstrap = RootCmd.Flags().BoolP("bootstrap", "b", false, "bootstrap a test suit for current package")
}

func Bootstarp(dir string) (err error) {
	utils.Print("bootstraping test suite")
	ansi.NewLine()
	defer func() {
		if err != nil {
			utils.PrintErr("failed to bootstrap test suite")
			ansi.NewLine()
			return
		} else {
			utils.PrintOk("bootstraped tests suite")
			ansi.NewLine()
			return
		}
	}()
	tmps, err := template.ParseFS(templates, "./**/*.gotmpl")
	if err != nil {
		return
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		return
	}
	dirName := filepath.Base(path)

	file, err := os.OpenFile(fmt.Sprintf(BOOTSTRAP_FILENAME_FORMAT, dirName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	err = tmps.ExecuteTemplate(file, "boot.gotmpl", map[string]any{
		"Package": dirName,
		"Suite":   strcase.UpperCamelCase(dirName),
	})
	return
}
