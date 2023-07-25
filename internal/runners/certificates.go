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

func GetCertificates(configData certwatch.ConfigFile) (map[string]certwatch.DomainQuery, error) {
	queries, err := config.ReadQueries()
	if err != nil {
		fmt.Println("error reading queries: ", err)
		return nil, err
	}
	refresh := configData.Refresh

	for _, domain := range configData.Domains {
		// Check if domain was not queried yet or if it was queried but in more than "refresh" seconds
		if (queries[domain.Name].LastCheck == time.Time{} || int(time.Since(queries[domain.Name].LastCheck).Seconds()) >= refresh) {
			certificate := Certificate(domain.Name)
			peerCertificate := certificate.PeerCertificates[0] // TODO if no certificate return error
			queries[domain.Name] = certwatch.DomainQuery{
				Issuer:    peerCertificate.Issuer,
				LastCheck: time.Now(),
				NotAfter:  peerCertificate.NotAfter,
			}
		}
	}

	err = config.WriteQueries(queries)
	if err != nil {
		fmt.Println("error writing query file err:", err)
		return nil, err
	}
	return queries, nil
}

func CalculateDaysToDeadline(certificates map[string]certwatch.DomainQuery, configData certwatch.ConfigFile) ([]certwatch.DomainDeadline, error) {
	domainsDeadlines := []certwatch.DomainDeadline{}
	for _, domain := range configData.Domains {
		timeHours := time.Until(certificates[domain.Name].NotAfter)
		timeDays := timeHours.Hours() / 24
		onDeadline := timeDays <= float64(domain.NotificationDays)

		deadline := certwatch.DomainDeadline{
			Domain:           domain.Name,
			DaysTillDeadline: timeDays,
			OnDeadline:       onDeadline,
		}
		domainsDeadlines = append(domainsDeadlines, deadline)
	}
	return domainsDeadlines, nil
}

func RunCheckCertificatesCommand(opts certwatch.CheckCertificatesOptions) {
	fmt.Println("domains", opts.Domains)
	panic("not implemented")
}

func RunCheckAllCertificatesCommand(opts certwatch.CheckAllCertificatesOptions) {
	configData, err := config.ReadYaml()
	if err != nil {
		fmt.Println("could not read config file")
		return
	}

	certificates, err := GetCertificates(configData)
	if err != nil {
		return
	}

	domainDeadlines, err := CalculateDaysToDeadline(certificates, configData)
	if err != nil {
		return
	}

	message, err := notifications.ComposeMessage(domainDeadlines)
	if err != nil {
		return
	}

	notifications.SendEmail(message, configData.Notifications.Email)
}
