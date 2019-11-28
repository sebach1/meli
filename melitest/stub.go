package melitest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Stub struct {
	Status            int
	Body              interface{}
	WantBodyReceive   []byte
	WantParamsReceive url.Values
}

type Client interface {
	SetClient(http.Client)
}

func (s *Stub) Serve(t *testing.T, cl Client) (Close func()) {
	t.Helper()
	if s == nil {
		return func() {}
	}
	sv := httptest.NewTLSServer(s.StubBody(t, s.Body, s.Status))
	cl.SetClient(StubClient(sv))
	return sv.Close
}

func StubClient(sv *httptest.Server) http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
			return net.Dial(network, sv.Listener.Addr().String())
		},
	}
	return http.Client{Transport: transport}
}

func (s *Stub) StubBody(t *testing.T, body interface{}, Status int) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		s.assertReceive(t, req)
		res.WriteHeader(Status)
		err := json.NewEncoder(res).Encode(body)
		if err != nil {
			t.Fatalf("couldn't marshal given Stub to body: %v", err)
		}
	})
}

func (s *Stub) assertReceive(t *testing.T, req *http.Request) {
	t.Helper()
	assertBody(t, req, s.WantBodyReceive)
	if s.WantParamsReceive == nil {
		s.WantParamsReceive = url.Values{}
	}
	assertParams(t, req, s.WantParamsReceive)
}

func assertBody(t *testing.T, req *http.Request, assertion []byte) {
	t.Helper()
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil && err != io.EOF {
		t.Fatalf("couldn't read body req: %v", err)
	}
	if diff := cmp.Diff(string(assertion), string(body)); diff != "" {
		t.Errorf("Received another BODY on server () mismatch (-want +got): %s", diff)
	}
}

func assertParams(t *testing.T, req *http.Request, assertion url.Values) {
	t.Helper()
	got := req.URL.Query()

	if diff := cmp.Diff(assertion, got); diff != "" {
		t.Errorf("Received another PARAMS on server () mismatch (-want +got): %s", diff)
	}
}
