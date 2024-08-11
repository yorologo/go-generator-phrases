package generator

import "testing"

func TestMain(t *testing.T) {
	t.Log(New().Generate())
}
