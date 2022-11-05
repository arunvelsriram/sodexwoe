package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/arunvelsriram/sodexwoe/constants"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config map[string]Bill

type Bill struct {
	Type      string `yaml:"type" binding:"required"`
	Label     string `yaml:"label" binding:"required"`
	KeepPages int    `yaml:"keep_pages"`
	Password  string `yaml:"password"`
}

func (c Config) Label(billName string) (string, error) {
	for name, bill := range c {
		if strings.EqualFold(billName, name) {
			log.Debugf("label identified: %v", bill.Label)
			return bill.Label, nil
		}
	}
	return "", fmt.Errorf("could not find label for bill name: %v", billName)
}

func (c Config) Labels(billNames []string) ([]string, error) {
	labels := make([]string, 0, len(billNames))
	for _, billName := range billNames {
		label, err := c.Label(billName)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}

func (c Config) BillNames() []string {
	names := make([]string, 0, len(c))
	for name := range c {
		names = append(names, name)
	}

	return names
}

func LoadConfig() (config Config, err error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Debug("unable to identify the home directory")
		return config, err
	}

	file, err := ioutil.ReadFile(filepath.Join(homeDir, constants.DEFAULT_CONFIG_FILE))
	if err != nil {
		log.WithField("file", file).Debug("unable to read config file")
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.WithField("file", file).Debug("unable to unmarshal yaml data")
		return config, err
	}

	return config, err
}
