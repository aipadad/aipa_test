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
 * file description: code database test
 * @Author: 
 * @Date:   2018-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package role

import (
	//	"encoding/json"
	log "github.com/cihub/seelog"
	"math/big"
	"testing"

	"github.com/aipadad/aipa/db"
)

func TestDelegateVotes_writedb(t *testing.T) {
	ins := db.NewDbService("./file2", "./file2/db.db", "")
	err := CreateDelegateVotesRole(ins)
	if err != nil {
		log.Error(err)
	}
	value := &DelegateVotes{
		OwnerAccount: "nodepad",
		Serve: Serve{
			Votes:          1,
			Position:       big.NewInt(2),
			TermUpdateTime: big.NewInt(2),
			TermFinishTime: big.NewInt(2),
		},
	}
	err = SetDelegateVotesRole(ins, value.OwnerAccount, value)
	if err != nil {
		log.Error("SetDelegateVotesRole", err)
	}

	value, err = GetDelegateVotesRole(ins, value.OwnerAccount)
	if err != nil {
		log.Error("GetDelegateVotesRole", err)
	}
	log.Info(value)

	value, err = GetDelegateVotesRoleByVote(ins, value.Serve.Votes)
	if err != nil {
		log.Error("GetDelegateVotesRoleByVote", err)
	}
	log.Info(value)

	value, err = GetDelegateVotesRoleByFinishTime(ins, value.Serve.TermFinishTime)
	if err != nil {
		log.Error("GetDelegateVotesRoleByFinishTime", err)
	}
	log.Info(value)

	values, nerr := GetAllDelegateVotesRole(ins)
	if nerr != nil {
		log.Error("GetAllDelegateVotes", nerr)
	}
	log.Info(len(values))

	svotes, nerr := GetAllSortVotesDelegates(ins)
	if nerr != nil {
		log.Error("GetAllSortVotesDelegates", nerr)
	}
	log.Info(len(svotes))
	tvotes, nerr := GetAllSortFinishTimeDelegates(ins)
	if nerr != nil {
		log.Error("GetAllSortFinishTimeDelegates", nerr)
	}
	log.Info(len(tvotes))
}
