package keymatch

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	lctrl  = 0xA2
	rctrl  = 0xA3
	lwin   = 0x5B
	rwin   = 0x5C
	lalt   = 0xA4
	lshift = 0xA0
	kz     = 0x5A
	ke     = 0x45
	kx     = 0x58
	f9     = 0x78
)

var (
	gCtrl  = []uint32{0xA2, 0xA3, 0x11}
	gShift = []uint32{0xA0, 0xA1, 0x10}
	gAlt   = []uint32{0xA4, 0xA5, 0x12}
	gWin   = []uint32{0x5B, 0x5C}
	gZ     = []uint32{kz}
	gE     = []uint32{ke}
	gX     = []uint32{kx}
	gF9    = []uint32{f9}
)

const settle = 200 * time.Millisecond
const grow = 1000 * time.Millisecond

func combo(id string, groups ...[]uint32) Combo { return Combo{ID: id, Groups: groups} }

func mainCombo() Combo { return combo("main", gCtrl, gWin) }

func subsCombo() Combo { return combo("subs", gCtrl, gWin, gZ) }

type rig struct {
	t      *testing.T
	m      *Matcher
	t0     time.Time
	phys   map[uint32]bool
	hidden map[uint32]bool
}

func newRig(t *testing.T, combos ...Combo) *rig {
	r := &rig{t: t, m: New(settle, grow), t0: time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC), phys: map[uint32]bool{}, hidden: map[uint32]bool{}}
	r.m.SetCombos(combos)
	return r
}

func (r *rig) at(ms int) time.Time { return r.t0.Add(time.Duration(ms) * time.Millisecond) }

func (r *rig) held(vk uint32) bool { return r.phys[vk] && !r.hidden[vk] }

func (r *rig) down(vk uint32, ms int) Result {
	fresh := !r.phys[vk]
	r.phys[vk] = true
	res := r.m.Key(vk, true, r.at(ms), r.held)
	if fresh {
		r.hidden[vk] = res.Swallow
	}
	return res
}

func (r *rig) up(vk uint32, ms int) Result {
	delete(r.phys, vk)
	delete(r.hidden, vk)
	return r.m.Key(vk, false, r.at(ms), r.held)
}

func (r *rig) lose(vk uint32) {
	delete(r.phys, vk)
	delete(r.hidden, vk)
}

func (r *rig) tick(ms int) Result { return r.m.Tick(r.at(ms)) }

func ev(k Kind, id string) Event { return Event{k, id} }

func events(res Result) string {
	if len(res.Events) == 0 {
		return "none"
	}
	return fmt.Sprint(res.Events)
}

func (r *rig) want(label string, res Result, want ...Event) {
	r.t.Helper()
	if len(want) == 0 {
		want = nil
	}
	if !reflect.DeepEqual(res.Events, want) {
		r.t.Errorf("%s: events %s, want %v", label, events(res), want)
	}
}

func (r *rig) swallow(label string, res Result, swallow, mask bool) {
	r.t.Helper()
	if res.Swallow != swallow || res.Mask != mask {
		r.t.Errorf("%s: swallow %v mask %v, want swallow %v mask %v", label, res.Swallow, res.Mask, swallow, mask)
	}
}

func notExtendable() []Combo {
	return []Combo{mainCombo(), combo("subs", gCtrl, gAlt, gF9)}
}

func TestNotExtendableFiresAtOnce(t *testing.T) {
	r := newRig(t, notExtendable()...)
	res := r.down(lctrl, 0)
	r.want("ctrl down", res)
	r.swallow("ctrl down", res, false, false)
	res = r.down(lwin, 10)
	r.want("win down", res, ev(Down, "main"))
	r.swallow("win down", res, false, false)
	if res.Wake != 0 {
		t.Errorf("no wait is asked for, got %v", res.Wake)
	}
	res = r.up(lwin, 500)
	r.want("win up", res, ev(Up, "main"))
	r.swallow("win up", res, false, false)
	res = r.up(lctrl, 510)
	r.want("ctrl up", res)
	r.swallow("ctrl up", res, false, false)
}

