package file

import (
	"testing"
)

func TestFileRetriever(t *testing.T) {
	testC := New("../test.xml")
	if testC == nil {
		t.Errorf("Failed to create test connection")
	}

	data, err := testC.Retrieve()
	if err != nil {
		t.Errorf("Failed to retrieve data: %s", err)
	}
	if data == nil {
		t.Errorf("Failed to retrieve data: data is nil!")
	}

}
