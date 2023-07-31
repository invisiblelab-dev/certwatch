package config

import (
	"fmt"
	"testing"

	"github.com/invisiblelab-dev/certwatch/test/helpers"
)

func TestReadYaml(t *testing.T) {
	configData, err := ReadYaml("../test/helpers/test.yaml")
	if err != nil {
		fmt.Println("not able to read file, err: ", err)
	}

	helpers.Equal(t, configData.Domains[0].Name, "https://www.invisiblelab.dev/")
	helpers.Equal(t, configData.Domains[0].NotificationDays, 10)
	helpers.Equal(t, configData.Refresh, 10)
	helpers.Equal(t, configData.Notifications.Username, "123abc")
	helpers.Equal(t, configData.Notifications.Password, "123abc")
}

func TestWriteYaml(t *testing.T) {
	configData, err := ReadYaml("../test/helpers/test.yaml")
	if err != nil {
		fmt.Println("not able to read file, err: ", err)
	}

	err = WriteYaml(configData, "../test/helpers/test.yaml")
	helpers.Equal(t, err, nil)
}
