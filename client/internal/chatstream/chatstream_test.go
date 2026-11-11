package chatstream

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func collect(body string) (string, []string, error) {
	var seen []string
	out, err := Read(strings.NewReader(body), func(s string) { seen = append(seen, s) })
	return out, seen, err
}

func TestEventStream(t *testing.T) {
	body := ": keep-alive\n\n" +
		"data: {\"choices\":[{\"delta\":{\"role\":\"assistant\"}}]}\n\n" +
		"data: {\"choices\":[{\"delta\":{\"content\":\"а потом\"}}]}\n\n" +
		"event: message\r\n" +
		"data: {\"choices\":[{\"delta\":{\"content\":\" мы вернулись\"}}]}\r\n\r\n" +
		"data:{\"choices\":[{\"delta\":{\"content\":\" в офис,\"}}]}\n\n" +
		"data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n" +
		"data: [DONE]\n\n" +
		"data: {\"choices\":[{\"delta\":{\"content\":\" after the end\"}}]}\n\n"
	out, seen, err := collect(body)
	if err != nil {
		t.Fatal(err)
	}
	if out != "а потом мы вернулись в офис," {
		t.Fatalf("text %q", out)
	}
	want := []string{"а потом", "а потом мы вернулись", "а потом мы вернулись в офис,"}
	if !reflect.DeepEqual(seen, want) {
		t.Fatalf("calls %q", seen)
	}
}

func TestStreamWithoutTheEndMark(t *testing.T) {
	out, seen, err := collect("data: {\"choices\":[{\"delta\":{\"content\":\"Hallo\"}}]}\n\ndata: {\"choices\":[{\"delta\":{\"content\":\" da\"}}]}")
	if err != nil || out != "Hallo da" || len(seen) != 2 {
		t.Fatalf("%q %q %v", out, seen, err)
	}
	out, seen, err = collect("")
	if err != nil || out != "" || len(seen) != 0 {
		t.Fatalf("empty body: %q %q %v", out, seen, err)
	}
	out, _, err = collect("data: not json\n\ndata: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\ndata: [DONE]\n")
	if err != nil || out != "ok" {
		t.Fatalf("a broken event: %q %v", out, err)
	}
}

func TestPlainJSONIsOnePiece(t *testing.T) {
	body := "\n {\"choices\":[{\"message\":{\"role\":\"assistant\",\"content\":\"Guten Morgen\\nzusammen\"}}]}\n"
	out, seen, err := collect(body)
	if err != nil || out != "Guten Morgen\nzusammen" {
		t.Fatalf("%q %v", out, err)
	}
	if !reflect.DeepEqual(seen, []string{"Guten Morgen\nzusammen"}) {
		t.Fatalf("calls %q", seen)
	}
	if _, _, err := collect(`{"choices":[]}`); err == nil {
		t.Fatal("an answer without choices was taken")
	}
	if _, _, err := collect(`<html>bad gateway</html>`); err != nil {
		t.Fatalf("a body that is neither: %v", err)
	}
}

func TestErrors(t *testing.T) {
	_, _, err := collect(`{"error":{"code":400,"message":"the grammar is invalid"}}`)
	if err == nil || err.Error() != "the grammar is invalid" {
		t.Fatalf("plain error %v", err)
	}
	out, seen, err := collect("data: {\"choices\":[{\"delta\":{\"content\":\"Bon\"}}]}\n\ndata: {\"error\":{\"message\":\"context overflow\"}}\n\n")
	if err == nil || err.Error() != "context overflow" || out != "Bon" || len(seen) != 1 {
		t.Fatalf("error in the stream: %q %v", out, err)
	}
	_, _, err = collect(`{"error":{}}`)
	if err == nil {
		t.Fatal("an empty error was taken as an answer")
	}
	broken := io.MultiReader(strings.NewReader("data: {\"choices\":[{\"delta\":{\"content\":\"Ciao\"}}]}\n"), errReader{})
	out, err = Read(broken, nil)
	if out != "Ciao" || !errors.Is(err, errBroken) {
		t.Fatalf("a broken connection: %q %v", out, err)
	}
}

var errBroken = errors.New("connection reset")

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errBroken }