func TestWinFirstThenCtrl(t *testing.T) {
	r := newRig(t, notExtendable()...)
	r.want("win down", r.down(lwin, 0))
	r.want("ctrl down", r.down(lctrl, 10), ev(Down, "main"))
}

func TestRightKeysDoTheSame(t *testing.T) {
	r := newRig(t, notExtendable()...)
	r.want("rctrl down", r.down(rctrl, 0))
	r.want("rwin down", r.down(rwin, 10), ev(Down, "main"))
	r.want("rwin up", r.up(rwin, 300), ev(Up, "main"))
}

func TestWinAgainWithCtrlDownStartsAgain(t *testing.T) {
	r := newRig(t, notExtendable()...)
	r.down(lctrl, 0)
	r.want("win down", r.down(lwin, 10), ev(Down, "main"))
	r.want("win up", r.up(lwin, 300), ev(Up, "main"))
	r.want("win down again", r.down(lwin, 400), ev(Down, "main"))
	r.want("win repeat", r.down(lwin, 900))
	r.want("ctrl repeat", r.down(lctrl, 930))
	r.want("win up", r.up(lwin, 1000), ev(Up, "main"))
}

func TestKeyHeldBeforeIsNotSwallowed(t *testing.T) {
	r := newRig(t, notExtendable()...)
	res := r.down(kx, 0)
	r.want("x down", res)
	r.swallow("x down", res, false, false)
	r.down(lctrl, 10)
	res = r.down(lwin, 20)
	r.want("win down", res, ev(Down, "main"))
	r.swallow("win down", res, false, false)
	res = r.up(kx, 100)
	r.want("x up", res)
	r.swallow("x up", res, false, false)
}

func TestNonModifierIsSwallowedWithMask(t *testing.T) {
	r := newRig(t, notExtendable()...)
	r.down(lctrl, 0)
	r.down(lalt, 10)
	res := r.down(f9, 20)
	r.want("f9 down", res, ev(Down, "subs"))
	r.swallow("f9 down", res, true, true)
	res = r.down(f9, 520)
	r.want("f9 repeat", res)
	r.swallow("f9 repeat", res, true, false)
	res = r.up(f9, 600)
	r.want("f9 up", res, ev(Up, "subs"))
	r.swallow("f9 up", res, true, false)
	r.want("alt up", r.up(lalt, 610))
	r.want("ctrl up", r.up(lctrl, 620))
}

func TestMainIsNotExtendableByAnotherShape(t *testing.T) {
	m := New(settle, grow)
	m.SetCombos(notExtendable())
	if m.extendable["main"] || m.extendable["subs"] {
		t.Errorf("nothing is extendable: %v", m.extendable)
	}
	m.SetCombos([]Combo{mainCombo(), subsCombo()})
	if !m.extendable["main"] || m.extendable["subs"] {
		t.Errorf("main is extendable, subs is not: %v", m.extendable)
	}
}

func TestGrowsIntoSubtitlesWhileWaiting(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	res := r.down(lwin, 0)
	r.want("win down", res, ev(Arm, "main"))
	r.swallow("win down", res, false, false)
	if res.Wake != settle {
		t.Errorf("wake %v, want %v", res.Wake, settle)
	}
	res = r.down(kz, 80)
	r.want("z down", res, ev(Disarm, "main"), ev(Down, "subs"))
	r.swallow("z down", res, true, true)
	res = r.up(kz, 300)
	r.want("z up", res, ev(Up, "subs"))
	r.swallow("z up", res, true, false)
	r.want("tick", r.tick(400))
	r.want("win up", r.up(lwin, 500))
	r.want("ctrl up", r.up(lctrl, 510))
	r.down(lctrl, 1000)
	r.want("ctrl+win again", r.down(lwin, 1010), ev(Arm, "main"))
}

func TestWaitEndsWithDown(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.want("win down", r.down(lwin, 0), ev(Arm, "main"))
	r.want("tick early", r.tick(199))
	r.want("tick", r.tick(200), ev(Down, "main"))
	r.want("tick again", r.tick(400))
	r.want("win up", r.up(lwin, 1500), ev(Up, "main"))
}

