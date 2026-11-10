package keymatch

import (
	"sort"
	"time"
)

type Combo struct {
	ID     string
	Groups [][]uint32
}

type Kind int

const (
	Arm Kind = iota + 1
	Down
	Up
	Disarm
	Retract
)

func (k Kind) String() string {
	switch k {
	case Arm:
		return "Arm"
	case Down:
		return "Down"
	case Up:
		return "Up"
	case Disarm:
		return "Disarm"
	case Retract:
		return "Retract"
	}
	return "?"
}

type Event struct {
	Kind Kind
	ID   string
}

type Result struct {
	Events    []Event
	Swallow   bool
	Mask      bool
	Wake      time.Duration
	Forgotten []uint32
}

const StaleSwallowed = 1500 * time.Millisecond

type phase int

const (
	idle phase = iota
	waiting
	reported
)

type keyState struct {
	swallowed bool
	last      time.Time
}

type Matcher struct {
	settle     time.Duration
	grow       time.Duration
	combos     []Combo
	extendable map[string]bool
	pressed    map[uint32]*keyState
	phase      phase
	activeID   string
	started    time.Time
}

func New(settle, grow time.Duration) *Matcher {
	return &Matcher{
		settle:     settle,
		grow:       grow,
		extendable: map[string]bool{},
		pressed:    map[uint32]*keyState{},
	}
}

func IsModifier(vk uint32) bool {
	switch vk {
	case 0x10, 0x11, 0x12, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0x5B, 0x5C:
		return true
	}
	return false
}

func masking(vk uint32) bool {
	switch vk {
	case 0x12, 0xA4, 0xA5, 0x5B, 0x5C:
		return true
	}
	return false
}

func (m *Matcher) SetCombos(c []Combo) Result {
	sorted := make([]Combo, 0, len(c))
	for _, x := range c {
		if len(x.Groups) > 0 {
			sorted = append(sorted, x)
		}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return len(sorted[i].Groups) > len(sorted[j].Groups)
	})
	m.combos = sorted
	m.extendable = map[string]bool{}
	for i := range sorted {
		for j := range sorted {
			if i != j && extends(&sorted[j], &sorted[i]) {
				m.extendable[sorted[i].ID] = true
				break
			}
		}
	}
	var r Result
	if m.phase != idle {
		if a := m.find(m.activeID); a == nil || !m.held(a) {
			m.end(&r)
		}
	}
	return r
}

func (m *Matcher) Key(vk uint32, down bool, now time.Time, held func(uint32) bool) Result {
	var r Result
	m.expire(now, &r)
	if !down {
		m.keyUp(vk, &r)
		return r
	}
	ks := m.pressed[vk]
	if ks != nil && ks.swallowed && now.Sub(ks.last) > StaleSwallowed {
		m.forgetReleased(vk, held, &r)
		if !m.inActive(vk) {
			delete(m.pressed, vk)
			r.Forgotten = append(r.Forgotten, vk)
			ks = nil
		}
	}
	if ks != nil {
		ks.last = now
		r.Swallow = ks.swallowed
		return r
	}
	m.forgetStale(vk, now, held, &r)
	ks = &keyState{last: now}
	m.pressed[vk] = ks
	m.fresh(vk, ks, now, &r)
	return r
}

func (m *Matcher) Tick(now time.Time) Result {
	var r Result
	m.expire(now, &r)
	return r
}

func (m *Matcher) Reset() Result {
	var r Result
	if m.phase != idle {
		m.end(&r)
	}
	m.pressed = map[uint32]*keyState{}
	return r
}

func (m *Matcher) Matched() bool {
	for i := range m.combos {
		if m.held(&m.combos[i]) {
			return true
		}
	}
	return false
}

func (m *Matcher) expire(now time.Time, r *Result) {
	if m.phase == waiting && !now.Before(m.started.Add(m.settle)) {
		m.phase = reported
		r.Events = append(r.Events, Event{Down, m.activeID})
	}
}

func (m *Matcher) keyUp(vk uint32, r *Result) {
	ks := m.pressed[vk]
	if ks == nil {
		return
	}
	delete(m.pressed, vk)
	r.Swallow = ks.swallowed
	if m.phase == idle {
		return
	}
	if a := m.find(m.activeID); a != nil && m.held(a) {
		return
	}
	if m.phase == waiting {
		r.Events = append(r.Events, Event{Down, m.activeID})
	}
	r.Events = append(r.Events, Event{Up, m.activeID})
	m.phase = idle
	m.activeID = ""
}

