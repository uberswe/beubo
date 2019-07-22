package beubo

import "testing"

func TestSetSetting(t *testing.T) {
	key := "key"
	value := "value"
	expected := key
	if result := setSetting(key, value); result != expected {
		t.Errorf("setSetting = %q, expected %q", result, expected)
	}
}
