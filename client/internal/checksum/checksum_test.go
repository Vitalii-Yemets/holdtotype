package checksum

import (
	"os"
	"path/filepath"
	"testing"
)

const emptySHA = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func write(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("файл не создан: %v", err)
	}
	return path
}

func TestFileHashesContent(t *testing.T) {
	path := write(t, "a.bin", "")
	got, err := File(path)
	if err != nil {
		t.Fatalf("хеш не посчитан: %v", err)
	}
	if got != emptySHA {
		t.Fatalf("хеш пустого файла неверен: %s", got)
	}
}

func TestFileMissing(t *testing.T) {
	if _, err := File(filepath.Join(t.TempDir(), "нет.bin")); err == nil {
		t.Fatal("отсутствующий файл не дал ошибки")
	}
}

func TestVerifyMatches(t *testing.T) {
	path := write(t, "b.bin", "")
	if err := Verify(path, emptySHA); err != nil {
		t.Fatalf("совпадающий хеш признан несовпадающим: %v", err)
	}
	if err := Verify(path, "SHA256:"+emptySHA); err != nil {
		t.Fatalf("префикс sha256: не распознан: %v", err)
	}
	if err := Verify(path, "  "+emptySHA+"  "); err != nil {
		t.Fatalf("пробелы вокруг хеша сломали проверку: %v", err)
	}
}

func TestVerifyMismatch(t *testing.T) {
	path := write(t, "c.bin", "не пусто")
	if err := Verify(path, emptySHA); err == nil {
		t.Fatal("несовпадение не поймано")
	}
}

func TestVerifyEmptyWantIsSkipped(t *testing.T) {
	path := write(t, "d.bin", "что угодно")
	if err := Verify(path, ""); err != nil {
		t.Fatalf("без эталона проверка должна пропускаться: %v", err)
	}
}

func TestNormalize(t *testing.T) {
	if got := Normalize(" SHA256:ABC "); got != "abc" {
		t.Fatalf("нормализация неверна: %q", got)
	}
}
