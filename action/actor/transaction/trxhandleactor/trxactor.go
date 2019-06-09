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

package trxactor

import (
	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/aipadad/aipa/action/message"
	"github.com/aipadad/aipa/transaction"
)

//TrxActorPid trx actor pid
var TrxActorPid *actor.PID

var trxPool *transaction.TrxPool

//TrxActor trx actor props
type TrxActor struct {
	props *actor.Props
}

//ContructTrxActor new a trx actor
func ContructTrxActor() *TrxActor {
	return &TrxActor{}
}

//NewTrxActor spawn a named actor
func NewTrxActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructTrxActor() })

	var err error
	TrxActorPid, err = actor.SpawnNamed(props, "TrxActor")

	if err != nil {
		panic(log.Errorf("TrxActor SpawnNamed error: %v", err))
	} else {
		return TrxActorPid
	}
}

//SetTrxPool set trx pool
func SetTrxPool(pool *transaction.TrxPool) {
	trxPool = pool
}

func handleSystemMsg(context actor.Context) bool {
	switch context.Message().(type) {
	case *actor.Started:
		log.Error("TrxActor received started msg")
	case *actor.Stopping:
		log.Error("TrxActor received stopping msg")
	case *actor.Restart:
		log.Error("TrxActor received restart msg")
	case *actor.Restarting:
		log.Error("TrxActor received restarting msg")
	case *actor.Stop:
		log.Error("TrxActor received Stop msg")
	case *actor.Stopped:
		log.Error("TrxActor received Stopped msg")
	default:
		return false
	}

	return true
}

//Receive process message
func (t *TrxActor) Receive(context actor.Context) {

	if handleSystemMsg(context) {
		return
	}

	switch msg := context.Message().(type) {
	case *message.PushTrxForP2PReq:

		trxPool.HandleTransactionFromFront(context, msg.P2PTrx)

	case *message.ReceiveTrx:

		trxPool.HandleTransactionFromP2P(context, msg.P2PTrx)


	default:
		log.Errorf("trx actor: Unknown msg ", msg)
	}
}
