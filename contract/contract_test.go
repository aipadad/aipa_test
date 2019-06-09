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
	//"fmt"
	"testing"

	"github.com/aipadad/aipa/contract/msgpack"
	log "github.com/cihub/seelog"
)

func TestTransfer(t *testing.T) {
	type transferparam struct {
		From  string
		To    string
		Value uint64
	}

	param := transferparam{
		From:  "delegate1",
		To:    "delegate2",
		Value: 100,
	}
	data, _ := msgpack.Marshal(param)
	log.Infof("transfer struct: %v, msgpack: %x\n", param, data)
}

func TestNewAccount(t *testing.T) {
	type newaccountparam struct {
		Name   string
		Pubkey string
	}

	param := newaccountparam{
		Name:   "testuser",
		Pubkey: "7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF",
	}
	data, _ := msgpack.Marshal(param)
	log.Infof("transfer struct: %v, msgpack: %x\n", param, data)
}

func TestGrantCredit(t *testing.T) {
	type GrantCreditParam struct {
		Name    string `json:"name"`
		Spender string `json:"spender"`
		Limit   uint64 `json:"limit"`
	}

	type CancelCreditParam struct {
		Name    string `json:"name"`
		Spender string `json:"spender"`
	}

	type TransferFromParam struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Value uint64 `json:"value"`
	}

	param := GrantCreditParam{
		Name:    "alice",
		Spender: "bob",
		Limit:   100,
	}
	data, _ := msgpack.Marshal(param)
	log.Infof("grant credit struct: %v, msgpack: %x\n", param, data)

	param1 := CancelCreditParam{
		Name:    "alice",
		Spender: "bob",
	}
	data, _ = msgpack.Marshal(param1)
	log.Infof("cancel credit struct: %v, msgpack: %x\n", param1, data)

	param2 := TransferFromParam{
		From:  "alice",
		To:    "toliman",
		Value: 150,
	}
	data, _ = msgpack.Marshal(param2)
	log.Infof("transfer from credit struct: %v, msgpack: %x\n", param2, data)

}
