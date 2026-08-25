package plexfont

import (
	_ "embed"
	"encoding/base64"
	"strings"
	"sync"
)

//go:embed fonts/IBMPlexSans-Regular.ttf
var sansRegular []byte

//go:embed fonts/IBMPlexSans-SemiBold.ttf
var sansSemiBold []byte

//go:embed fonts/IBMPlexMono-Regular.ttf
var monoRegular []byte

//go:embed fonts/IBMPlexMono-SemiBold.ttf
var monoSemiBold []byte

const (
	Sans = "IBM Plex Sans"
	Mono = "IBM Plex Mono"
)

type Face struct {
	Family string
	Weight string
	Data   []byte
}

func Faces() []Face {
	return []Face{
		{Sans, "400", sansRegular},
		{Sans, "600", sansSemiBold},
		{Mono, "400", monoRegular},
		{Mono, "600", monoSemiBold},
	}
}

var (
	cssOnce sync.Once
	css     string
)

func FaceCSS() string {
	cssOnce.Do(func() {
		var b strings.Builder
		for _, f := range Faces() {
			b.WriteString("@font-face{font-family:'")
			b.WriteString(f.Family)
			b.WriteString("';font-style:normal;font-display:block;font-weight:")
			b.WriteString(f.Weight)
			b.WriteString(";src:url(data:font/ttf;base64,")
			b.WriteString(base64.StdEncoding.EncodeToString(f.Data))
			b.WriteString(") format('truetype')}")
		}
		css = b.String()
	})
	return css
}
