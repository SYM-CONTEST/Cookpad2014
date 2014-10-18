package crawler

import "testing"

func TestHoge(t *testing.T) {
	a := wordCount([]string{"foo", "bar", "foo"})
	if a["foo"] != 2 {
		t.Error("error")
	}
	if a["bar"] != 1 {
		t.Error("error")
	}
}
