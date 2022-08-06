package ftp

import (
	"testing"
)

func TestFtpConnection(t *testing.T) {
	f := New("IDStemp.xml")

	if f.id != "IDStemp.xml" {
		t.Errorf("Failed to create ID 'IDStemp.xml'")
	}

	if f.address == "" {
		t.Errorf("Failed to create correct address.")
	}

	if f.path == "" {
		t.Errorf("Failed to create correct path.")
	}
}
