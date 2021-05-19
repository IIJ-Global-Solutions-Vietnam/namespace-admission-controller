package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type ControllerConfig struct {
	ClusterID       string `required:"true" envconfig:"cluster_id"`
	RancherAPIToken string `required:"true" envconfig:"rancher_api_token"`
	RancherURL      string `required:"true" envconfig:"rancher_url"`
	CertFilePath    string `required:"true" envconfig:"certfile_path"`
	KeyFilePath     string `required:"true" envconfig:"keyfile_path"`
	Debug           bool   `envconfig:"debug" default:"false"`
}

var Config ControllerConfig

func init() {
	err := envconfig.Process("", &Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
