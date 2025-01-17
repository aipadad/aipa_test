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

// Package operators provides all operators used by WebAssembly bytecode,
// together with their parameter and return type(s).
package operators

import (
	"fmt"
	"github.com/aipadad/aipa/vm/wasm/wasm"
	log "github.com/cihub/seelog"
)

var (
	ops      [256]Op // an array of Op values mapped by wasm opcodes, used by New().
	noReturn = wasm.ValueType(wasm.BlockTypeEmpty)
)

// Op describes a WASM operator.
type Op struct {
	Code byte   // The single-byte opcode
	Name string // The name of the operator

	// Whether this operator is polymorphic.
	// A polymorphic operator has a variable arity. call, call_indirect, and
	// drop are examples of polymorphic operators.
	Polymorphic bool
	Args        []wasm.ValueType // an array of value types used by the operator as arguments, is nil for polymorphic operators
	Returns     wasm.ValueType   // the value returned (pushed) by the operator, is 0 for polymorphic operators
}

//IsValid check name
func (o Op) IsValid() bool {
	return o.Name != ""
}

func newOp(code byte, name string, args []wasm.ValueType, returns wasm.ValueType) byte {
	if ops[code].IsValid() {
		panic(log.Errorf("Opcode %#x is already assigned to %s", code, ops[code].Name))
	}

	op := Op{
		Code:        code,
		Name:        name,
		Polymorphic: false,
		Args:        args,
		Returns:     returns,
	}
	ops[code] = op
	return code
}

func newPolymorphicOp(code byte, name string) byte {
	if ops[code].IsValid() {
		panic(log.Errorf("Opcode %#x is already assigned to %s", code, ops[code].Name))
	}

	op := Op{
		Code:        code,
		Name:        name,
		Polymorphic: true,
	}
	ops[code] = op
	return code
}

//InvalidOpcodeError op code error type
type InvalidOpcodeError byte

func (e InvalidOpcodeError) Error() string {
	return fmt.Sprintf("Invalid opcode: %#x", byte(e))
}

// New returns the Op object for a valid given opcode.
// If code is invalid, an ErrInvalidOpcode is returned.
func New(code byte) (Op, error) {
	var op Op
	if int(code) >= len(ops) {
		return op, InvalidOpcodeError(code)
	}

	op = ops[code]
	if !op.IsValid() {
		return op, InvalidOpcodeError(code)
	}
	return op, nil
}
