package theme

import (
	"strings"
	"testing"
)

func TestTerminalButtonsWearBrackets(t *testing.T) {
	v := Current("terminal", "green").CSSVars()
	if !strings.Contains(v, `--btnbo:"[`) || !strings.Contains(v, `]"`) {
		t.Fatalf("терминал обязан одевать кнопки в скобки, получено: %.120s", v[strings.Index(v, "--btnbo"):])
	}
}

func TestOtherSkinsKeepButtonsPlain(t *testing.T) {
	for _, id := range []string{"editor", "neon", "soft", "paper"} {
		v := Current(id, "").CSSVars()
		if !strings.Contains(v, `--btnbo:""`) {
			t.Errorf("облик %s не должен рисовать скобки", id)
		}
	}
}
