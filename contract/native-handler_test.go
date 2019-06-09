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
	 "regexp"
	 "github.com/aipadad/aipa/config"
	 "github.com/aipadad/aipa/common"

 )
 
 func TestAnalyzeName(t *testing.T) {
	nc := &NativeContract{
		
	}
	nc.AccountReg = regexp.MustCompile(config.ACCOUNT_NAME_REGEXP)
	nc.ContractReg = regexp.MustCompile(config.CONTRACT_NAME_REGEXP)
	nc.ExContractReg = regexp.MustCompile(config.EX_CONTRACT_NAME_REGEXP)

	type1, account1 := nc.analyzeName("aaa@bbb")
	if (type1 != common.NameTypeExContract || account1 != "bbb") {
		t.Error("aaa@bbb")
	}

	type2, account2 := nc.analyzeName("aaabbb")
	if (type2 != common.NameTypeAccount || account2 != "aaabbb") {
		t.Error("aaabbb")
	}
 }
 