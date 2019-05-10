package progress

import (
	"testing"
)

func TestNewBarID(t *testing.T) {
	NewBarID()
}

func Test_newUUID(t *testing.T) {
	_, err := newUUID()
	if err != nil {
		t.Error(err)
	}
}
