package dot

import (
	"bytes"
	"testing"
)

func TestAttributesMapWrite(t *testing.T) {
	cfg := &AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attr("l", "v").Attr("l2", "v2")

	buf := bytes.NewBufferString("")
	cfg.Write(buf, true)

	want := `[l="v",l2="v2"]`
	if got := buf.String(); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestAttributesMapGet(t *testing.T) {
	cfg := &AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attr("k1", "v1").Attr("k2", "v2")

	if got, want := cfg.attributes["k1"], "v1"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	if got, want := cfg.attributes["k2"], "v2"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	if got := cfg.Value("k3"); got != nil {
		t.Errorf("got [%v:%T] want [%v]", got, got, nil)
	}
}

func TestAttributesMapDelete(t *testing.T) {
	cfg := &AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attr("k1", 100).Attr("k2", "v2")

	if got, want := cfg.Value("k1"), 100; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	if got, want := cfg.Value("k2"), "v2"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	cfg.Delete("k1")
	if got := cfg.Value("k1"); got != nil {
		t.Errorf("got [%v:%T] want [%v]", got, got, nil)
	}
}
