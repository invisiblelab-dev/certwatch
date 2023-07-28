package certwatch

import (
	"context"
	"crypto/tls"
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

type DomainSSLInfo struct {
	Domain string
	SSL    SSLInfo
}

type DomainDeadline struct {
	Domain           string
	DaysTillDeadline float64
	OnDeadline       bool
}

type DomainQuery struct {
	Issuer    string
	LastCheck time.Time
	NotAfter  time.Time
}

func (certInfo CertificateInfo) String() string {
	return fmt.Sprintf(" Subject: %s \n Issuer: %s \n NotBefore: %s \n NotAfter: %s\n\n",
		certInfo.Subject,
		certInfo.Issuer,
		certInfo.NotBefore,
		certInfo.NotAfter,
	)
}

func Certificate(domain string) (SSLInfo, error) {
	// Create a new client with a timeout of 5 seconds
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			// TODO: check the need for skipping checks
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url, err := AddHTTPS(domain)
	if err != nil {
		return SSLInfo{}, err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return SSLInfo{}, fmt.Errorf("failed to setup certificate request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return SSLInfo{}, fmt.Errorf("failed to get certificate for %s: %w", url, err)
	}

	defer resp.Body.Close()

	// Create a new instance of SSLInfo
	sslInfo := SSLInfo{
		Version:           resp.TLS.Version,
		HandshakeComplete: resp.TLS.HandshakeComplete,
		DidResume:         resp.TLS.DidResume,
		CipherSuite:       resp.TLS.CipherSuite,
	}

	if len(resp.TLS.PeerCertificates) == 0 {
		return SSLInfo{}, fmt.Errorf("failed to get peer certificates %v", url)
	}

	// Retrieve information about the peer certificates
	cert := resp.TLS.PeerCertificates[0]
	peerCertificate := CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
	}

	sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)

	return sslInfo, nil
}
