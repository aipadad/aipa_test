// Copyright 2018~2022 The aipa Authors
// This file is part of the aipa Chain library.
// Created by  Team of aipa.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with aipa.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2018 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package operators

import (
	"reflect"
	"testing"

	"github.com/aipadad/aipa/vm/wasm/wasm"
)

func TestNewConversionOp(t *testing.T) {
	origOps := ops
	defer func() {
		ops = origOps
	}()

	ops = [256]Op{}
	testCases := []struct {
		name    string
		args    []wasm.ValueType
		returns wasm.ValueType
	}{
		{"i32.wrap/i64", []wasm.ValueType{wasm.ValueTypeI64}, wasm.ValueTypeI32},
		{"i32.trunc_s/f32", []wasm.ValueType{wasm.ValueTypeF32}, wasm.ValueTypeI32},
	}

	for i, testCase := range testCases {
		op, err := New(newConversionOp(byte(i), testCase.name))
		if err != nil {
			t.Fatalf("%s: unexpected error from New: %v", testCase.name, err)
		}

		if !reflect.DeepEqual(op.Args, testCase.args) {
			t.Fatalf("%s: unexpected param types: got=%v, want=%v", testCase.name, op.Args, testCase.args)
		}

		if op.Returns != testCase.returns {
			t.Fatalf("%s: unexpected return type: got=%v, want=%v", testCase.name, op.Returns, testCase.returns)
		}
	}
}
