package runners

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

func CreateSSLJSON(domainsArray []string) {
	sslInfoArray := SSLInfoArray{}

	for _, domain := range domainsArray {
		certificate := Cert_getter(domain, true)
		domainSSLInfo := DomainSSLInfo{Domain: domain, SSL: certificate}
		sslInfoArray.DomainsSSLs = append(sslInfoArray.DomainsSSLs, domainSSLInfo)
	}

	file, _ := json.MarshalIndent(sslInfoArray, "", " ")
	_ = os.WriteFile("UserSSLs/ssls.json", file, 0644)
}
