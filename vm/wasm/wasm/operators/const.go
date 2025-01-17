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
	"github.com/aipadad/aipa/vm/wasm/wasm"
)

// DO NOT EDIT. This file define op code
var (
	I32Const = newOp(0x41, "i32.const", nil, wasm.ValueTypeI32)
	I64Const = newOp(0x42, "i64.const", nil, wasm.ValueTypeI64)
	F32Const = newOp(0x43, "f32.const", nil, wasm.ValueTypeF32)
	F64Const = newOp(0x44, "f64.const", nil, wasm.ValueTypeF64)
)
