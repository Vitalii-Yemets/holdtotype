package advisor

import (
	"strings"
	"testing"
)

func catalog() []Model {
	return []Model{
		{ID: "base", Engine: "whisper", Langs: []string{"*"}, SizeMB: 142, RAMMB: 273, Translate: true, Speed: 5, Accuracy: 2},
		{ID: "small", Engine: "whisper", Langs: []string{"*"}, SizeMB: 466, RAMMB: 759, Translate: true, Speed: 3, Accuracy: 3},
		{ID: "medium", Engine: "whisper", Langs: []string{"*"}, SizeMB: 539, RAMMB: 869, Translate: true, Speed: 2, Accuracy: 4},
		{ID: "turbo", Engine: "whisper", Langs: []string{"*"}, SizeMB: 574, RAMMB: 921, Translate: false, Speed: 4, Accuracy: 4},
		{ID: "gigaam", Engine: "sherpa", Langs: []string{"ru"}, SizeMB: 232, RAMMB: 278, Punct: true, Speed: 5, Accuracy: 5},
	}
}

func has(why []string, want string) bool {
	for _, w := range why {
		if w == want {
			return true
		}
	}
	return false
}

func TestRussianPrefersDedicatedModel(t *testing.T) {
	r := Recommend(Input{Lang: "ru", Priority: PriorityBalance, RAMFreeMB: 8000}, catalog())
	if r.Primary != "gigaam" {
		t.Fatalf("получено %q, для русского ожидалась gigaam", r.Primary)
	}
	if !has(r.Why, WhyLanguage) {
		t.Errorf("причина «язык» потеряна: %v", r.Why)
	}
}

func TestRussianGetsWhisperCompanionForOtherLanguages(t *testing.T) {
	r := Recommend(Input{Lang: "ru", Priority: PriorityBalance, RAMFreeMB: 8000}, catalog())
	if r.Companion == "" {
		t.Fatal("к русской модели должна прилагаться универсальная")
	}
	if r.Companion == "turbo" {
		t.Error("в спутники не годится модель без перевода")
	}
	if !has(r.Why, WhyCompanion) {
		t.Errorf("причина «спутник» потеряна: %v", r.Why)
	}
}

func TestMultiLanguageStaysOnWhisper(t *testing.T) {
	r := Recommend(Input{Lang: "multi", Priority: PriorityBalance, RAMFreeMB: 8000}, catalog())
	if r.Primary == "gigaam" {
		t.Fatal("для нескольких языков нельзя рекомендовать одноязычную модель")
	}
	if r.Companion != "" {
		t.Error("универсальной модели спутник не нужен")
	}
}

func TestSpeedAndAccuracyChangeThePick(t *testing.T) {
	fast := Recommend(Input{Lang: "multi", Priority: PrioritySpeed, RAMFreeMB: 8000}, catalog())
	slowButGood := Recommend(Input{Lang: "multi", Priority: PriorityAccuracy, RAMFreeMB: 8000}, catalog())
	if fast.Primary != "base" {
		t.Errorf("при упоре на скорость получено %q, ожидалась base", fast.Primary)
	}
	if slowButGood.Primary == "base" {
		t.Error("при упоре на точность base — неправильный выбор")
	}
	if !has(fast.Why, WhySpeed) || !has(slowButGood.Why, WhyAccuracy) {
		t.Errorf("причины не отражают приоритет: %v / %v", fast.Why, slowButGood.Why)
	}
}

func TestTightMemoryPicksSmallerModel(t *testing.T) {
	r := Recommend(Input{Lang: "multi", Priority: PriorityAccuracy, RAMFreeMB: 400}, catalog())
	if r.Primary != "base" {
		t.Fatalf("при 400 МБ свободной памяти получено %q, ожидалась base", r.Primary)
	}
	if !has(r.Why, WhyRAM) {
		t.Errorf("нехватка памяти должна попадать в причины: %v", r.Why)
	}
}

func TestNoMemoryAtAll(t *testing.T) {
	r := Recommend(Input{Lang: "multi", Priority: PriorityBalance, RAMFreeMB: 100}, catalog())
	if r.Primary != "" {
		t.Fatalf("при 100 МБ ничего рекомендовать нельзя, получено %q", r.Primary)
	}
	if !has(r.Why, WhyRAM) {
		t.Errorf("причина должна объяснять отказ: %v", r.Why)
	}
}

func TestTranslationExcludesNonTranslatingModels(t *testing.T) {
	r := Recommend(Input{Lang: "ru", Priority: PriorityBalance, RAMFreeMB: 8000, Translate: true}, catalog())
	if r.Primary == "gigaam" || r.Primary == "turbo" {
		t.Fatalf("получено %q — эта модель не переводит", r.Primary)
	}
}

func TestUnknownRequirementsFallBackToUniversal(t *testing.T) {
	r := Recommend(Input{RAMFreeMB: 8000}, catalog())
	if r.Primary == "" || r.Primary == "gigaam" {
		t.Fatalf("получено %q, без указания языка ожидалась универсальная модель", r.Primary)
	}
}

func TestEmptyCatalog(t *testing.T) {
	r := Recommend(Input{Lang: "ru", RAMFreeMB: 8000}, nil)
	if r.Primary != "" || !has(r.Why, WhyNothing) {
		t.Fatalf("пустой каталог должен давать пустой ответ: %+v", r)
	}
}

func TestZeroMemoryMeansUnknownNotZero(t *testing.T) {
	r := Recommend(Input{Lang: "multi", Priority: PriorityAccuracy, RAMFreeMB: 0}, catalog())
	if r.Primary == "" {
		t.Fatal("неизвестный объём памяти не должен блокировать рекомендацию")
	}
	if has(r.Why, WhyRAM) {
		t.Errorf("не надо жаловаться на память, которую не измерили: %v", r.Why)
	}
	if strings.TrimSpace(r.Primary) != r.Primary {
		t.Error("идентификатор модели не должен содержать пробелов")
	}
}
