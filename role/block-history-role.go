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
 * file description:  block history role
 * @Author: 
 * @Date:   2018-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	"strconv"

	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/db"
)

//BlockHistoryObjectName is definition of block history object name
const BlockHistoryObjectName string = "block_history"

// BlockHistory is definition of block history
type BlockHistory struct {
	BlockNumber uint64      `json:"block_number"`
	BlockHash   common.Hash `json:"block_hash"`
}

func blockNumberToKey(blockNumber uint64) string {
	id := blockNumber & 0xFFFF
	key := strconv.Itoa(int(id))
	return key
}

// CreateBlockHistoryRole is to init block history
func CreateBlockHistoryRole(ldb *db.DBService) error {
	return nil
}

// SetBlockHistoryRole is to save block history
func SetBlockHistoryRole(ldb *db.DBService, blockNumber uint64, blockHash common.Hash) error {
	key := blockNumberToKey(blockNumber)
	value := &BlockHistory{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
	}
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ldb.SetObject(BlockHistoryObjectName, key, string(jsonvalue))
}

// GetBlockHistoryRole is to get block history
func GetBlockHistoryRole(ldb *db.DBService, blockNumber uint64) (*BlockHistory, error) {
	key := blockNumberToKey(blockNumber)
	value, err := ldb.GetObject(BlockHistoryObjectName, key)
	if err != nil {
		return nil, err
	}
	res := &BlockHistory{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
