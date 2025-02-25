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
	"testing"

	"github.com/MangoDB-io/MangoDB/internal/types"
)

var binaryTestcases = []fuzzTestCase{{
	name: "foo",
	v: &Binary{
		Subtype: types.BinaryUser,
		B:       []byte("foo"),
	},
	b: []byte{0x03, 0x00, 0x00, 0x00, 0x80, 0x66, 0x6f, 0x6f},
	j: `{"$b":"Zm9v","s":128}`,
}, {
	name: "empty",
	v: &Binary{
		B: []byte{},
	},
	b: []byte{0x00, 0x00, 0x00, 0x00, 0x00},
	j: `{"$b":"","s":0}`,
}, {
	name: "invalid subtype",
	v: &Binary{
		Subtype: 0xff,
		B:       []byte{},
	},
	b: []byte{0x00, 0x00, 0x00, 0x00, 0xff},
	j: `{"$b":"","s":255}`,
}}

func TestBinary(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, binaryTestcases, func() bsontype { return new(Binary) })
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		testJSON(t, binaryTestcases, func() bsontype { return new(Binary) })
	})
}

func FuzzBinaryBinary(f *testing.F) {
	fuzzBinary(f, binaryTestcases, func() bsontype { return new(Binary) })
}

func FuzzBinaryJSON(f *testing.F) {
	fuzzJSON(f, binaryTestcases, func() bsontype { return new(Binary) })
}
