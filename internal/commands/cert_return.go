package main

import (
	"encoding/json"
	"os"
)

type DomainSSLInfo struct {
	Domain string
	SSL    SSLInfo
}

type SSLInfoArray struct {
	DomainsSSLs []DomainSSLInfo
}

func main() {
	domainsString := []string{"https://www.youtube.com/", "https://www.facebook.com/", "https://www.invisiblelab.dev/", "https://www.golangprograms.com/"}
	sslInfoArray := SSLInfoArray{}

	for _, domain := range domainsString {
		certificate := cert_getter(domain, true)
		domainSSLInfo := DomainSSLInfo{Domain: domain, SSL: certificate}
		sslInfoArray.DomainsSSLs = append(sslInfoArray.DomainsSSLs, domainSSLInfo)
	}

	file, _ := json.MarshalIndent(sslInfoArray, "", " ")
	_ = os.WriteFile("UserSSLs/ssls.json", file, 0644)

}
