package certwatch

import "time"

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

type Domain struct {
	Name             string `yaml:"name"`
	NotificationDays int    `yaml:"days"`
}

type DomainQuery struct {
	Issuer    string
	LastCheck time.Time
	NotAfter  time.Time
}
