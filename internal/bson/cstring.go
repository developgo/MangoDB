// Copyright 2021 Baltoro OÜ.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bson

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/MangoDB-io/MangoDB/internal/util/lazyerrors"
)

type CString string

func (cstr *CString) bsontype() {}

func (cstr *CString) ReadFrom(r *bufio.Reader) error {
	b, err := r.ReadBytes(0)
	if err != nil {
		return lazyerrors.Errorf("bson.CString.ReadFrom: %w", err)
	}

	*cstr = CString(b[:len(b)-1])
	return nil
}

func (cstr CString) WriteTo(w *bufio.Writer) error {
	v, err := cstr.MarshalBinary()
	if err != nil {
		return lazyerrors.Errorf("bson.CString.WriteTo: %w", err)
	}

	_, err = w.Write(v)
	if err != nil {
		return lazyerrors.Errorf("bson.CString.WriteTo: %w", err)
	}

	return nil
}

func (cstr CString) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(cstr)+1)
	copy(b, cstr)
	return b, nil
}

type cstringJSON struct {
	CString string `json:"$c"`
}

func (cstr *CString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o cstringJSON
	if err := dec.Decode(&o); err != nil {
		return err
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Errorf("bson.CString.UnmarshalJSON: %s", err)
	}

	*cstr = CString(o.CString)
	return nil
}

func (cstr CString) MarshalJSON() ([]byte, error) {
	return json.Marshal(cstringJSON{
		CString: string(cstr),
	})
}

// check interfaces
var (
	_ bsontype = (*CString)(nil)
)
