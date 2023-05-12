package routes

import "testing"

func TestNewAuth(t *testing.T) {
	a := NewAuth("123", "111", false)
	if a.Url != "123" && a.port != "111" && a.secure != false {
		t.Errorf("Result was incorrect, got: %v, want: other things.", a)
	}
}
