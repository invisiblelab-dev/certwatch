package runners

type DomainSSLInfo struct {
	Domain string
	SSL    SSLInfo
}

type SSLInfoArray struct {
	DomainsSSLs []DomainSSLInfo
}

func CertReturn(domainsArray []string) SSLInfoArray {
	sslInfoArray := SSLInfoArray{}

	for _, domain := range domainsArray {
		certificate := Cert_getter(domain, true)
		domainSSLInfo := DomainSSLInfo{Domain: domain, SSL: certificate}
		sslInfoArray.DomainsSSLs = append(sslInfoArray.DomainsSSLs, domainSSLInfo)
	}

	return sslInfoArray
}
