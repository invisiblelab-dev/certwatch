package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	certwatch "github.com/invisiblelab-dev/certwatch"
	"gopkg.in/yaml.v3"
)

func ReadYaml(path string) (certwatch.ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error: ", err)
		return certwatch.ConfigFile{}, err
	}

	var domains certwatch.ConfigFile
	err = yaml.Unmarshal(data, &domains)
	if err != nil {
		fmt.Println("File parsing error: ", err)
		return certwatch.ConfigFile{}, err
	}
	return domains, nil
}

func WriteYaml(data []byte, path string) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("not writing file, error: ", err)
		return err
	}
	return nil
}

func ReadQueries() (map[string]certwatch.DomainQuery, error) {
	queries := make(map[string]certwatch.DomainQuery)

	fileName := "queries.yaml"
	data, err := os.ReadFile(fileName)
	if errors.Is(err, fs.ErrNotExist) {
		return queries, nil
	}
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &queries)
	if err != nil {
		fmt.Println("File unmarshall error: ", err)
		return nil, err
	}

	return queries, nil
}

func WriteQueries(data map[string]certwatch.DomainQuery) error {
	marshalData, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile("queries.yaml", marshalData, 0644)
	if err != nil {
		return err
	}
	return nil
}
