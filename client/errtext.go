package main

import (
	"log"
	"strings"

	"holdtotype/internal/errkind"
)

func errText(err error) string {
	if err == nil {
		return ""
	}
	s := strings.TrimSpace(err.Error())
	if r := []rune(s); len(r) > 200 {
		s = string(r[:200]) + "…"
	}
	return s
}

func humanError(err error) string {

	if err == nil {
		return ""
	}
	log.Printf("error detail: %v", err)
	switch errkind.Of(err) {
	case errkind.DNS:
		host := errkind.Host(err)
		if host == "" {
			host = "huggingface.co"
		}
		return trf("err.net.dns", host)
	case errkind.Timeout:
		return tr("err.net.timeout")
	case errkind.Down:
		return tr("err.net.down")
	case errkind.Cert:
		return tr("err.net.cert")
	case errkind.Answer:
		return tr("err.answer")
	case errkind.Missing:
		return tr("err.file.missing")
	case errkind.Denied:
		return tr("err.file.denied")
	case errkind.DiskFull:
		return tr("err.disk.full")
	case errkind.Cancelled:
		return tr("err.cancelled")
	}
	return tr("err.generic")
}
