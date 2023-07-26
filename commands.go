package certwatch

type AddDomainOptions struct {
	Domain     string
	DaysBefore int32
}

type CheckCertificatesOptions struct {
	Domains []string
}

type CheckAllCertificatesOptions struct {
	Force bool
}
