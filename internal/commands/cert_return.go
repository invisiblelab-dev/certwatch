package main

import "fmt"

func main() {
	domainsString := []string{"https://www.youtube.com/", "https://www.facebook.com/", "https://www.invisiblelab.dev/", "https://www.golangprograms.com/"}

	for _, domain := range domainsString {
		fmt.Println(string(cert_getter(domain, true)))
	}
}
