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
 * file description:  transaction history role
 * @Author: 
 * @Date:   2018-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/db"
)

//TransactionHistoryObjectName is transaction_history
const TransactionHistoryObjectName string = "transaction_history"

//TransactionHistory is to store the history of transaction
type TransactionHistory struct {
	TxHash    common.Hash `json:"trx_hash"`
	BlockHash common.Hash `json:"block_hash"`
}

//CreateTransactionHistoryObjectRole is creating a transaction history role
func CreateTransactionHistoryObjectRole(ldb *db.DBService) error {
	return nil
}

//AddTransactionHistoryRole is adding transaction history role
func AddTransactionHistoryRole(ldb *db.DBService, txhash common.Hash, blockhash common.Hash) error {
	value := &TransactionHistory{
		TxHash:    txhash,
		BlockHash: blockhash,
	}
	return setTransactionHistoryObjectRole(ldb, txhash, value)
}

func setTransactionHistoryObjectRole(ldb *db.DBService, txhash common.Hash, value *TransactionHistory) error {
	key := hashToKey(txhash)
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(TransactionHistoryObjectName, key, string(jsonvalue))
}

//GetTransactionHistoryRole is geting transaction history role
func GetTransactionHistoryRole(ldb *db.DBService, txhash common.Hash) (common.Hash, error) {
	history, err := getTransactionHistoryByHash(ldb, txhash)
	if err != nil {
		return common.Hash{}, err
	}
	return history.BlockHash, nil
}

func getTransactionHistoryByHash(ldb *db.DBService, hash common.Hash) (*TransactionHistory, error) {
	key := hashToKey(hash)
	value, err := ldb.GetObject(TransactionHistoryObjectName, key)
	if err != nil {
		return nil, err
	}
	res := &TransactionHistory{}
	json.Unmarshal([]byte(value), res)
	return res, err
}
