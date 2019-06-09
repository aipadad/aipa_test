// Copyright 2018~2022 The aipa Authors
// This file is part of the aipa Chain library.
// Created by  Team of aipa.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with aipa.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  context definition
 * @Author: 
 * @Date:   2018-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/role"
)

//Context for contracts
type Context struct {
	RoleIntf   role.RoleInterface
	Trx        *types.Transaction
}

//GetTrxParam for contracts
func (ctx *Context) GetTrxParam() []byte {
	return ctx.Trx.Param
}

//GetTrxParamSize for contracts
func (ctx *Context) GetTrxParamSize() uint32 {
	size := len(ctx.Trx.Param)
	return uint32(size)
}
