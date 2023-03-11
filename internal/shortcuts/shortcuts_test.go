package shortcuts

import (
	"path/filepath"
	"testing"
)

const (
	testdataDir = "../../testdata"
)

func TestExpand(t *testing.T) {
	testMap := map[string]string{
		"exp":   "Expand the following text. Use concise language, an academic tone, avoid unecessary words",
		"tr-fr": "Translate the following text from Italian to French",
	}

	got := Expand(testMap, "@@tr-fr Ciao Mondo!")
	t.Logf(got)

	got = Expand(testMap, "@@exp Ciao Mondo!")
	t.Logf(got)

	got = Expand(testMap, "@@pty Ciao Mondo!")
	t.Logf(got)

}

func TestFromFile(t *testing.T) {
	expectedValues := map[string]string{
		"exp":  "AAAA",
		"fix":  "BBBB",
		"tren": "CCC ddd",
	}

	fn := filepath.Join(testdataDir, "shortcuts")
	res, err := FromFile(fn)
	if err != nil {
		t.Error("Error reading file")
	}

	if len(res) != len(expectedValues) {
		t.Error("Didn't get the right size map back")
	}

	for key, value := range expectedValues {
		if res[key] != value {
			t.Errorf("expected %s to be %s, got %s", key, value, res[key])
		}
	}
}
