package runners

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
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

func Certificate(url string, roots bool) SSLInfo {

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
	if roots {
		for _, cert := range resp.TLS.PeerCertificates {
			peerCertificate := peerCertificate(cert)
			sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)
		}
	} else {
		cert := resp.TLS.PeerCertificates[0]
		peerCertificate := peerCertificate(cert)
		sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)
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

func GetCertificates() SSLInfoArray {
	domainsArray := ReadYaml()
	sslInfoArray := SSLInfoArray{}
	roots := domainsArray.Roots
	for _, domain := range domainsArray.Domains {
		certificate := Certificate(domain.Name, roots)
		domainSSLInfo := DomainSSLInfo{Domain: domain.Name, SSL: certificate}
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
	Name             string `yaml:"name"`
	NotificationDays int    `yaml:"days"`
}

type ConfigFile struct {
	Domains       []Domain `yml:"domains"`
	Roots         bool     `yml:"roots"`
	Notifications struct {
		Email Email `yml:"email"`
	} `yml:"notifications"`
}

type Email struct {
	Mailtrap struct {
		Username string `yml:"username"`
		Password string `yml:"password"`
		SmtpHost string `yml:"smtpHost"`
	} `yml:"mailtrap"`
	From string `yml:"from"`
	To   string `yml:"to"`
}

func AddDomain(domain string, daysToNotify int) error {
	domains := ReadYaml()
	newDomain := Domain{Name: domain, NotificationDays: daysToNotify}

	for _, listedDomain := range domains.Domains {
		if listedDomain.Name == domain {
			return errors.New("Domain already added")
		} else {
			continue
		}
	}

	domains.Domains = append(domains.Domains, newDomain)

	marshalData, err := yaml.Marshal(&domains)
	if err != nil {
		fmt.Println("not marshalling file, error: ", err)
		return err
	}

	err = os.WriteFile("certwatch.yaml", marshalData, 0644)
	if err != nil {
		fmt.Println("not writing file, error: ", err)
		return err
	}
	return nil
}

func ReadYaml() ConfigFile {
	data, err := os.ReadFile("certwatch.yaml")
	if err != nil {
		fmt.Println("File reading error: ", err)
	}

	var domains ConfigFile
	err = yaml.Unmarshal(data, &domains)
	if err != nil {
		fmt.Println("File parsing error: ", err)
	}
	return domains
}
