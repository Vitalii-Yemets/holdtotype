package subtitles

import (
	"strings"
	"unicode"
)

var linkingLists = map[string]string{
	"en": "and but so because which that when where if or while although though whether then whereas since unless until who what why how " +
		"on in at from with for by about into than as after before over between through without during against across " +
		"within under around toward towards upon per via like including behind beyond above below inside outside onto despite except among along throughout versus",
	"ru": "и а но или либо что чтобы потому поскольку когда где если как чем хотя пока однако который которая которое которые которого которой которую которым которых " +
		"в во на с со к ко по из от до для за под над при про без через между перед после около " +
		"из-за из-под внутри вокруг возле среди вдоль против вместо кроме согласно благодаря несмотря сквозь ради вне мимо вроде насчёт",
	"uk": "і й та а але або чи що щоб бо тому коли де якщо як ніж хоча поки проте який яка яке які якого якої яку яким яких " +
		"в у на з із зі до від по для за під над при про без через між перед після " +
		"з-за з-під всередині навколо біля серед уздовж вздовж проти замість крім окрім згідно завдяки попри крізь заради поза повз щодо",
	"de": "und aber oder denn weil dass wenn als ob obwohl während damit sondern wo wie was nachdem bevor sobald falls sodass " +
		"in im an am auf aus bei mit nach von vom für über unter vor durch gegen ohne um zwischen seit bis " +
		"ins ans aufs ums zum zur beim innerhalb außerhalb oberhalb unterhalb wegen trotz statt anstatt entlang gegenüber hinter neben jenseits laut gemäß dank per via inklusive einschließlich ab",
	"fr": "et mais ou donc car que qu' qui quand où si parce puisque puisqu' lorsque lorsqu' comme alors tandis " +
		"à au aux dans en sur sous avec sans pour par chez vers entre avant après pendant depuis contre " +
		"selon malgré parmi envers durant hors outre derrière devant jusqu' jusque auprès dès via hormis excepté",
	"es": "y e o u pero sino porque que cuando donde si como aunque mientras pues " +
		"a al en con sin para por sobre entre desde hasta hacia durante contra tras " +
		"bajo ante según mediante dentro fuera detrás delante encima debajo alrededor cerca lejos excepto salvo vía versus",
	"it": "e ed o ma perché che quando dove se come mentre quindi però " +
		"a al alla all' in nel nella nell' con su sul sulla sull' per tra fra da dal dalla dall' senza verso durante contro " +
		"ai agli alle allo nei negli nelle nello sui sugli sulle sullo dai dagli dalle dallo sotto sopra dentro fuori dietro davanti oltre presso lungo attraverso secondo tramite mediante entro eccetto tranne nonostante intorno",
	"pl": "i a ale lub albo czy że żeby bo ponieważ kiedy gdy gdzie jeśli jeżeli jak niż choć chociaż więc lecz który która które którzy którego której którą którym których " +
		"w we na z ze do od po dla za pod nad przy o bez przez między przed " +
		"u ku zza spod znad sprzed wewnątrz wokół obok wśród wzdłuż przeciw przeciwko zamiast oprócz poza według dzięki mimo pomimo spośród ponad poprzez naprzeciw koło około",
}

var linkingSets = func() map[string]map[string]bool {
	out := map[string]map[string]bool{}
	all := map[string]bool{}
	for lang, list := range linkingLists {
		set := map[string]bool{}
		for _, w := range strings.Fields(list) {
			set[w] = true
			all[w] = true
		}
		out[lang] = set
	}
	out[""] = all
	return out
}()

func linkingSet(lang string) map[string]bool {
	if set, ok := linkingSets[strings.ToLower(strings.TrimSpace(lang))]; ok {
		return set
	}
	return linkingSets[""]
}

func linkingKey(w string) string {
	w = strings.ToLower(strings.ReplaceAll(w, "’", "'"))
	if i := strings.IndexRune(w, '\''); i > 0 {
		if head := strings.TrimLeftFunc(w[:i+1], func(r rune) bool { return !unicode.IsLetter(r) && r != '\'' }); head != "'" {
			return head
		}
	}
	return strings.TrimFunc(w, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}

func IsLinkingWord(lang, w string) bool {
	return linkingSet(lang)[linkingKey(w)]
}
