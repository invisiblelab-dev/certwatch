package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

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

func ReadQueries(path string) (map[string]certwatch.DomainQuery, error) {
	queries := make(map[string]certwatch.DomainQuery)
	pathBuilder := strings.Builder{}

	if path != "" {
		pathBuilder.WriteString(path)
	} else {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return queries, fmt.Errorf("failed to determine cache location: %w", err)
		}
		pathBuilder.WriteString(cacheDir)
	}

	pathBuilder.WriteString("/certwatch/cache.json")

	filePath := pathBuilder.String()

	data, err := os.ReadFile(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		return queries, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read queries file [%s]: %w", filePath, err)
	}

	err = json.Unmarshal(data, &queries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file [%s]: %w", filePath, err)
	}

	return queries, nil
}

func WriteQueries(data map[string]certwatch.DomainQuery, path string) error {
	marshalData, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal queries payload: %w", err)
	}

	pathBuilder := strings.Builder{}

	if path != "" {
		pathBuilder.WriteString(path)
	} else {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("failed to determine cache location: %w", err)
		}
		pathBuilder.WriteString(cacheDir)
	}

	pathBuilder.WriteString("/certwatch")

	cachePath := pathBuilder.String()
	if _, err := os.Stat(cachePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cachePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create cache folder: %w", err)
		}
	}

	pathBuilder.WriteString("/cache.json")
	filePath := pathBuilder.String()

	err = os.WriteFile(filePath, marshalData, 0600)
	if err != nil {
		return fmt.Errorf("failed to writes queries [%s]: %w", filePath, err)
	}
	return nil
}
