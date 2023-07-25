package runners

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
	"github.com/invisiblelab-dev/certwatch/internal/config"

	"github.com/invisiblelab-dev/certwatch/internal/notifications"
)

func Certificate(url string) certwatch.SSLInfo {
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
	cert := resp.TLS.PeerCertificates[0]
	peerCertificate := certwatch.CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
	}

	sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)

	return sslInfo
}

func GetCertificates() (map[string]certwatch.DomainQuery, error) {
	// Read config file
	domains, err := config.ReadYaml()
	if err != nil {
		return nil, err
	}

	// Read past queries file
	queries, err := config.ReadQueries()
	if err != nil {
		fmt.Println("error reading queries: ", err)
		return nil, err
	}

	for _, domain := range domains.Domains {
		// Check if already queried or if query was done in more than deadline "days"
		if (queries[domain.Name].LastCheck == time.Time{}) || (int(time.Until(queries[domain.Name].LastCheck).Hours()) >= domain.NotificationDays) {
			certificate := Certificate(domain.Name)
			queries[domain.Name] = certwatch.DomainQuery{Issuer: certificate.PeerCertificates[0].Issuer, LastCheck: time.Now(), NotAfter: certificate.PeerCertificates[0].NotAfter}
		}
	}

	err = config.WriteQueries(queries)
	if err != nil {
		fmt.Println("error writing query file err:", err)
		return nil, err
	}
	return queries, nil
}

func CalculateDaysToDeadline(certificates map[string]certwatch.DomainQuery) ([]certwatch.DomainDeadline, error) {
	domainsDeadlines := []certwatch.DomainDeadline{}
	file, err := config.ReadYaml()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(file.Domains); i++ {
		timeHours := time.Until(certificates[file.Domains[i].Name].NotAfter)
		timeDays := timeHours.Hours() / 24
		var onDeadline bool
		if timeDays <= float64(file.Domains[i].NotificationDays) {
			onDeadline = true
		} else {
			onDeadline = false
		}
		deadline := certwatch.DomainDeadline{Domain: file.Domains[i].Name, DaysTillDeadline: timeDays, OnDeadline: onDeadline}
		domainsDeadlines = append(domainsDeadlines, deadline)
	}
	return domainsDeadlines, nil
}

func RunCheckCertificatesCommand(opts certwatch.CheckCertificatesOptions) {
	fmt.Println("domains", opts.Domains)
	panic("not implemented")
}

func RunCheckAllCertificatesCommand(opts certwatch.CheckAllCertificatesOptions) {
	certificates, err := GetCertificates()
	if err != nil {
		return
	}

	domainDeadlines, err := CalculateDaysToDeadline(certificates)
	if err != nil {
		return
	}

	message, err := notifications.ComposeMessage(domainDeadlines)
	if err != nil {
		return
	}

	file, err := config.ReadYaml()
	if err != nil {
		return
	}

	notifications.SendEmail(message, file.Notifications.Email)
}
