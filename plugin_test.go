package pluggo

import "testing"

func TestRegister1(t *testing.T) {
	err := Register("ep1", func() interface{} {
		return "1"
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegister2(t *testing.T) {
	err := Register("ep2", func() interface{} {
		return "2"
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDuplicated(t *testing.T) {
	err := Register("ep1", func() interface{} {
		return "3"
	})
	if err == nil {
		t.Fatal("expected to fail duplicated registration")
	}
}

func TestGet(t *testing.T) {
	ep1 := Get("ep1").(string)
	if ep1 != "1" {
		t.Fatal("plugin returned unexpected instance")
	}
}

func TestGetUnknown(t *testing.T) {
	ep := Get("ep-unknown")
	if ep != nil {
		t.Fatal("unexpected plugin instance returned for unknown extension point")
	}
}
