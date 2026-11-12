package fastpack

import (
	"path/filepath"

	"holdtotype/internal/trmodel"
)

const (
	Folder = "fast"
	Pivot  = "en"
)

var Languages = []string{"de", "en", "es", "fr", "it", "pl", "uk", "ru"}

type File struct {
	Kind   string
	Name   string
	Size   int64
	SHA256 string
	URL    string
}

type Pack struct {
	From    string
	To      string
	Version string
	Files   []File
}

func (p Pack) Key() string { return p.From + "-" + p.To }

func (p Pack) Dir(models string) string { return filepath.Join(models, Folder, p.Key()) }

func (p Pack) Size() int64 {
	var n int64
	for _, f := range p.Files {
		n += f.Size
	}
	return n
}

func (p Pack) SizeMB() int { return int((p.Size() + 512*1024) / (1024 * 1024)) }

func (p Pack) File(kind string) (File, bool) {
	for _, f := range p.Files {
		if f.Kind == kind {
			return f, true
		}
	}
	return File{}, false
}

func (p Pack) Path(models, kind string) string {
	f, _ := p.File(kind)
	return filepath.Join(p.Dir(models), f.Name)
}

func (f File) Spec() trmodel.Spec { return trmodel.Spec{Name: f.Name, Size: f.Size, Hash: f.SHA256} }

func All() []Pack { return append([]Pack(nil), catalog...) }

func Find(from, to string) (Pack, bool) {
	for _, p := range catalog {
		if p.From == from && p.To == to {
			return p, true
		}
	}
	return Pack{}, false
}

func Route(from, to string) ([]Pack, bool) {
	if from == "" || to == "" || from == to {
		return nil, false
	}
	if p, ok := Find(from, to); ok {
		return []Pack{p}, true
	}
	if from == Pivot || to == Pivot {
		return nil, false
	}
	a, ok1 := Find(from, Pivot)
	b, ok2 := Find(Pivot, to)
	if !ok1 || !ok2 {
		return nil, false
	}
	return []Pack{a, b}, true
}

type State string

const (
	Absent   State = "absent"
	Checking State = "checking"
	Ready    State = "ready"
	Bad      State = "bad"
)

func Known(models string, p Pack) (State, bool) {
	absent, unknown := 0, 0
	for _, f := range p.Files {
		st, ok := trmodel.Known(p.Dir(models), f.Spec())
		if !ok {
			unknown++
			continue
		}
		switch st {
		case trmodel.Absent:
			absent++
		case trmodel.WrongSize, trmodel.WrongHash:
			return Bad, true
		}
	}
	if absent > 0 {
		return Absent, true
	}
	if unknown > 0 {
		return Checking, false
	}
	return Ready, true
}

func Check(models string, p Pack) (State, error) {
	absent := 0
	for _, f := range p.Files {
		st, err := trmodel.Check(p.Dir(models), f.Spec())
		if err != nil {
			return "", err
		}
		switch st {
		case trmodel.Absent:
			absent++
		case trmodel.WrongSize, trmodel.WrongHash:
			return Bad, nil
		}
	}
	if absent > 0 {
		return Absent, nil
	}
	return Ready, nil
}

func Remove(models string, p Pack) error {
	var first error
	for _, f := range p.Files {
		if err := trmodel.Remove(p.Dir(models), f.Spec()); err != nil && first == nil {
			first = err
		}
	}
	return first
}
