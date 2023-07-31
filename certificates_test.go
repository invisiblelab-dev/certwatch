package certwatch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invisiblelab-dev/certwatch/test/helpers"
)

func TestCertificate(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	certificateInfo, err := Certificate(ts.URL)
	fmt.Println(certificateInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	helpers.Equal(t, certificateInfo.PeerCertificates[0].NotAfter, ts.Certificate().NotAfter)
	helpers.Equal(t, certificateInfo.PeerCertificates[0].NotBefore, ts.Certificate().NotBefore)
}
