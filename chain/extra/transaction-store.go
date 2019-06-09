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
 * @Date:   2018-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package txstore

import (
	"github.com/aipadad/aipa/chain"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/role"
)

var (
	//TrxBlockHashPrefix trx block hash prefix
	TrxBlockHashPrefix = []byte("txbh-")
)

//TransactionStore transaction store
type TransactionStore struct {
	roleIntf role.RoleInterface
	bc       chain.BlockChainInterface
}

//NewTransactionStore new a transaction store
func NewTransactionStore(bc chain.BlockChainInterface, roleIntf role.RoleInterface) *TransactionStore {
	ts := &TransactionStore{
		roleIntf: roleIntf,
		bc:       bc,
	}
	bc.RegisterHandledBlockCallback(ts.ReceiveHandledBlock)
	return ts
}

//GetTransaction get trx from block
func (t *TransactionStore) GetTransaction(txhash common.Hash) *types.Transaction {
	blockHash, err := t.roleIntf.GetTransactionHistory(txhash)
	if err != nil {
		return nil
	}

	block := t.bc.GetBlockByHash(blockHash)
	if block == nil {
		return nil
	}

	return block.GetTransactionByHash(txhash)
}

func (t *TransactionStore) addTx(txhash common.Hash, blockhash common.Hash) error {
	return t.roleIntf.AddTransactionHistory(txhash, blockhash)
}

func (t *TransactionStore) delTx(txhash common.Hash) error {
	return nil
}

//ReceiveHandledBlock receive a block
func (t *TransactionStore) ReceiveHandledBlock(block *types.Block) {
	blockHash := block.Hash()

	for _, tx := range block.Transactions {
		txHash := tx.Hash()
		t.addTx(txHash, blockHash)
	}
}

//RemoveBlock remove block
func (t *TransactionStore) RemoveBlock(block *types.Block) {
	for _, tx := range block.Transactions {
		txHash := tx.Hash()
		t.delTx(txHash)
	}
}
