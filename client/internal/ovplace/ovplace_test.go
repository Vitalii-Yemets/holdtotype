package ovplace

import "testing"

var wa = Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1040}

func TestAnchorCorners(t *testing.T) {
	cases := []struct {
		mode string
		x, y int32
	}{
		{PosTopLeft, 28, 28},
		{PosTop, (1920 - 640) / 2, 28},
		{PosTopRight, 1920 - 640 - 28, 28},
		{PosLeft, 28, (1040 - 120) / 2},
		{PosRight, 1920 - 640 - 28, (1040 - 120) / 2},
		{PosBottomLeft, 28, 1040 - 120 - 28},
		{PosBottom, (1920 - 640) / 2, 1040 - 120 - 28},
		{PosBottomRight, 1920 - 640 - 28, 1040 - 120 - 28},
	}
	for _, c := range cases {
		x, y := Anchor(c.mode, wa, 640, 120, 28)
		if x != c.x || y != c.y {
			t.Errorf("%s: получено %d,%d, ожидалось %d,%d", c.mode, x, y, c.x, c.y)
		}
	}
}

func TestAnchorOffsetMonitor(t *testing.T) {
	second := Rect{Left: 1920, Top: -200, Right: 3840, Bottom: 880}
	x, y := Anchor(PosTopLeft, second, 640, 120, 28)
	if x != 1948 || y != -172 {
		t.Fatalf("на втором мониторе получено %d,%d — координаты обязаны считаться от его угла", x, y)
	}
}

func TestCustomCentresThePlate(t *testing.T) {
	x, y := Custom(Frac{X: 0.5, Y: 0.5}, wa, 640, 120)
	if x != (1920-640)/2 || y != (1040-120)/2 {
		t.Fatalf("получено %d,%d — середина обязана быть серединой", x, y)
	}
}

func TestCustomNeverLeavesTheScreen(t *testing.T) {
	x, y := Custom(Frac{X: 1, Y: 1}, wa, 640, 120)
	if x != 1920-640 || y != 1040-120 {
		t.Fatalf("получено %d,%d — плашка не должна вылезать за край", x, y)
	}
	x, y = Custom(Frac{X: 0, Y: 0}, wa, 640, 120)
	if x != 0 || y != 0 {
		t.Fatalf("получено %d,%d — и за левый верх тоже", x, y)
	}
}

func TestFracRoundTrips(t *testing.T) {
	f := Frac{X: 0.31, Y: 0.72}
	x, y := Custom(f, wa, 640, 120)
	back := FracOf(x, y, 640, 120, wa)
	if back.X < 0.30 || back.X > 0.32 || back.Y < 0.71 || back.Y > 0.73 {
		t.Fatalf("дробь не вернулась: %v → %d,%d → %v", f, x, y, back)
	}
}

func TestResKeyNamesTheResolution(t *testing.T) {
	if k := ResKey(Rect{Left: 1920, Top: -200, Right: 3840, Bottom: 880}); k != "1920x1080" {
		t.Fatalf("получено %q — ключ должен быть разрешением, не координатами", k)
	}
}

func TestCleanCustomDropsGarbage(t *testing.T) {
	in := map[string]Frac{
		"1920x1080": {X: 0.5, Y: 0.9},
		"junk":      {X: 0.5, Y: 0.5},
		"800x600":   {X: 2, Y: 0.5},
	}
	out := CleanCustom(in)
	if len(out) != 1 || out["1920x1080"].Y != 0.9 {
		t.Fatalf("получено %v — мусорные ключи и дроби вне экрана должны отбрасываться", out)
	}
	if CleanCustom(map[string]Frac{"junk": {}}) != nil {
		t.Fatal("из одного мусора должна получаться пустота")
	}
}

func TestValidKnowsAllNinePlusCustom(t *testing.T) {
	for _, m := range []string{PosBottom, PosTop, PosCaret, PosTopLeft, PosTopRight, PosLeft, PosRight, PosBottomLeft, PosBottomRight, PosCustom} {
		if !Valid(m) {
			t.Errorf("режим %q должен быть допустимым", m)
		}
	}
	if Valid("середина") {
		t.Error("чужое слово не должно проходить")
	}
}
