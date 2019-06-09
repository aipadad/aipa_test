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
 * file description:  transaction actor
 * @Author:
 * @Date:   2018-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package trxprehandleactor

import (
	log "github.com/cihub/seelog"
	
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/aipadad/aipa/action/env"
	"github.com/aipadad/aipa/action/message"
	aipaErr "github.com/aipadad/aipa/common/errors"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/context"
	"github.com/aipadad/aipa/transaction"
	"github.com/aipadad/aipa/config"
)

//TrxPreHandleActorPid trx actor pid
var TrxPreHandleActorPid *actor.PID

var trxActorPid *actor.PID


const maxConcurrency = 100

var trxPool *transaction.TrxPool

var actorEnv *env.ActorEnv

var protocolInterface context.ProtocolInterface

//TrxActor trx actor props
// type TrxActor struct {
// 	props *actor.Props
// }

// //ContructTrxActor new a trx actor
// func ContructTrxActor() *TrxActor {
// 	return &TrxActor{}
// }

func handleSystemMsg(context actor.Context) bool {
	switch context.Message().(type) {
	case *actor.Started:
		log.Error("TrxPreHandleActor received started msg")
	case *actor.Stopping:
		log.Error("TrxPreHandleActor received stopping msg")
	case *actor.Restart:
		log.Error("TrxPreHandleActor received restart msg")
	case *actor.Restarting:
		log.Error("TrxPreHandleActor received restarting msg")
	case *actor.Stop:
		log.Error("TrxPreHandleActor received Stop msg")
	case *actor.Stopped:
		log.Error("TrxPreHandleActor received Stopped msg")
		
	default:
		return false
	}

	return true
}

func preHandleCommon(trx *types.Transaction) (bool, aipaErr.ErrCode, bool) {
	putInCache := false
	if checkResult, err := trxPool.CheckTransactionBaseCondition(trx); true != checkResult {
		return false, err, putInCache
	}

	if false == actorEnv.Protocol.GetBlockSyncState() {
		distance := actorEnv.Protocol.GetBlockSyncDistance()
		log.Errorf("TRX rcv trx when block is syncing, trx %x, distance %v", trx.Hash(), distance)
		if distance > config.DEFAULT_MAX_SYNC_DISTANCE_PUT_TRX_IN_CACHE {
			return false, aipaErr.ErrTrxBlockSyncingError, putInCache
		} else {
			putInCache = true
		}
	}

	sender, err := actorEnv.RoleIntf.GetAccount(trx.Sender)
	if nil != err {
		return false, aipaErr.ErrTrxAccountError, putInCache
	}

	if !trx.VerifySignature(sender.PublicKey) {
		return false, aipaErr.ErrTrxSignError, putInCache
	}

	return true, aipaErr.ErrNoError, putInCache
}
func initP2PTrxMsg(msg *message.PushTrxReq) (msgp *message.PushTrxForP2PReq) {
	//set trx TTL
	var TTL uint16
	switch  actorEnv.RoleIntf.IsMyselfDelegate() {
	case true:
		TTL = config.TRX_IN_TTL
	case false:
		TTL = config.TRX_OUT_TTL
	}
	var p2pTrx types.P2PTransaction
	p2pTrx.Transaction = msg.Trx
	p2pTrx.TTL = TTL
	msgp = &message.PushTrxForP2PReq{P2PTrx: &p2pTrx}
	return msgp
}

func preHandlePushTrxReq(msg *message.PushTrxReq, ctx actor.Context) {

	preHandleResult, err, putInCache := preHandleCommon(msg.Trx)
	
	if !preHandleResult {			
		log.Errorf("TRX pre handle trx from front failed, trx %x", msg.Trx.Hash())
		ctx.Respond(err)
	} else {
		msgP2P := initP2PTrxMsg(msg)
		if putInCache || trxPool.IsCacheEmpty() == false {
			trxPool.AddTransactionToCache(msgP2P)
		} else {
		trxActorPid.Tell(msgP2P)
		}
		ctx.Respond(aipaErr.ErrNoError)
	}
}

func preHandleReceiveTrx(msg *message.ReceiveTrx, ctx actor.Context) {

	preHandleResult, _, putInCache := preHandleCommon(msg.P2PTrx.Transaction)
	
	if preHandleResult {
		if putInCache || trxPool.IsCacheEmpty() == false {
			trxPool.AddTransactionToCache(msg)
		} else {
		trxActorPid.Tell(msg)
		}
	} else {
		if actorEnv.RoleIntf.IsMyselfDelegate() == true{
			log.Info("TRX pre handle trx from producer node failed, trx %x", msg.P2PTrx.Transaction.Hash())
		}else{
			log.Errorf("TRX pre handle trx from service node failed, trx %x", msg.P2PTrx.Transaction.Hash())
		}

	}
}

func doWork(ctx actor.Context) {

	if handleSystemMsg(ctx) {
		return
	}

	switch msg := ctx.Message().(type) {
	case *message.PushTrxReq:

		log.Infof("rcv trx %x in PushTrxReq\n", msg.Trx.Hash())

		preHandlePushTrxReq(msg, ctx)
		
	case *message.ReceiveTrx:

		log.Infof("rcv trx %x in ReceiveTrx\n", msg.P2PTrx.Transaction.Hash())

		preHandleReceiveTrx(msg, ctx)		

	default:
		log.Errorf("trx pool actor: Unknown msg ", msg)
	}

}

//NewTrxPreHandleActor spawn a named actor
func NewTrxPreHandleActor(env *env.ActorEnv) *actor.PID {

	actorEnv = env

	TrxPreHandleActorPid := actor.Spawn(router.NewRoundRobinPool(maxConcurrency).WithFunc(doWork))

	return TrxPreHandleActorPid
}

//SetTrxPool set trx pool
func SetTrxPool(pool *transaction.TrxPool) {
	trxPool = pool
}

func SetTrxActor(trxactorPid *actor.PID) {
	trxActorPid = trxactorPid
}

