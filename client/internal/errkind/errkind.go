package errkind

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"strings"
)

const (
	None      = ""
	DNS       = "dns"
	Timeout   = "timeout"
	Down      = "down"
	Cert      = "cert"
	Answer    = "answer"
	Missing   = "missing"
	Denied    = "denied"
	DiskFull  = "full"
	Cancelled = "cancelled"
	Generic   = "generic"
)

func Of(err error) string {
	if err == nil {
		return None
	}
	if errors.Is(err, context.Canceled) {
		return Cancelled
	}
	if errors.Is(err, fs.ErrNotExist) {
		return Missing
	}
	if errors.Is(err, fs.ErrPermission) {
		return Denied
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return DNS
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return Timeout
	}
	s := strings.ToLower(err.Error())
	switch {
	case strings.Contains(s, "no such host"):
		return DNS
	case strings.Contains(s, "timeout") || strings.Contains(s, "deadline exceeded"):
		return Timeout
	case strings.Contains(s, "certificate") || strings.Contains(s, "x509"):
		return Cert
	case strings.Contains(s, "connection refused") || strings.Contains(s, "connectex") ||
		strings.Contains(s, "network is unreachable") || strings.Contains(s, "dial tcp") ||
		strings.Contains(s, "connection reset") || strings.Contains(s, "unexpected eof"):
		return Down
	case strings.Contains(s, "no space") || strings.Contains(s, "not enough space"):
		return DiskFull
	case strings.Contains(s, "being used by another process") || strings.Contains(s, "access is denied") ||
		strings.Contains(s, "permission denied"):
		return Denied
	case strings.Contains(s, "cannot find the file") || strings.Contains(s, "cannot find the path") ||
		strings.Contains(s, "no such file"):
		return Missing
	case strings.Contains(s, "invalid character") || strings.Contains(s, "unexpected end of json") ||
		strings.Contains(s, "cannot unmarshal") || strings.Contains(s, "http 5"):
		return Answer
	}
	return Generic
}

func Host(err error) string {
	if err == nil {
		return ""
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) && dnsErr.Name != "" {
		return dnsErr.Name
	}
	s := err.Error()
	i := strings.Index(s, "lookup ")
	if i < 0 {
		return ""
	}
	rest := s[i+len("lookup "):]
	for j, r := range rest {
		if r == ' ' || r == ':' || r == '"' {
			return rest[:j]
		}
	}
	return rest
}
