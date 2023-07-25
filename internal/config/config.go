package config

import (
	"fmt"
	"os"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
	"gopkg.in/yaml.v3"
)

func ReadYaml() (certwatch.ConfigFile, error) {
	data, err := os.ReadFile("certwatch.yaml")
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

func WriteYaml(data []byte) error {
	err := os.WriteFile("certwatch.yaml", data, 0644)
	if err != nil {
		fmt.Println("not writing file, error: ", err)
		return err
	}
	return nil
}

func ReadQueries() (map[string]certwatch.DomainQuery, error) {
	fileName := "queries.yaml"
	data, err := os.ReadFile(fileName)
	certwatch.MissingFile(fileName, data, err)

	queries := make(map[string]certwatch.DomainQuery)

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
