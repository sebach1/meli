package meli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")

	svErrFooBar = &Error{ResponseErr: errFoo.Error(), Message: errBar.Error()}
)

type stub struct {
	status            int
	body              interface{}
	wantBodyReceive   []byte
	wantParamsReceive url.Values
}

func (s *stub) serve(t *testing.T, ml *MeLi) (Close func()) {
	t.Helper()
	if s == nil {
		return func() {}
	}
	sv := httptest.NewTLSServer(s.stubBody(t, s.body, s.status))
	ml.Client = stubClient(sv)
	return sv.Close
}

func stubClient(sv *httptest.Server) http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
			return net.Dial(network, sv.Listener.Addr().String())
		},
	}
	return http.Client{Transport: transport}
}

func (s *stub) stubBody(t *testing.T, body interface{}, status int) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		s.assertReceive(t, req)
		res.WriteHeader(status)
		err := json.NewEncoder(res).Encode(body)
		if err != nil {
			t.Fatalf("couldn't marshal given stub to body: %v", err)
		}
	})
}

func (s *stub) assertReceive(t *testing.T, req *http.Request) {
	t.Helper()
	assertBody(t, req, s.wantBodyReceive)
	if s.wantParamsReceive == nil {
		s.wantParamsReceive = url.Values{}
	}
	assertParams(t, req, s.wantParamsReceive)
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

func JSONMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("couldn't marshal given value: %v", err)
	}
	return bytes
}
