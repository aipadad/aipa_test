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
 * file description:  blockchain utility
 * @Author: 
 * @Date:   2018-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	"github.com/aipadad/aipa/bpl"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/db"
	log "github.com/cihub/seelog"
)

var (
	//BlockHashPrefix prefix of block hash
	BlockHashPrefix = []byte("bh-")
	//BlockNumberPrefix prefix of block number
	BlockNumberPrefix = []byte("bn-")
	//LastBlockKey prefix of block key
	LastBlockKey = []byte("lb")
)

//HasBlock check block in db
func HasBlock(db *db.DBService, hash common.Hash) bool {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) != 0 {
		return true
	}

	return false
}

//GetBlock get block from db by hash
func GetBlock(db *db.DBService, hash common.Hash) *types.Block {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) == 0 {
		return nil
	}

	block := types.Block{}
	if err := bpl.Unmarshal(data, &block); err != nil {
		return nil
	}
	return &block
}

//GetBlockHashByNumber get block from db by number
func GetBlockHashByNumber(db *db.DBService, number uint64) common.Hash {
	hash, _ := db.Get(append(BlockNumberPrefix, common.NumberToBytes(number, 32)...))
	if len(hash) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(hash)
}

//GetLastBlock get lastest block from db
func GetLastBlock(db *db.DBService) *types.Block {
	data, _ := db.Get(LastBlockKey)
	if len(data) == 0 {
		return nil
	}

	return GetBlock(db, common.BytesToHash(data))
}

//WriteGenesisBlock write the first block in db
func WriteGenesisBlock(db *db.DBService, block *types.Block) error {
	if err := WriteBlock(db, block); err != nil {
		return err
	}

	return nil
}

func writeHead(db *db.DBService, block *types.Block) error {
	key := append(BlockNumberPrefix, common.NumberToBytes(block.GetNumber(), 32)...)
	err := db.Put(key, block.Hash().Bytes())
	if err != nil {
		return err
	}

	err = db.Put(LastBlockKey, block.Hash().Bytes())
	if err != nil {
		return err
	}
	return nil
}

//WriteBlock write block in db
func WriteBlock(db *db.DBService, block *types.Block) error {
	start := common.MeasureStart()
	key := append(BlockHashPrefix, block.Hash().Bytes()...)
	data, _ := bpl.Marshal(block)

	err := db.Put(key, data)
	if err != nil {
		return err
	}
	span := common.Elapsed(start)
	log.Infof("WriteBlock span:%v", span)
	return writeHead(db, block)
}
