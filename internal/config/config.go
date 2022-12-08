package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arunvelsriram/sodexwoe/internal/constants"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type BillConfigs map[string]BillConfig

type Config struct {
	DownloadDir string      `yaml:"download_dir" binding:"required"`
	BillConfigs BillConfigs `yaml:"bills"`
}

type BillConfig struct {
	Type      string `yaml:"type" binding:"required"`
	Label     string `yaml:"label" binding:"required"`
	KeepPages int    `yaml:"keep_pages"`
	Password  string `yaml:"password"`
}

func (c Config) Label(billName string) (string, error) {
	for name, bill := range c.BillConfigs {
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
	names := make([]string, 0, len(c.BillConfigs))
	for name := range c.BillConfigs {
		names = append(names, name)
	}

	return names
}

func LoadConfig() (config Config, err error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return config, err
	}

	file, err := os.ReadFile(filepath.Join(homeDir, constants.DEFAULT_CONFIG_FILE))
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	var downloadDir string
	if downloadDir, err = homedir.Expand(config.DownloadDir); err != nil {
		return config, err
	}
	config.DownloadDir = downloadDir

	return config, err
}
