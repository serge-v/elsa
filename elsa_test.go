package elsa

import "testing"

func TestSummary(t *testing.T) {
	s, err := Summary()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestDistance(t *testing.T) {
	d, err := Distance(41, -74)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("distance: %.0f miles", d)
}
