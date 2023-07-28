package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/invisiblelab-dev/certwatch"
	"gopkg.in/yaml.v3"
)

func ReadYaml(path string) (certwatch.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return certwatch.Config{}, fmt.Errorf("failed to read config file [%s]: %w", path, err)
	}

	var domains certwatch.Config
	err = yaml.Unmarshal(data, &domains)
	if err != nil {
		return certwatch.Config{}, fmt.Errorf("failed to unmarshal config file [%s]: %w", path, err)
	}

	return domains, nil
}

func WriteYaml(data []byte, path string) error {
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("[config.WriteFile] failed to write yaml file %s: %w", path, err)
	}

	return nil
}

func ReadQueries() (map[string]certwatch.DomainQuery, error) {
	queries := make(map[string]certwatch.DomainQuery)

	filename := "queries.yaml"
	data, err := os.ReadFile(filename)

	if errors.Is(err, fs.ErrNotExist) {
		return queries, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read queries file [%s]: %w", filename, err)
	}

	err = yaml.Unmarshal(data, &queries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file [%s]: %w", filename, err)
	}

	return queries, nil
}

func WriteQueries(data map[string]certwatch.DomainQuery) error {
	marshalData, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal queries payload: %w", err)
	}

	err = os.WriteFile("queries.yaml", marshalData, 0600)
	if err != nil {
		return fmt.Errorf("failed to writes queries [%s]: %w", "queries.yaml", err)
	}

	return nil
}
