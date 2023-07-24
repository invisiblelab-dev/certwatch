package runners

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"time"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
	"github.com/invisiblelab-dev/certwatch/internal/config"

	"github.com/invisiblelab-dev/certwatch/internal/notifications"
)

func Certificate(url string, roots bool) certwatch.SSLInfo {
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
		return certwatch.SSLInfo{}
	}
	defer resp.Body.Close()

	// Create a new instance of SSLInfo
	sslInfo := certwatch.SSLInfo{
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

func peerCertificate(cert *x509.Certificate) certwatch.CertificateInfo {
	certificate := certwatch.CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
	}
	return certificate
}

func GetCertificates() []certwatch.DomainSSLInfo {
	domains := config.ReadYaml()
	sslInfos := []certwatch.DomainSSLInfo{}
	roots := domains.Roots
	for _, domain := range domains.Domains {
		certificate := Certificate(domain.Name, roots)
		domainSSLInfo := certwatch.DomainSSLInfo{Domain: domain.Name, SSL: certificate}
		sslInfos = append(sslInfos, domainSSLInfo)
	}

	return sslInfos
}

func CalculateDaysToDeadline(certificates []certwatch.DomainSSLInfo) ([]certwatch.DomainDeadline, error) {
	domainsDeadlines := []certwatch.DomainDeadline{}
	file := config.ReadYaml()
	for i := 0; i < len(certificates); i++ {
		timeHours := time.Until(certificates[i].SSL.PeerCertificates[0].NotAfter)
		timeDays := timeHours.Hours() / 24
		var onDeadline bool
		if file.Domains[i].Name != certificates[i].Domain {
			return domainsDeadlines, errors.New("domains don't match")
		}
		if timeDays <= float64(file.Domains[i].NotificationDays) {
			onDeadline = true
		} else {
			onDeadline = false
		}
		deadline := certwatch.DomainDeadline{Domain: certificates[i].Domain, DaysTillDeadline: timeDays, OnDeadline: onDeadline}
		domainsDeadlines = append(domainsDeadlines, deadline)
	}
	return domainsDeadlines, nil
}

func RunCheckCertificatesCommand(opts certwatch.CheckCertificatesOptions) {
	fmt.Println("domains", opts.Domains)
	panic("not implemented")
}

func RunCheckAllCertificatesCommand(opts certwatch.CheckAllCertificatesOptions) {
	certificates := GetCertificates()
	domainDeadlines, err := CalculateDaysToDeadline(certificates)
	if err != nil {
		return
	}

	message, err := notifications.ComposeMessage(domainDeadlines)
	if err != nil {
		return
	}

	notifications.SendEmail(message, config.ReadYaml())
}
