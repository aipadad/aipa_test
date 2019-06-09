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
 * file description:  voter role
 * @Author: 
 * @Date:   2018-08-23
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/aipadad/aipa/db"
	log "github.com/cihub/seelog"
)

const (
	//VoterObjectName is the table name of voter object
	VoterObjectName string = "voter"
)

// Voter is definition of voter
type Voter struct {
	Owner       string      `json:"owner"`
	Delegate    string      `json:"delegate"`
}

// CreateVoterRole is create voter role
func CreateVoterRole(ldb *db.DBService) error {
	ldb.AddObject(VoterObjectName)
	return nil
}

// SetVoterRole is common func to set role for voter
func SetVoterRole(ldb *db.DBService, accountName string, value *Voter) error {
	key := accountNameToKey(accountName)
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		log.Error("ROLE Marshal failed ", accountName, err)
		return err
	}
	return ldb.SetObject(VoterObjectName, key, string(jsonvalue))
}

// GetVoterRole is common func to get role for voter
func GetVoterRole(ldb *db.DBService, accountName string) (*Voter, error) {
	key := accountNameToKey(accountName)
	value, err := ldb.GetObject(VoterObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &Voter{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		log.Error("ROLE Unmarshal failed ", accountName, err)
		return nil, err
	}

	return res, nil
}
