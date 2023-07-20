package main

import "github.com/invisiblelab-dev/certwatch/internal/runners"

func main() {
	domainsString := []string{"https://www.youtube.com/", "https://www.facebook.com/", "https://www.invisiblelab.dev/", "https://www.golangprograms.com/"}

	runners.CreateSSLJSON(domainsString)
}
