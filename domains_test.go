package certwatch

import (
	"fmt"
	"testing"

	"github.com/invisiblelab-dev/certwatch/test/helpers"
)

func TestAddHttps(t *testing.T) {
	tests := []string{"www.invisiblelab.dev", "www.invisiblelab.dev/", "invisiblelab.dev"}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			httpsUrl, err := AddHTTPS(tt)
			if err != nil {
				fmt.Println(err)
				return
			}

			helpers.Equal(t, httpsUrl, "https://"+tt)
		})
	}

	edgeCases := []string{"https://www.invisiblelab.dev", "http://www.invisiblelab.dev", "foo://www.invisiblelab.dev", "foo://invisiblelab.dev"}
	for _, tt := range edgeCases {
		t.Run(tt, func(t *testing.T) {
			httpsUrl, err := AddHTTPS(tt)
			if err != nil {
				fmt.Println(err)
				return
			}

			helpers.Equal(t, httpsUrl, tt)
		})
	}
}

func TestRemoveHttps(t *testing.T) {
	tests := []string{"https://www.invisiblelab.dev/", "https://invisiblelab.dev/", "http://www.invisiblelab.dev/"}
	expected := []string{"www.invisiblelab.dev/", "invisiblelab.dev/", "www.invisiblelab.dev/"}
	for i, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			url, err := RemoveHTTPS(tt)
			if err != nil {
				fmt.Println(err)
				return
			}

			helpers.Equal(t, url, expected[i])
		})
	}
}
