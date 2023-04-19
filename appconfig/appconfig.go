package appconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Directory struct {
	Name        string      `json:"name" yaml:"name"`
	Directories []Directory `json:"directories,omitempty" yaml:"directories,omitempty"`
}

type Config struct {
	Directories []Directory `json:"directories" yaml:"directories"`
}

//func ParseConfig(configPath string) (Config, error) {
//	var config Config
//	content, err := ioutil.ReadFile(configPath)
//	if err != nil {
//		return config, err
//	}
//
//	var m map[string]interface{}
//	err = yaml.Unmarshal(content, &m)
//	if err != nil {
//		return config, err
//	}
//
//	config.Directories = make([]Directory, 0)
//	for name, value := range m {
//		directories := parseDirectories(value)
//		config.Directories = append(config.Directories, Directory{Name: name, Directories: directories})
//	}
//
//	return config, nil
//}

func ParseConfig(configPath string) (Config, error) {
	var config Config
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func CreateDirectories(basePath string, directories []Directory) error {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir.Name)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		if len(dir.Directories) > 0 {
			err = CreateDirectories(path, dir.Directories)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//func parseDirectories(value interface{}) []Directory {
//	directories := make([]Directory, 0)
//	if subMap, ok := value.(map[interface{}]interface{}); ok {
//		for name, subValue := range subMap {
//			subDirs := parseDirectories(subValue)
//			directories = append(directories, Directory{Name: name.(string), Directories: subDirs})
//		}
//	} else if value == nil {
//		directories = append(directories, Directory{Name: "", Directories: nil})
//	}
//	return directories
//}

func parseDirectories(value interface{}) []Directory {
	directories := make([]Directory, 0)
	if subMap, ok := value.(map[interface{}]interface{}); ok {
		for name, subValue := range subMap {
			subDirs := parseDirectories(subValue)
			directories = append(directories, Directory{Name: name.(string), Directories: subDirs})
		}
	} else if subList, ok := value.([]interface{}); ok {
		for _, listItem := range subList {
			subMap, ok := listItem.(map[interface{}]interface{})
			if ok {
				for name, subValue := range subMap {
					subDirs := parseDirectories(subValue)
					directories = append(directories, Directory{Name: name.(string), Directories: subDirs})
				}
			}
		}
	} else if value == nil {
		return nil
	}
	return directories
}

func createDirectoriesFromMap(path string, m map[interface{}]interface{}) ([]Directory, error) {
	directories := make([]Directory, 0)
	for name, value := range m {
		if subMap, ok := value.(map[interface{}]interface{}); ok {
			subDirs, err := createDirectoriesFromMap(filepath.Join(path, name.(string)), subMap)
			if err != nil {
				return nil, err
			}
			directories = append(directories, Directory{Name: name.(string), Directories: subDirs})
		} else if value == nil {
			directories = append(directories, Directory{Name: name.(string), Directories: nil})
		} else {
			return nil, fmt.Errorf("unsupported directory type")
		}
	}
	return directories, nil
}
