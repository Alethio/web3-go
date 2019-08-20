package thelper

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func Load(t *testing.T, fn string, out interface{}) {
	t.Helper()

	raw := LoadFile(t, fn)

	err := json.Unmarshal(raw, out)
	if err != nil {
		t.Fatal(err)
	}
}

func Save(t *testing.T, fn string, out interface{}) {
	t.Helper()
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(fn, b, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func LoadFile(t *testing.T, fn string) []byte {
	t.Helper()

	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}

	return raw
}

func SaveOnUpdate(t *testing.T, u *bool, fn string, g interface{}) {
	t.Helper()

	if *u {
		Save(t, fn, g)
	}
}