func TestTapGivesDownThenUp(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("win up", r.up(lwin, 90), ev(Down, "main"), ev(Up, "main"))
	r.want("tick", r.tick(200))
}

func TestLateKeyRetracts(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("tick", r.tick(200), ev(Down, "main"))
	res := r.down(kz, 500)
	r.want("z down", res, ev(Retract, "main"), ev(Down, "subs"))
	r.swallow("z down", res, true, true)
	r.want("z up", r.up(kz, 700), ev(Up, "subs"))
	r.want("win up", r.up(lwin, 800))
}

func TestTooLateKeyIsSwallowedOnly(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("tick", r.tick(200), ev(Down, "main"))
	res := r.down(kz, 1200)
	r.want("z down", res)
	r.swallow("z down", res, true, true)
	for ms := 1700; ms < 2000; ms += 30 {
		res = r.down(kz, ms)
		r.want("z repeat", res)
		r.swallow("z repeat", res, true, false)
	}
	res = r.up(kz, 2000)
	r.want("z up", res)
	r.swallow("z up", res, true, false)
	r.want("win up", r.up(lwin, 2100), ev(Up, "main"))
}

func TestGrowBoundary(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("z at the limit", r.down(kz, 1000), ev(Down, "main"), ev(Retract, "main"), ev(Down, "subs"))
	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("z past the limit", r.down(kz, 1001), ev(Down, "main"))
}

func TestLateTimerOrderFollowsTimestamps(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("z down", r.down(kz, 260), ev(Down, "main"), ev(Retract, "main"), ev(Down, "subs"))
	r.want("late tick", r.tick(270))
}

func TestUnrelatedKeyWhileWaiting(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	res := r.down(kx, 50)
	r.want("x down", res)
	r.swallow("x down", res, false, false)
	r.want("tick", r.tick(200), ev(Down, "main"))
	res = r.up(kx, 300)
	r.want("x up", res)
	r.swallow("x up", res, false, false)
}

func TestSecondTapOfZ(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.down(kz, 80)
	r.want("z up", r.up(kz, 200), ev(Up, "subs"))
	r.want("tick", r.tick(280))
	r.want("tick later", r.tick(1000))
	res := r.down(kz, 1200)
	r.want("z again", res, ev(Down, "subs"))
	r.swallow("z again", res, true, true)
	r.want("z up again", r.up(kz, 1300), ev(Up, "subs"))
	r.want("win up", r.up(lwin, 1400))
	r.want("win down", r.down(lwin, 1500), ev(Arm, "main"))
}

func TestExtraModifierInTheMiddle(t *testing.T) {
	tr := combo("tr", gCtrl, gShift, gWin)
	r := newRig(t, mainCombo(), tr)
	r.down(lctrl, 0)
	r.down(lshift, 10)
	res := r.down(lwin, 20)
	r.want("ctrl shift win", res, ev(Down, "tr"))
	r.swallow("ctrl shift win", res, false, false)
	r.up(lwin, 300)
	r.up(lshift, 300)
	r.up(lctrl, 300)
	r.down(lctrl, 1000)
	r.want("ctrl win", r.down(lwin, 1000), ev(Arm, "main"))
	res = r.down(lshift, 1080)
	r.want("shift", res, ev(Disarm, "main"), ev(Down, "tr"))
	r.swallow("shift", res, false, false)
}

func TestChain(t *testing.T) {
	cs := []Combo{mainCombo(), combo("lang", gCtrl, gWin, gE), combo("far", gCtrl, gWin, gE, gF9)}
	r := newRig(t, cs...)
	r.down(lctrl, 0)
	r.want("ctrl win", r.down(lwin, 0), ev(Arm, "main"))
	res := r.down(ke, 60)
	r.want("e", res, ev(Disarm, "main"), ev(Arm, "lang"))
	r.swallow("e", res, true, true)
	if res.Wake != settle {
		t.Errorf("e: wake %v, want %v", res.Wake, settle)
	}
	r.want("stale tick of main", r.tick(200))
	res = r.down(f9, 120)
	r.want("f9", res, ev(Disarm, "lang"), ev(Down, "far"))
	r.swallow("f9", res, true, true)

	r = newRig(t, cs...)
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.down(ke, 60)
	r.want("tick of main", r.tick(200))
	r.want("tick of lang", r.tick(260), ev(Down, "lang"))
	r.want("f9 in time", r.down(f9, 900), ev(Retract, "lang"), ev(Down, "far"))
}