func (m *Matcher) forgetReleased(current uint32, held func(uint32) bool, r *Result) {
	m.forget(current, r, func(vk uint32, ks *keyState) bool {
		return !ks.swallowed && held != nil && !held(vk)
	})
}

func (m *Matcher) forgetStale(current uint32, now time.Time, held func(uint32) bool, r *Result) {
	m.forget(current, r, func(vk uint32, ks *keyState) bool {
		if ks.swallowed {
			return now.Sub(ks.last) > StaleSwallowed && !m.inActive(vk)
		}
		return held != nil && !held(vk)
	})
}

func (m *Matcher) forget(current uint32, r *Result, gone func(uint32, *keyState) bool) {
	var dropped []uint32
	for vk, ks := range m.pressed {
		if vk != current && gone(vk, ks) {
			dropped = append(dropped, vk)
		}
	}
	sort.Slice(dropped, func(i, j int) bool { return dropped[i] < dropped[j] })
	for _, vk := range dropped {
		delete(m.pressed, vk)
	}
	r.Forgotten = append(r.Forgotten, dropped...)
	if m.phase != idle {
		if a := m.find(m.activeID); a == nil || !m.held(a) {
			m.end(r)
		}
	}
}

func (m *Matcher) inActive(vk uint32) bool {
	if m.phase == idle {
		return false
	}
	a := m.find(m.activeID)
	return a != nil && contains(a, vk)
}

func (m *Matcher) fresh(vk uint32, ks *keyState, now time.Time, r *Result) {
	l := m.longestHeldWith(vk)
	if l == nil {
		return
	}
	switch m.phase {
	case idle:
		m.start(l, vk, ks, now, r)
	case waiting:
		a := m.find(m.activeID)
		if a != nil && l.ID != a.ID && extends(l, a) {
			r.Events = append(r.Events, Event{Disarm, a.ID})
			m.phase = idle
			m.start(l, vk, ks, now, r)
		}
	case reported:
		a := m.find(m.activeID)
		if a == nil || l.ID == a.ID || !extends(l, a) {
			return
		}
		if now.Sub(m.started) <= m.grow {
			r.Events = append(r.Events, Event{Retract, a.ID})
			m.phase = idle
			m.start(l, vk, ks, now, r)
			return
		}
		m.swallow(vk, ks, r)
	}
}

func (m *Matcher) start(c *Combo, vk uint32, ks *keyState, now time.Time, r *Result) {
	m.activeID = c.ID
	m.started = now
	if m.extendable[c.ID] {
		m.phase = waiting
		r.Events = append(r.Events, Event{Arm, c.ID})
		r.Wake = m.settle
	} else {
		m.phase = reported
		r.Events = append(r.Events, Event{Down, c.ID})
	}
	m.swallow(vk, ks, r)
}

func (m *Matcher) swallow(vk uint32, ks *keyState, r *Result) {
	if IsModifier(vk) {
		return
	}
	ks.swallowed = true
	r.Swallow = true
	for other := range m.pressed {
		if masking(other) {
			r.Mask = true
			return
		}
	}
}

func (m *Matcher) end(r *Result) {
	if m.phase == waiting {
		r.Events = append(r.Events, Event{Disarm, m.activeID})
	} else {
		r.Events = append(r.Events, Event{Up, m.activeID})
	}
	m.phase = idle
	m.activeID = ""
}

func (m *Matcher) find(id string) *Combo {
	for i := range m.combos {
		if m.combos[i].ID == id {
			return &m.combos[i]
		}
	}
	return nil
}

func (m *Matcher) longestHeldWith(vk uint32) *Combo {
	for i := range m.combos {
		if contains(&m.combos[i], vk) && m.held(&m.combos[i]) {
			return &m.combos[i]
		}
	}
	return nil
}

func (m *Matcher) held(c *Combo) bool {
	if len(c.Groups) == 0 {
		return false
	}
	for _, g := range c.Groups {
		ok := false
		for _, vk := range g {
			if _, down := m.pressed[vk]; down {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func contains(c *Combo, vk uint32) bool {
	for _, g := range c.Groups {
		for _, x := range g {
			if x == vk {
				return true
			}
		}
	}
	return false
}

func sameGroup(a, b []uint32) bool {
	sa := map[uint32]bool{}
	for _, x := range a {
		sa[x] = true
	}
	sb := map[uint32]bool{}
	for _, x := range b {
		if !sa[x] {
			return false
		}
		sb[x] = true
	}
	return len(sa) == len(sb)
}

func extends(b, a *Combo) bool {
	if len(b.Groups) <= len(a.Groups) {
		return false
	}
	for _, ga := range a.Groups {
		found := false
		for _, gb := range b.Groups {
			if sameGroup(ga, gb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
