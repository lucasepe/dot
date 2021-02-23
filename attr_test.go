package dot

import (
	"bytes"
	"testing"
)

func TestAttributesMapWrite(t *testing.T) {
	cfg := AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attrs("l", "v", "l2", "v2")

	buf := bytes.NewBufferString("")
	cfg.Write(buf, true)

	want := `[l="v",l2="v2"]`
	if got := buf.String(); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestAttributesMapGetHappy(t *testing.T) {
	cfg := AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attrs("l", "v", "l2", "v2")

	if got, want := cfg.attributes["l"], "v"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	if got, want := cfg.attributes["l2"], "v2"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestAttributesMapPanics(t *testing.T) {
	// No need to check whether `recover()` is nil.
	// Just turn off the panic.
	defer func() { recover() }()

	cfg := AttributesMap{
		attributes: make(map[string]interface{}),
	}
	cfg.Attrs("l", "v", "l2")

	// Never reaches here if `cfg.Attrs` panics.
	t.Errorf("must provide pairs did not panic")
}
