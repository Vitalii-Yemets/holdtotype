package errkind

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"testing"
)

func TestOfRecognisesNetworkTrouble(t *testing.T) {
	cases := []struct {
		err  error
		want string
	}{
		{nil, None},
		{&net.DNSError{Err: "no such host", Name: "huggingface.co", IsNotFound: true}, DNS},
		{fmt.Errorf(`Get "https://huggingface.co/api/models": dial tcp: lookup huggingface.co: no such host`), DNS},
		{fmt.Errorf("context deadline exceeded (Client.Timeout exceeded while awaiting headers)"), Timeout},
		{fmt.Errorf("dial tcp 127.0.0.1:8910: connectex: No connection could be made"), Down},
		{fmt.Errorf("x509: certificate signed by unknown authority"), Cert},
		{fmt.Errorf("invalid character 'p' after top-level value"), Answer},
		{context.Canceled, Cancelled},
		{fs.ErrNotExist, Missing},
		{fs.ErrPermission, Denied},
		{fmt.Errorf("write models/ggml-small.bin: no space left on device"), DiskFull},
		{fmt.Errorf("fork/exec whisper-server.exe: The system cannot find the file specified."), Missing},
		{errors.New("something else entirely"), Generic},
	}
	for _, c := range cases {
		if got := Of(c.err); got != c.want {
			t.Errorf("Of(%v) = %q, want %q", c.err, got, c.want)
		}
	}
}

func TestHostNamesTheServer(t *testing.T) {
	if got := Host(&net.DNSError{Err: "no such host", Name: "huggingface.co"}); got != "huggingface.co" {
		t.Errorf("Host(DNSError) = %q", got)
	}
	err := fmt.Errorf(`Get "https://example.org/x": dial tcp: lookup example.org: no such host`)
	if got := Host(err); got != "example.org" {
		t.Errorf("Host(wrapped) = %q, want example.org", got)
	}
	if got := Host(errors.New("plain")); got != "" {
		t.Errorf("Host(plain) = %q, want empty", got)
	}
}
