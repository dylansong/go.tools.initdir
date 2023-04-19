package appconfig

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	Name        string      `json:"name" yaml:"name"`
	Directories []Directory `json:"directories,omitempty" yaml:"directories,omitempty"`
}

type Config struct {
	Directories []Directory `json:"directories" yaml:"directories"`
}

func ParseConfig(configPath string) (Config, error) {
	var config Config
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	switch {
	case strings.HasSuffix(configPath, ".json"):
		err = json.Unmarshal(content, &config)
	case strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml"):
		err = yaml.Unmarshal(content, &config)
	default:
		return config, fmt.Errorf("unsupported configuration file format")
	}

	return config, err
}

func createDirectories(basePath string, directories []Directory) error {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir.Name)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		if len(dir.Directories) > 0 {
			err = createDirectories(path, dir.Directories)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
