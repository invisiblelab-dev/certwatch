package certwatch

type CheckCertificatesOptions struct {
	Domains []string
}

type CheckAllCertificatesOptions struct {
	Force bool
	Path  string
}
