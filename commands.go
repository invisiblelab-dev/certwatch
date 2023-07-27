package certwatch

type AddDomainOptions struct {
	Domain     string
	DaysBefore int32
	Path       string
}

type CheckCertificatesOptions struct {
	Domains []string
}

type CheckAllCertificatesOptions struct {
	Force bool
	Path  string
}
