﻿// Copyright 2018~2022 The aipa Authors
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
Package types is definition of common type
 * file description:  transaction
 * @Author: 
 * @Date:   2018-12-07
 * @Last Modified by:
 * @Last Modified time:
 */
package types

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/aipadad/aipa/bpl"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/common/signature"
	log "github.com/cihub/seelog"
)

// Transaction define transaction struct for aipa protocol
type Transaction struct {
	Version     uint32
	CursorNum   uint64
	CursorLabel uint32
	Lifetime    uint64
	Sender      string // max length 21
	Contract    string // max length 21
	Method      string // max length 21
	Param       []byte
	SigAlg      uint32
	Signature   []byte
}
type P2PTransaction struct {
	Transaction *Transaction
	TTL  uint16
}
type BlockTransaction struct {
	Transaction     *Transaction
	ResourceReceipt  *ResourceReceipt
}

// ResourceReceipt
type ResourceReceipt struct {
	AccountName    string `json:"account_name"`
	SpaceTokenCost uint64 `json:"space_token_cost"`
	TimeTokenCost  uint64 `json:"time_token_cost"`
}

// HandledTransaction define transaction which is handled
type HandledTransaction struct {
	Transaction *Transaction
	DerivedTrx  []*DerivedTransaction
}

// DerivedTransaction define transaction which is derived from raw transaction
type DerivedTransaction struct {
	Transaction *Transaction
	DerivedTrx  []*DerivedTransaction
}

// Hash transaction hash
func (trx *Transaction) Hash() common.Hash {
	data, _ := bpl.Marshal(trx)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

// Hash transaction hash
func (trx *BlockTransaction) Hash() common.Hash {
	data, _ := bpl.Marshal(trx)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

// BasicTransaction define transaction struct for transaction signature
type BasicTransaction struct {
	Version     uint32
	CursorNum   uint64
	CursorLabel uint32
	Lifetime    uint64
	Sender      string
	Contract    string
	Method      string
	Param       []byte
	SigAlg      uint32
}

// VerifySignature verify signature
func (trx *Transaction) VerifySignature(pubkey []byte) bool {
	data, err := bpl.Marshal(BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       trx.Param,
		SigAlg:      trx.SigAlg,
	})

	if nil != err {
		return false
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	h.Write([]byte(hex.EncodeToString(config.GetChainID())))
	hash := h.Sum(nil)

	ok := signature.VerifySign(pubkey, hash, trx.Signature)

	if false == ok {
		log.Errorf("COMMON trx verify signature failed, hash %x, sender %s, pubkey %x", trx.Hash(), trx.Sender, pubkey)
	}

	return ok
}

// Sign sign a transaction with privkey
func (trx *Transaction) Sign(param []byte, privkey []byte) ([]byte, error) {
	data, err := bpl.Marshal(BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       param,
		SigAlg:      trx.SigAlg,
	})
	if nil != err {
		return []byte{}, err
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	h.Write([]byte(hex.EncodeToString(config.GetChainID())))
	hash := h.Sum(nil)
	signdata, err := signature.Sign(hash, privkey)

	return signdata, err
}
