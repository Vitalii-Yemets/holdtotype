package plexfont

import (
	_ "embed"
	"encoding/base64"
	"strings"
	"sync"
)

//go:embed fonts/IBMPlexSans-Regular.ttf
var Regular []byte

//go:embed fonts/IBMPlexSans-SemiBold.ttf
var SemiBold []byte

const Family = "IBM Plex Sans"

var (
	cssOnce sync.Once
	css     string
)

func FaceCSS() string {
	cssOnce.Do(func() {
		var b strings.Builder
		face := func(weight string, data []byte) {
			b.WriteString("@font-face{font-family:'")
			b.WriteString(Family)
			b.WriteString("';font-style:normal;font-display:block;font-weight:")
			b.WriteString(weight)
			b.WriteString(";src:url(data:font/ttf;base64,")
			b.WriteString(base64.StdEncoding.EncodeToString(data))
			b.WriteString(") format('truetype')}")
		}
		face("400", Regular)
		face("600", SemiBold)
		css = b.String()
	})
	return css
}