func TestKeyDownBeforeTheCombination(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	res := r.down(kz, 0)
	r.want("z", res)
	r.swallow("z", res, false, false)
	r.down(lctrl, 10)
	res = r.down(lwin, 20)
	r.want("win", res, ev(Down, "subs"))
	r.swallow("win", res, false, false)
	res = r.up(kz, 300)
	r.want("z up", res, ev(Up, "subs"))
	r.swallow("z up", res, false, false)
}

func TestSetCombosEndsTheActiveOne(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("without main while waiting", r.m.SetCombos([]Combo{subsCombo()}), ev(Disarm, "main"))

	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.tick(200)
	r.want("without main while reported", r.m.SetCombos([]Combo{subsCombo()}), ev(Up, "main"))
	r.want("win up after", r.up(lwin, 500))

	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("the same shortcuts again", r.m.SetCombos([]Combo{mainCombo(), subsCombo()}))
	r.want("the wait goes on", r.tick(200), ev(Down, "main"))

	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.tick(200)
	r.want("main moved to other keys", r.m.SetCombos([]Combo{combo("main", gCtrl, gAlt)}), ev(Up, "main"))
}

func TestResetEndsTheActiveOne(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.want("reset while waiting", r.m.Reset(), ev(Disarm, "main"))
	if r.m.Matched() {
		t.Errorf("pressed keys are forgotten")
	}
	r.want("tick after reset", r.tick(300))

	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.tick(200)
	r.want("reset while reported", r.m.Reset(), ev(Up, "main"))
	r.want("reset again", r.m.Reset())
	r.want("win up after reset", r.up(lwin, 500))
}

func TestStaleSwallowedKey(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.down(kz, 80)
	r.up(lwin, 300)
	r.up(lctrl, 310)
	r.lose(kz)
	res := r.down(kz, 2080)
	r.want("stale z alone", res)
	r.swallow("stale z alone", res, false, false)
	if !reflect.DeepEqual(res.Forgotten, []uint32{kz}) {
		t.Errorf("forgotten %v, want z", res.Forgotten)
	}
	r.up(kz, 2100)

	r = newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.down(kz, 80)
	r.up(lwin, 300)
	r.up(lctrl, 310)
	r.lose(kz)
	r.down(lctrl, 2080)
	r.want("ctrl win with a stale z", r.down(lwin, 2090), ev(Arm, "main"))

	r = newRig(t, combo("main", gF9))
	r.want("f9", r.down(f9, 0), ev(Down, "main"))
	for ms := 500; ms <= 3500; ms += 30 {
		res := r.down(f9, ms)
		r.want("f9 repeat", res)
		r.swallow("f9 repeat", res, true, false)
	}
	if !r.m.Matched() {
		t.Errorf("a repeating swallowed key stays pressed")
	}
	r.want("f9 up", r.up(f9, 3600), ev(Up, "main"))
}

func TestReleasedUnseenKeyIsForgotten(t *testing.T) {
	r := newRig(t, mainCombo(), combo("subs", gCtrl, gWin, gX))
	r.down(kx, 0)
	r.lose(kx)
	res := r.down(lctrl, 500)
	if !reflect.DeepEqual(res.Forgotten, []uint32{kx}) {
		t.Errorf("forgotten %v, want x", res.Forgotten)
	}
	r.want("ctrl win", r.down(lwin, 510), ev(Arm, "main"))
}

