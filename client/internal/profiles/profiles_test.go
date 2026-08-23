package profiles

import "testing"

var sample = []Profile{
	{ID: "clean", Name: "Чистка", Prompt: "p1"},
	{ID: "formal", Name: "Деловой", Prompt: "p2"},
	{ID: "empty", Name: "Без промпта", Prompt: ""},
}

func ids(list []Profile) []string {
	out := make([]string, 0, len(list))
	for _, p := range list {
		out = append(out, p.ID)
	}
	return out
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestChainFollowsListOrderNotCheckOrder(t *testing.T) {
	got := ids(Chain(sample, []string{"formal", "clean"}, ""))
	if !eq(got, []string{"clean", "formal"}) {
		t.Fatalf("порядок берётся не из списка профилей: %v", got)
	}
}

func TestChainSkipsProfilesWithoutPrompt(t *testing.T) {
	got := ids(Chain(sample, []string{"clean", "empty"}, ""))
	if !eq(got, []string{"clean"}) {
		t.Fatalf("профиль без промпта попал в цепочку: %v", got)
	}
}

func TestChainForcedProfileWinsAlone(t *testing.T) {
	got := ids(Chain(sample, []string{"clean", "formal"}, "formal"))
	if !eq(got, []string{"formal"}) {
		t.Fatalf("хоткей профиля не даёт одиночную цепочку: %v", got)
	}
}

func TestChainForcedProfileIgnoresActiveList(t *testing.T) {
	got := ids(Chain(sample, nil, "clean"))
	if !eq(got, []string{"clean"}) {
		t.Fatalf("профиль по хоткею не сработал без отмеченных: %v", got)
	}
}

func TestChainTranslateHasNoPrompts(t *testing.T) {
	if got := Chain(sample, []string{"clean"}, Translate); len(got) != 0 {
		t.Fatalf("перевод по хоткею потащил промпты: %v", ids(got))
	}
}

func TestChainUnknownIDsAreDropped(t *testing.T) {
	got := ids(Chain(sample, []string{"clean", "нет-такого"}, ""))
	if !eq(got, []string{"clean"}) {
		t.Fatalf("несуществующий профиль не выброшен: %v", got)
	}
	if got := Chain(sample, nil, "нет-такого"); len(got) != 0 {
		t.Fatalf("несуществующий профиль по хоткею дал цепочку: %v", ids(got))
	}
}

func TestChainDeduplicates(t *testing.T) {
	got := ids(Chain(sample, []string{"clean", "clean"}, ""))
	if !eq(got, []string{"clean"}) {
		t.Fatalf("повтор в списке продублировал профиль: %v", got)
	}
}

func TestChainEmptyInputs(t *testing.T) {
	if got := Chain(nil, []string{"clean"}, ""); len(got) != 0 {
		t.Fatalf("пустой каталог дал цепочку: %v", ids(got))
	}
	if got := Chain(sample, nil, ""); len(got) != 0 {
		t.Fatalf("без отмеченных профилей цепочка не пуста: %v", ids(got))
	}
}

func TestByID(t *testing.T) {
	if p := ByID(sample, "formal"); p == nil || p.Name != "Деловой" {
		t.Fatalf("поиск по id не нашёл профиль: %+v", p)
	}
	if p := ByID(sample, "нет"); p != nil {
		t.Fatalf("поиск нашёл несуществующий профиль: %+v", p)
	}
}
