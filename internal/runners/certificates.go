package runners

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func CertArrayGetter() SSLInfoArray {
	domainsArray := ReadDomains()
	sslInfoArray := SSLInfoArray{}

	for _, domain := range domainsArray.Domains {
		certificate := CertGetter(domain.Domain, true)
		domainSSLInfo := DomainSSLInfo{Domain: domain.Domain, SSL: certificate}
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

func CalculateDaysToDeadline(certificates SSLInfoArray) DomainsDeadlines {
	domainsDeadlines := DomainsDeadlines{}
	for i := 0; i < len(certificates.DomainsSSLs); i++ {
		timeHours := time.Until(certificates.DomainsSSLs[i].SSL.PeerCertificates[0].NotAfter)
		timeDays := timeHours.Hours() / 24
		deadline := DomainDeadline{Domain: certificates.DomainsSSLs[i].Domain, DaysTillDeadline: timeDays}
		domainsDeadlines.Deadlines = append(domainsDeadlines.Deadlines, deadline)
	}
	return domainsDeadlines
}

type Domain struct {
	Domain  string    `json:"domain"`
	AddTime time.Time `json:"addTime"`
}

type DomainsArray struct {
	Domains []Domain `json:"domains"`
}

func AddDomain(domain string) error {
	domains := ReadDomains()
	newDomain := Domain{Domain: domain, AddTime: time.Now()}

	for _, listedDomain := range domains.Domains {
		if listedDomain.Domain == domain {
			return errors.New("Domain already added")
		} else {
			continue
		}
	}

	domains.Domains = append(domains.Domains, newDomain)

	file, err := json.MarshalIndent(domains, "", " ")
	if err != nil {
		fmt.Println("not indenting file, error: ", err)
		return err
	}

	err = os.WriteFile("domainsStorage.json", file, 0644)
	if err != nil {
		fmt.Println("not writing file, error: ", err)
		return err
	}
	return nil
}

func ReadDomains() DomainsArray {
	data, err := os.ReadFile("domainsStorage.json")
	if err != nil {
		fmt.Println("File reading error: ", err)
	}

	var domains DomainsArray
	err = json.Unmarshal(data, &domains)
	if err != nil {
		fmt.Println("File parsing error: ", err)
	}
	return domains
}
