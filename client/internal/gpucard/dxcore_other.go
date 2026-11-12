//go:build !windows

package gpucard

import "errors"

func Adapters() ([]Adapter, error) {
	return nil, errors.New("graphics adapters are read on Windows only")
}
