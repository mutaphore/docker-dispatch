package dockerdispatch

import (
	"net"
	"net/http"
	"net/http/httputil"
)

func NewUnixTransport(socketPath string) *http.Transport {
	unixTransport := &http.Transport{}
	unixTransport.RegisterProtocol("unix", NewUnixRoundTripper(socketPath))
	return unixTransport
}

func NewUnixRoundTripper(path string) *UnixRoundTripper {
	return &UnixRoundTripper{path: path}
}

type UnixRoundTripper struct {
	path string
	conn httputil.ClientConn
}

// The RoundTripper (http://golang.org/pkg/net/http/#RoundTripper) for
// the socket transport dials the socket each time a request is made.
func (roundTripper UnixRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := net.Dial("unix", roundTripper.path)
	if err != nil {
		return nil, err
	}
	socketClientConn := httputil.NewClientConn(conn, nil)
	defer socketClientConn.Close()
	return socketClientConn.Do(req)
}
