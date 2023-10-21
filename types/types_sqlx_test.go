// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"testing"
)

type fakeJson struct {
	json.RawMessage
}

func (j *fakeJson) MarshalBinary() (res []byte, err error) {
	if j == nil {
		return []byte{}, nil
	}
	var m json.RawMessage
	if len(j.RawMessage) == 0 {
		j.RawMessage = json.RawMessage("{}")
	}
	err, res = json.Unmarshal([]byte(j.RawMessage), &m), m
	return
}

func (j *fakeJson) UnmarshalBinary(data []byte) (res *fakeJson, err error) {
	if j == nil {
		res = &fakeJson{}
	} else {
		res = j
	}
	err = res.RawMessage.UnmarshalJSON(data)
	return
}

func TestBinary(t *testing.T) {
	j := Binary[*fakeJson]{
		Data: &fakeJson{
			RawMessage: json.RawMessage(`{"foo": 1, "bar": 2}`),
		},
	}
	v, err := j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}
	err = j.Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}

	j = Binary[*fakeJson]{
		Data: &fakeJson{
			RawMessage: json.RawMessage(`{"foo": 1, invalid, false}`),
		},
	}
	_, err = j.Value()
	if err == nil {
		t.Errorf("Was expecting invalid json to fail!")
	}

	j = Binary[*fakeJson]{
		Data: &fakeJson{
			RawMessage: json.RawMessage(""),
		},
	}
	v, err = j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}

	err = j.Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}

	j = Binary[*fakeJson]{
		Data: &fakeJson{
			RawMessage: nil,
		},
	}
	v, err = j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}

	err = (&j).Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}
}

func TestNullBinary(t *testing.T) {
	j := NullBinary[*fakeJson]{}
	err := j.Scan(`{"foo": 1, "bar": 2}`)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}
	v, err := j.Value()
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}
	err = (&j).Scan(v)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}

	j = NullBinary[*fakeJson]{}
	err = j.Scan(nil)
	if err != nil {
		t.Errorf("Was not expecting an error: %s", err)
	}
	if j.Valid != false {
		t.Errorf("Expected valid to be false, but got true")
	}
}
