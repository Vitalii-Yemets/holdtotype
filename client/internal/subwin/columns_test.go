package subwin

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLayoutReadsAnythingElseAsRows(t *testing.T) {
	for in, want := range map[string]string{"rows": "rows", "columns": "columns", "": "rows", "cols": "rows", "Columns": "rows", " columns": "rows"} {
		if got := Layout(in); got != want {
			t.Errorf("Layout(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestColumnsOnlyWithTwoParts(t *testing.T) {
	cases := []struct {
		layout string
		parts  int
		want   bool
	}{
		{"columns", 2, true},
		{"columns", 1, false},
		{"columns", 0, false},
		{"rows", 2, false},
		{"other", 2, false},
	}
	for _, c := range cases {
		if got := InColumns(c.layout, c.parts); got != c.want {
			t.Errorf("InColumns(%q, %d) = %v", c.layout, c.parts, got)
		}
	}
}

func TestMinRows(t *testing.T) {
	if MinRows(2, false) != 2 || MinRows(2, true) != 1 || MinRows(1, false) != 1 || MinRows(1, true) != 1 || MinRows(0, false) != 1 {
		t.Fatal("two rows only for two parts one above the other")
	}
}

func TestColumnBorderDefaultsToTheMiddle(t *testing.T) {
	for _, col := range []int32{0, -5, 1000, 4000} {
		if x := ColX(900, 16, 120, col); x != 450 {
			t.Errorf("ColX with %d = %d, want the middle", col, x)
		}
	}
	if x := ColX(901, 16, 120, ColDefault); x != 451 {
		t.Errorf("odd width: %d", x)
	}
}

func TestColumnBorderKeepsItsShareWhenTheWindowWidens(t *testing.T) {
	col := ColAt(300, 900, 16, 120)
	if col != 333 {
		t.Fatalf("share %d", col)
	}
	for _, c := range []struct{ w, want int32 }{{900, 300}, {1800, 599}, {1200, 400}} {
		if x := ColX(c.w, 16, 120, col); x != c.want {
			t.Errorf("width %d: border at %d, want %d", c.w, x, c.want)
		}
	}
}

func TestColumnBorderStopsAtTheNarrowestColumn(t *testing.T) {
	lo, hi := int32(2*16+120), int32(900-2*16-120)
	if x := ColX(900, 16, 120, 10); x != lo {
		t.Errorf("far left: %d, want %d", x, lo)
	}
	if x := ColX(900, 16, 120, 990); x != hi {
		t.Errorf("far right: %d, want %d", x, hi)
	}
	if c := ColAt(0, 900, 16, 120); ColX(900, 16, 120, c) != lo {
		t.Errorf("dragged to the left edge: %d", c)
	}
	if c := ColAt(900, 900, 16, 120); ColX(900, 16, 120, c) != hi {
		t.Errorf("dragged to the right edge: %d", c)
	}
	if x := ColX(200, 16, 120, 300); x != 100 {
		t.Errorf("too narrow for two columns: %d, want the middle", x)
	}
	if c := ColAt(30, 200, 16, 120); c != ColDefault {
		t.Errorf("too narrow to drag: %d", c)
	}
	if c := ColAt(30, 0, 16, 120); c != ColDefault {
		t.Errorf("no width: %d", c)
	}
}

func TestColumnBorderFollowsTheMouse(t *testing.T) {
	for _, x := range []int32{200, 300, 450, 511, 600, 700} {
		c := ColAt(x, 900, 16, 120)
		if got := ColX(900, 16, 120, c); got < x-1 || got > x+1 {
			t.Errorf("the mouse at %d puts the border at %d", x, got)
		}
	}
}

func TestColumnAreasLeaveThePaddingAtTheBorder(t *testing.T) {
	l, r, x := ColAreas(900, 62, 200, 16, 120, ColDefault)
	if x != 450 {
		t.Fatalf("border %d", x)
	}
	if l != (Area{Left: 16, Top: 62, Right: 434, Bottom: 200}) || r != (Area{Left: 466, Top: 62, Right: 884, Bottom: 200}) {
		t.Fatalf("%+v %+v", l, r)
	}
	if l.Right-l.Left != r.Right-r.Left {
		t.Fatal("the halves differ")
	}
	l, r, _ = ColAreas(900, 62, 200, 16, 120, 1)
	if l.Right-l.Left != 120 || r.Right-r.Left < 120 {
		t.Fatalf("narrowest %+v %+v", l, r)
	}
}

func TestHitColumn(t *testing.T) {
	cases := []struct {
		x, y int32
		want bool
	}{
		{450, 100, true},
		{444, 100, true},
		{456, 100, true},
		{443, 100, false},
		{457, 100, false},
		{450, 52, true},
		{450, 51, false},
		{450, 189, true},
		{450, 190, false},
	}
	for _, c := range cases {
		if got := HitColumn(c.x, c.y, 450, 52, 190, 6); got != c.want {
			t.Errorf("HitColumn(%d, %d) = %v", c.x, c.y, got)
		}
	}
}

func TestOldSizeRecordsReadAsBefore(t *testing.T) {
	var boxes map[string]Box
	if err := json.Unmarshal([]byte(`{"1920x1080":{"w":900,"rows":6,"top":2},"2560x1440":{"w":1000,"rows":4}}`), &boxes); err != nil {
		t.Fatal(err)
	}
	want := map[string]Box{"1920x1080": {W: 900, Rows: 6, Top: 2}, "2560x1440": {W: 1000, Rows: 4}}
	if got := CleanBoxes(boxes); !reflect.DeepEqual(got, want) {
		t.Fatalf("%+v", got)
	}
	if ColX(900, 16, 120, want["1920x1080"].Col) != 450 {
		t.Fatal("an old record has its border in the middle")
	}
}

func TestTheColumnBorderIsSavedBesideTheRows(t *testing.T) {
	in := map[string]Box{
		"1920x1080": {W: 900, Rows: 6, Top: 2, Col: 333},
		"2560x1440": {W: 900, Rows: 6, Top: 2, Col: 1000},
		"1280x720":  {W: 900, Rows: 6, Col: -4},
		"1600x900":  {W: 900, Rows: 1, Top: 0, Col: 700},
	}
	want := map[string]Box{
		"1920x1080": {W: 900, Rows: 6, Top: 2, Col: 333},
		"2560x1440": {W: 900, Rows: 6, Top: 2},
		"1280x720":  {W: 900, Rows: 6},
		"1600x900":  {W: 900, Rows: 1, Col: 700},
	}
	if got := CleanBoxes(in); !reflect.DeepEqual(got, want) {
		t.Fatalf("%+v", got)
	}
	out, _ := json.Marshal(Box{W: 900, Rows: 6, Top: 2})
	if string(out) != `{"w":900,"rows":6,"top":2}` {
		t.Fatalf("a record without a border is written as before: %s", out)
	}
	out, _ = json.Marshal(Box{W: 900, Rows: 6, Top: 2, Col: 333})
	if string(out) != `{"w":900,"rows":6,"top":2,"col":333}` {
		t.Fatalf("%s", out)
	}
}
