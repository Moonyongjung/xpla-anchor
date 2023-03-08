package app

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var once sync.Once
var instance *singleton

type singleton struct {
	app AppType
}

// Read and get app.yaml
func AppFile() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func (s *singleton) Get() AppType {
	return s.app
}

func (s *singleton) Read(filePath string) error {
	appType, err := readAppFile(filePath)
	if err != nil {
		return err
	}
	s.app = *appType

	return nil
}

func readAppFile(filePath string) (*AppType, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var appType AppType

	err = yaml.Unmarshal(yamlFile, &appType)
	if err != nil {
		return nil, err
	}

	return &appType, nil
}

func saveAppFile(appFilePath string, appType *AppType) error {
	bytes, err := yaml.Marshal(appType)
	if err != nil {
		return err
	}

	f, err := os.Create(appFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(bytes); err != nil {
		return err
	}

	return nil
}
