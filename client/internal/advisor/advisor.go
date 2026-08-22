package advisor

import "sort"

const (
	PriorityAccuracy = "accuracy"
	PrioritySpeed    = "speed"
	PriorityBalance  = "balance"
)

const (
	WhyLanguage  = "language"
	WhyAccuracy  = "accuracy"
	WhySpeed     = "speed"
	WhyRAM       = "ram"
	WhyCompanion = "companion"
	WhyNothing   = "nothing"
)

type Model struct {
	ID        string
	Engine    string
	Langs     []string
	SizeMB    int
	RAMMB     int
	Punct     bool
	Translate bool
	Speed     int
	Accuracy  int
}

type Input struct {
	Lang      string
	Priority  string
	RAMFreeMB int
	Translate bool
}

type Result struct {
	Primary   string
	Companion string
	Why       []string
}

func (m Model) covers(lang string) bool {
	for _, l := range m.Langs {
		if l == "*" || l == lang {
			return true
		}
	}
	return false
}

func (m Model) fits(freeMB int) bool {
	if freeMB <= 0 || m.RAMMB <= 0 {
		return true
	}
	return m.RAMMB*10 <= freeMB*8
}

func score(m Model, priority string) int {
	switch priority {
	case PrioritySpeed:
		return m.Speed*3 + m.Accuracy
	case PriorityAccuracy:
		return m.Accuracy*3 + m.Speed
	default:
		return m.Accuracy*2 + m.Speed*2
	}
}

func best(cands []Model, priority string) (Model, bool) {
	if len(cands) == 0 {
		return Model{}, false
	}
	sorted := make([]Model, len(cands))
	copy(sorted, cands)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := score(sorted[i], priority), score(sorted[j], priority)
		if si != sj {
			return si > sj
		}
		return sorted[i].RAMMB < sorted[j].RAMMB
	})
	return sorted[0], true
}

func Recommend(in Input, catalog []Model) Result {
	lang := in.Lang
	if lang == "" {
		lang = "multi"
	}
	multi := lang == "multi"

	var fitting, tooBig []Model
	for _, m := range catalog {
		if in.Translate && !m.Translate {
			continue
		}
		if multi {
			if !m.covers("*") {
				continue
			}
		} else if !m.covers(lang) {
			continue
		}
		if m.fits(in.RAMFreeMB) {
			fitting = append(fitting, m)
		} else {
			tooBig = append(tooBig, m)
		}
	}

	why := []string{}
	pick, ok := best(fitting, in.Priority)
	if !ok {
		if _, over := best(tooBig, in.Priority); over {
			return Result{Why: []string{WhyRAM}}
		}
		return Result{Why: []string{WhyNothing}}
	}
	if len(tooBig) > 0 {
		why = append(why, WhyRAM)
	}
	if !multi && pick.Engine != "whisper" {
		why = append(why, WhyLanguage)
	}
	switch in.Priority {
	case PrioritySpeed:
		why = append(why, WhySpeed)
	case PriorityAccuracy:
		why = append(why, WhyAccuracy)
	}

	res := Result{Primary: pick.ID, Why: why}
	if !multi && !pick.Translate {
		var companions []Model
		for _, m := range catalog {
			if m.ID != pick.ID && m.Translate && m.covers("*") && m.fits(in.RAMFreeMB) {
				companions = append(companions, m)
			}
		}
		if c, found := best(companions, PriorityBalance); found {
			res.Companion = c.ID
			res.Why = append(res.Why, WhyCompanion)
		}
	}
	return res
}
