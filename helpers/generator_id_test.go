package helpers

import "testing"

func TestGeneratorID(t *testing.T) {
	sid := GenerateID(22)
	t.Log(sid)
}
