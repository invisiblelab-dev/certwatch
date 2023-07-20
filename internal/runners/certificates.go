package runners

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

type SSLInfo struct {
	Version           uint16
	HandshakeComplete bool
	DidResume         bool
	CipherSuite       uint16
	PeerCertificates  []CertificateInfo
}

type CertificateInfo struct {
	Subject            string
	Issuer             string
	NotBefore          time.Time
	NotAfter           time.Time
	SignatureAlgorithm string
	PublicKeyAlgorithm string
}

func CertGetter(url string, onlyLeaf bool) SSLInfo {

	// Create a new client with a timeout of 5 seconds
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Make a GET request to the URL
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return SSLInfo{}
	}
	defer resp.Body.Close()

	// Create a new instance of SSLInfo
	sslInfo := SSLInfo{
		Version:           resp.TLS.Version,
		HandshakeComplete: resp.TLS.HandshakeComplete,
		DidResume:         resp.TLS.DidResume,
		CipherSuite:       resp.TLS.CipherSuite,
	}

	// Retrieve information about the peer certificates
	if onlyLeaf {
		cert := resp.TLS.PeerCertificates[0]
		peerCertificate := peerCertificate(cert)
		sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)

	} else {
		for _, cert := range resp.TLS.PeerCertificates {
			peerCertificate := peerCertificate(cert)
			sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)
		}
	}

	return sslInfo
}

func peerCertificate(cert *x509.Certificate) CertificateInfo {
	certificate := CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
	}
	return certificate
}

type DomainSSLInfo struct {
	Domain string
	SSL    SSLInfo
}

type SSLInfoArray struct {
	DomainsSSLs []DomainSSLInfo
}

func CertArrayGetter(domainsArray []string) SSLInfoArray {
	sslInfoArray := SSLInfoArray{}

	for _, domain := range domainsArray {
		certificate := CertGetter(domain, true)
		domainSSLInfo := DomainSSLInfo{Domain: domain, SSL: certificate}
		sslInfoArray.DomainsSSLs = append(sslInfoArray.DomainsSSLs, domainSSLInfo)
	}

	return sslInfoArray
}

type DomainDeadline struct {
	Domain           string
	DaysTillDeadline float64
}

type DomainsDeadlines struct {
	Deadlines []DomainDeadline
}

func CalculateDeadline(certificates SSLInfoArray) DomainsDeadlines {
	domainsDeadlines := DomainsDeadlines{}
	for i := 0; i < len(certificates.DomainsSSLs); i++ {
		timeHours := time.Until(certificates.DomainsSSLs[i].SSL.PeerCertificates[0].NotAfter)
		timeDays := timeHours.Hours() / 24
		deadline := DomainDeadline{Domain: certificates.DomainsSSLs[i].Domain, DaysTillDeadline: timeDays}
		domainsDeadlines.Deadlines = append(domainsDeadlines.Deadlines, deadline)
	}
	return domainsDeadlines
}
