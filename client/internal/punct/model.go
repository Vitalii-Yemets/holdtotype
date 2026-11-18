package punct

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	ModelName = "Horizon Punctuation Small"
	Repo      = "Horizon-Labs/punctuation-restoration-small"
	Revision  = "a1818fae73a1d09b5496e67747dda7bd667b79a8"
	ModelDir  = "punct/horizon"
	RAMMB     = 340
	LangCount = 89
	License   = "Apache 2.0"
)

type File struct {
	Name   string
	Remote string
	Bytes  int64
	SHA256 string
}

var Files = []File{
	{Name: "model_quantized.onnx", Remote: "onnx/model_quantized.onnx", Bytes: 268414381, SHA256: "5679870be92c7fa47c9be9b2872fac389ef28d43d2c33f35e1315b93a82f4ace"},
	{Name: "tokenizer.json", Remote: "tokenizer.json", Bytes: 34363286, SHA256: "f34b44a6dc8a0eef651bfa10e8325057530152606a69d9c3b6ae02efb1558f0a"},
}

func BaseURL() string {
	return "https://huggingface.co/" + Repo + "/resolve/" + Revision + "/"
}

func SizeMB() int {
	var total int64
	for _, f := range Files {
		total += f.Bytes
	}
	return int((total + 1<<20 - 1) >> 20)
}

func Missing(dir string) string {
	for _, f := range Files {
		p := filepath.Join(dir, f.Name)
		if _, err := os.Stat(p); err != nil {
			return fmt.Sprintf("the punctuation model is missing: %s", p)
		}
	}
	return ""
}

func OffReason(on, translating bool, missing string) string {
	switch {
	case !on:
		return "whole sentences are turned off"
	case !translating:
		return "translation is off"
	}
	return missing
}
