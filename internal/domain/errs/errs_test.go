package errs

import (
	"errors"
	"testing"
)

func TestWrapKeepsKind(t *testing.T) {
	base := NotFoundf("courier not found")
	wrapped := Wrap(Wrap(base, "repo: get"), "service: create")

	if KindOf(wrapped) != NotFound {
		t.Errorf("kind = %v, want NotFound", KindOf(wrapped))
	}
	if !errors.Is(wrapped, base) {
		t.Error("errors.Is should find the sentinel through wraps")
	}
	if Public(wrapped) != "courier not found" {
		t.Errorf("public = %q", Public(wrapped))
	}
	if wrapped.Error() != "service: create: repo: get: courier not found" {
		t.Errorf("error = %q", wrapped.Error())
	}
}

func TestWrapUnknownIsInternal(t *testing.T) {
	raw := errors.New("connection refused")
	wrapped := Wrap(raw, "repo: list")

	if KindOf(wrapped) != Internal {
		t.Errorf("kind = %v, want Internal", KindOf(wrapped))
	}
	if Public(wrapped) != "repo: list" {
		t.Errorf("public = %q", Public(wrapped))
	}
	if !errors.Is(wrapped, raw) {
		t.Error("cause should still be reachable")
	}
}

func TestWrapNil(t *testing.T) {
	if Wrap(nil, "x") != nil {
		t.Error("wrapping nil must stay nil")
	}
}

func TestPublicRaw(t *testing.T) {
	if Public(errors.New("boom")) != "internal error" {
		t.Error("raw errors must not leak their message")
	}
}
