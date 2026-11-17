package punct

import (
	"encoding/json"
	"math"
	"os"
	"reflect"
	"testing"
)

type reference struct {
	Bos   int32 `json:"bos"`
	Eos   int32 `json:"eos"`
	Lines []struct {
		Lang  string   `json:"lang"`
		Words []string `json:"words"`
		IDs   []int32  `json:"ids"`
	} `json:"lines"`
}

func readReference(t *testing.T, path string) reference {
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var r reference
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func checkTokens(t *testing.T, tok *Tokenizer, ref reference) {
	if len(ref.Lines) == 0 {
		t.Fatal("no lines")
	}
	for _, l := range ref.Lines {
		ids, last := tok.Words(l.Words)
		if ids[0] != ref.Bos || ids[len(ids)-1] != ref.Eos {
			t.Fatalf("%s: the special tokens %d and %d", l.Lang, ids[0], ids[len(ids)-1])
		}
		if !reflect.DeepEqual(ids[1:len(ids)-1], l.IDs) {
			t.Fatalf("%s: got %v, the reference library gives %v", l.Lang, ids[1:len(ids)-1], l.IDs)
		}
		if len(last) != len(l.Words) || last[len(last)-1] != len(ids)-2 {
			t.Fatalf("%s: the last tokens of the words %v", l.Lang, last)
		}
	}
}

func TestTheTokensMatchTheReferenceLibraryInEightLanguages(t *testing.T) {
	tok, err := LoadTokenizer("testdata/tokenizer-mini.json")
	if err != nil {
		t.Fatal(err)
	}
	ref := readReference(t, "testdata/tokens-reference.json")
	langs := map[string]bool{}
	for _, l := range ref.Lines {
		langs[l.Lang] = true
	}
	for _, l := range []string{"en", "ru", "uk", "de", "fr", "es", "it", "pl"} {
		if !langs[l] {
			t.Fatalf("no reference line in %s", l)
		}
	}
	checkTokens(t, tok, ref)
}

func TestTheTokensMatchAReferenceFileWithTheFullTokenizer(t *testing.T) {
	tf, rf := os.Getenv("HT_PUNCT_TOKENIZER"), os.Getenv("HT_PUNCT_REFERENCE")
	if tf == "" || rf == "" {
		t.Skip("HT_PUNCT_TOKENIZER and HT_PUNCT_REFERENCE are not set")
	}
	tok, err := LoadTokenizer(tf)
	if err != nil {
		t.Fatal(err)
	}
	checkTokens(t, tok, readReference(t, rf))
}

func TestAWordOutsideTheVocabularyFallsBackToBytes(t *testing.T) {
	raw := `{"added_tokens":[{"id":1,"content":"<eos>"},{"id":2,"content":"<bos>"}],
	"model":{"type":"BPE","vocab":{"▁":3,"a":4,"▁a":5,"<0xE2>":6,"<0x82>":7,"<0xAC>":8},"merges":[["▁","a"]]}}`
	tok, err := ParseTokenizer([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if got := tok.Word("a€"); !reflect.DeepEqual(got, []int32{5, 6, 7, 8}) {
		t.Fatalf("got %v", got)
	}
}

func TestNormLikeTheProbe(t *testing.T) {
	for in, want := range map[string]string{"Hello,": "hello", "it’s": "it's", "—": "", "(Repaid": "repaid", "3.5.": "3.5", "Kohorte.": "kohorte", "«когорта»": "когорта"} {
		if got := Norm(in); got != want {
			t.Fatalf("%q: %q, want %q", in, got, want)
		}
	}
}

func TestDecideTakesTheLastTokenOfEachWord(t *testing.T) {
	logits := make([]float32, 4*Labels)
	logits[1*Labels+0] = 5
	logits[2*Labels+1] = 5
	logits[3*Labels+3] = 1
	logits[3*Labels+0] = 1
	end, prob := Decide(logits, []int{1, 2, 3})
	if !reflect.DeepEqual(end, []bool{false, true, false}) {
		t.Fatalf("ends %v", end)
	}
	if prob[1] < 0.9 || math.Abs(prob[2]-(math.E+1)/(2*math.E+4)) > 1e-9 {
		t.Fatalf("probabilities %v", prob)
	}
}
