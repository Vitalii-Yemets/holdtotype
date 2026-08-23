package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func autorunEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runRegKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(runRegValue)
	return err == nil && v != ""
}

func setAutorun(enable bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runRegKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if !enable {
		err := k.DeleteValue(runRegValue)
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return k.SetStringValue(runRegValue, `"`+filepath.Clean(exe)+`"`)
}
