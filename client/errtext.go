package main

import (
	"errors"
	"log"
	"strings"

	"holdtotype/internal/errkind"
)

// uiError carries two texts: the English one goes into the log, the localised
// one goes to the window. The log is read by whoever fixes the program, so it
// stays in one language whatever the interface is set to.
type uiError struct {
	en string
	ui string
}

func (e *uiError) Error() string { return e.ui }

func (e *uiError) Log() string { return e.en }

func uiErr(en, ui string) error { return &uiError{en: en, ui: ui} }

// logText picks the English side of an error when it has one.
func logText(err error) string {
	if err == nil {
		return ""
	}
	var ue *uiError
	if errors.As(err, &ue) {
		return ue.en
	}
	return err.Error()
}

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

// uiErrModel says a model file is missing: English for the log, the user's
// language for the window.
func uiErrModel(path string) error {
	return uiErr("model file missing: "+path, trf("err.sherpa.model", path))
}
