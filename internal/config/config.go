package config

import (
	"fmt"
	"os"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
	"gopkg.in/yaml.v3"
)

func ReadYaml() certwatch.ConfigFile {
	data, err := os.ReadFile("certwatch.yaml")
	if err != nil {
		fmt.Println("File reading error: ", err)
	}

	var domains certwatch.ConfigFile
	err = yaml.Unmarshal(data, &domains)
	if err != nil {
		fmt.Println("File parsing error: ", err)
	}
	return domains
}

func WriteYaml(data []byte) error {
	err := os.WriteFile("certwatch.yaml", data, 0644)
	if err != nil {
		fmt.Println("not writing file, error: ", err)
		return err
	}
	return nil
}