func TestReleasedUnseenKeyEndsTheActiveOne(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.lose(lwin)
	r.want("x while win is lost", r.down(kx, 100), ev(Disarm, "main"))
	r.want("tick", r.tick(200))

	r = newRig(t, notExtendable()...)
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.lose(lwin)
	r.want("x while win is lost after the report", r.down(kx, 100), ev(Up, "main"))
}

func TestMaskOnlyWithWinOrAlt(t *testing.T) {
	r := newRig(t, combo("main", gCtrl, gShift, gF9))
	r.down(lctrl, 0)
	r.down(lshift, 0)
	res := r.down(f9, 10)
	r.want("ctrl shift f9", res, ev(Down, "main"))
	r.swallow("ctrl shift f9", res, true, false)
}

func TestLoneKey(t *testing.T) {
	r := newRig(t, combo("main", gF9))
	res := r.down(f9, 0)
	r.want("f9", res, ev(Down, "main"))
	r.swallow("f9", res, true, false)
	res = r.down(f9, 500)
	r.want("f9 repeat", res)
	r.swallow("f9 repeat", res, true, false)
	res = r.up(f9, 600)
	r.want("f9 up", res, ev(Up, "main"))
	r.swallow("f9 up", res, true, false)
	res = r.down(f9, 1000)
	r.want("f9 again", res, ev(Down, "main"))
}

func TestMatched(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	if r.m.Matched() {
		t.Errorf("nothing held")
	}
	r.down(lctrl, 0)
	if r.m.Matched() {
		t.Errorf("ctrl alone is not a combination")
	}
	r.down(lwin, 0)
	if !r.m.Matched() {
		t.Errorf("matched while waiting")
	}
	r.tick(200)
	if !r.m.Matched() {
		t.Errorf("matched while reported")
	}
	r.up(lwin, 300)
	if r.m.Matched() {
		t.Errorf("not matched after the release")
	}
}

func TestKeyOutsideTheHeldCombination(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	r.down(lctrl, 0)
	r.down(lwin, 0)
	r.down(kz, 80)
	r.up(kz, 200)
	res := r.down(kx, 300)
	r.want("x", res)
	r.swallow("x", res, false, false)
}

func TestOtherCombinationWhileOneIsActive(t *testing.T) {
	r := newRig(t, notExtendable()...)
	r.down(lctrl, 0)
	r.want("ctrl win", r.down(lwin, 0), ev(Down, "main"))
	r.down(lalt, 100)
	res := r.down(f9, 110)
	r.want("f9 while main is held", res)
	r.swallow("f9 while main is held", res, false, false)
	r.want("win up", r.up(lwin, 200), ev(Up, "main"))
}

func TestKeyUpNeverSeenDown(t *testing.T) {
	r := newRig(t, mainCombo(), subsCombo())
	res := r.up(kz, 0)
	r.want("z up", res)
	r.swallow("z up", res, false, false)
}

func TestTiesKeepRegistrationOrder(t *testing.T) {
	r := newRig(t, combo("first", gCtrl, gAlt), combo("second", gCtrl, gAlt), combo("long", gCtrl, gShift, gWin))
	r.down(lctrl, 0)
	r.want("ctrl alt", r.down(lalt, 0), ev(Down, "first"))
	r = newRig(t, combo("first", gCtrl, gAlt), combo("second", gCtrl, gAlt), combo("long", gCtrl, gShift, gWin))
	r.down(lshift, 0)
	r.down(lwin, 0)
	r.down(lalt, 0)
	r.want("the longest held wins", r.down(lctrl, 10), ev(Down, "long"))
}

func TestDuplicateGroupsAreSets(t *testing.T) {
	a := combo("a", []uint32{0xA2, 0xA3, 0x11}, gWin)
	b := combo("b", []uint32{0x11, 0xA3, 0xA2}, gWin, gZ)
	if !extends(&b, &a) || extends(&a, &b) || extends(&a, &a) {
		t.Errorf("extends compares groups as sets and needs more groups")
	}
	c := combo("c", gAlt, gWin, gZ)
	if extends(&c, &a) {
		t.Errorf("a different group does not extend")
	}
}
