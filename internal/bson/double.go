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
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/MangoDB-io/MangoDB/internal/util/lazyerrors"
)

type Double float64

func (d *Double) bsontype() {}

func (d *Double) ReadFrom(r *bufio.Reader) error {
	var bits uint64
	if err := binary.Read(r, binary.LittleEndian, &bits); err != nil {
		return lazyerrors.Errorf("bson.Double.ReadFrom (binary.Read): %w", err)
	}

	*d = Double(math.Float64frombits(bits))
	return nil
}

func (d Double) WriteTo(w *bufio.Writer) error {
	v, err := d.MarshalBinary()
	if err != nil {
		return lazyerrors.Errorf("bson.Double.WriteTo: %w", err)
	}

	_, err = w.Write(v)
	if err != nil {
		return lazyerrors.Errorf("bson.Double.WriteTo: %w", err)
	}

	return nil
}

func (d Double) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, math.Float64bits(float64(d)))

	return buf.Bytes(), nil
}

type doubleJSON struct {
	F float64 `json:"$f,string"`
}

func (d *Double) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	// TODO "Infinity", "-Infinity", "NaN"

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o doubleJSON
	if err := dec.Decode(&o); err != nil {
		return err
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Errorf("bson.Double.UnmarshalJSON: %s", err)
	}

	*d = Double(o.F)
	return nil
}

func (d Double) MarshalJSON() ([]byte, error) {
	// TODO "Infinity", "-Infinity", "NaN"

	return json.Marshal(doubleJSON{
		F: float64(d),
	})
}

// check interfaces
var (
	_ bsontype = (*Double)(nil)
)
