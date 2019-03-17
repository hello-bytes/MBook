package app

import (
	"testing"
)

func TestRandomKey(t *testing.T) {
	rnd := randomKey(16)
	if len(rnd) == 16 {
		t.Log("success", rnd)
	} else {
		t.Error("Failured", rnd)
	}

}

func TestMakeSk(t *testing.T) {
	MakeSecurityKey()
	if len(SecurityKey()) == 48 {
		t.Log("success", SecurityKey())
	} else {
		t.Error("Failured", SecurityKey())
	}
}

func TestHasSk(t *testing.T) {
	MakeSecurityKey()
	if HasSecurityKey() {
		t.Log("success", SecurityKey())
	} else {
		t.Error("Failured", SecurityKey())
	}
}
