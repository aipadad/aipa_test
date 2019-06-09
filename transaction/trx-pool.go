﻿//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with aipa.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  trx pool
 * @Author: 
 * @Date:   2018-12-15
 * @Last Modified by:
 * @Last Modified time:
 */

package transaction

import (
	"fmt"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/aipadad/aipa/action/message"
	"github.com/aipadad/aipa/bpl"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/contract"
	"github.com/aipadad/aipa/db"
	"github.com/aipadad/aipa/role"

	"crypto/sha256"
	"encoding/hex"

	aipaErr "github.com/aipadad/aipa/common/errors"
	"github.com/aipadad/crypto-go/crypto"
	log "github.com/cihub/seelog"
)

var (
	trxExpirationCheckInterval = 60 * time.Second // Time interval for check expiration pending transactions
)

// TrxPoolInst is local var of TrxPool
var TrxPoolInst *TrxPool

// TrxPool is definition of trx pool
type TrxPool struct {
	pending     map[common.Hash]*types.Transaction
	roleIntf    role.RoleInterface
	netActorPid *actor.PID

	dbInst *db.DBService
	mu     sync.RWMutex
	quit   chan struct{}
}

// InitTrxPool is init trx pool process when system start
func InitTrxPool(dbInstance *db.DBService, roleIntf role.RoleInterface, nc contract.NativeContractInterface, netActorPid *actor.PID) *TrxPool {

	TrxPoolInst := &TrxPool{
		pending:     make(map[common.Hash]*types.Transaction),
		roleIntf:    roleIntf,
		netActorPid: netActorPid,
		dbInst:      dbInstance,

		quit: make(chan struct{}),
	}

	CreateTrxApplyService(roleIntf, nc)

	go TrxPoolInst.expirationCheckLoop()

	return TrxPoolInst
}

func (trxPool *TrxPool) expirationCheckLoop() {

	expire := time.NewTicker(trxExpirationCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:

			var currentTime = common.Now()
			for trxHash := range trxPool.pending {
				if currentTime >= (trxPool.pending[trxHash].Lifetime) {
					log.Info("remove expirate trx, hash is: ", trxHash, "curtime", currentTime, "lifeTime", trxPool.pending[trxHash].Lifetime)
					trxPool.RemoveSingleTransactionbyHash(trxHash)
				}
			}

		case <-trxPool.quit:
			return
		}
	}
}

func (trxPool *TrxPool) isTransactionExist(trx *types.Transaction) {
	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	_, ok := trxPool[trx.Hash()]
	return ok
}

func (trxPool *TrxPool) addTransaction(trx *types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	trxHash := trx.Hash()
	trxPool.pending[trxHash] = trx
}

// Stop is processing when system stop
func (trxPool *TrxPool) Stop() {

	close(trxPool.quit)

	log.Info("Transaction pool stopped")
}

// CheckTransactionBaseCondition is checking trx
func (trxPool *TrxPool) CheckTransactionBaseCondition(trx *types.Transaction) (bool, aipaErr.ErrCode) {
	if isTransactionExist(trx) {
		return false, aipa.ErrTrxAlreadyInPool
	}
	if config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(trxPool.pending)) {
		log.Errorf("trx %x pending num over", trx.Hash())
		return false, aipaErr.ErrTrxPendingNumLimit
	}

	if !trxPool.VerifySignature(trx) {
		return false, aipaErr.ErrTrxSignError
	}

	return true, aipaErr.ErrNoError
}

// HandleTransactionCommon is processing trx
func (trxPool *TrxPool) HandleTransactionCommon(context actor.Context, trx *types.Transaction) {

	if checkResult, err := trxPool.CheckTransactionBaseCondition(trx); true != checkResult {
		context.Respond(err)
		return
	}

	result, err, _ := trxApplyServiceInst.ApplyTransaction(trx)
	if !result {
		context.Respond(err)
		return
	}

	trxPool.addTransaction(trx)

	notify := &message.NotifyTrx{
		Trx: trx,
	}
	trxPool.netActorPid.Tell(notify)
}

// HandleTransactionFromFront is handling trx from front
func (trxPool *TrxPool) HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {
	log.Infof("rcv trx %x from front,sender %v, contract %v, method %v", trx.Hash(), trx.Sender, trx.Contract, trx.Method)
	trxPool.HandleTransactionCommon(context, trx)
}

// HandleTransactionFromP2P is handling trx from P2P
func (trxPool *TrxPool) HandleTransactionFromP2P(context actor.Context, trx *types.Transaction) {
	log.Tracef("rcv trx %x from P2P,sender %v, contract %v method %v", trx.Hash(), trx.Sender, trx.Contract, trx.Method)
	trxPool.HandleTransactionCommon(context, trx)
}

// GetAllPendingTransactions is interface to get all pending trxs in trx pool
func (trxPool *TrxPool) GetAllPendingTransactions(context actor.Context) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	rsp := &message.GetAllPendingTrxRsp{}
	for trxHash := range trxPool.pending {
		rsp.Trxs = append(rsp.Trxs, trxPool.pending[trxHash])
	}

	context.Respond(rsp)
}

// RemoveTransactions is interface to remove trxs in trx pool
func (trxPool *TrxPool) RemoveTransactions(trxs []*types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	for _, trx := range trxs {
		delete(trxPool.pending, trx.Hash())
	}
}

// RemoveSingleTransaction is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransaction(trx *types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	delete(trxPool.pending, trx.Hash())
}

// RemoveSingleTransactionbyHash is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransactionbyHash(trxHash common.Hash) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	delete(trxPool.pending, trxHash)
}

func (trxPool *TrxPool) getPubKey(accountName string) ([]byte, error) {

	account, err := trxPool.roleIntf.GetAccount(accountName)
	if nil != err {
		return nil, fmt.Errorf("get account failed")
	}

	return account.PublicKey, nil
}

// VerifySignature is verify signature from trx whether it is valid
func (trxPool *TrxPool) VerifySignature(trx *types.Transaction) bool {

	trxToVerify := &types.BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       trx.Param,
		SigAlg:      trx.SigAlg,
	}

	serializeData, err := bpl.Marshal(trxToVerify)
	if nil != err {
		return false
	}

	senderPubKey, err := trxPool.getPubKey(trx.Sender)
	if nil != err {
		log.Errorf("trx %x get pub key error", trx.Hash())
		return false
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(serializeData)))
	h.Write([]byte(hex.EncodeToString(config.GetChainID())))
	hashData := h.Sum(nil)

	verifyResult := crypto.VerifySign(senderPubKey, hashData, trx.Signature)

	if false == verifyResult {
		log.Errorf("trx %x verify signature failed, sender %v, pubkey %v", trx.Hash(), trx.Sender, senderPubKey)
	}

	return verifyResult
}
