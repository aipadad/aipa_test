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
 * file description:  producer actor
 * @Author: 
 * @Date:   2018-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package produceractor

import (
	"time"
	"unsafe"

	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/aipadad/aipa/action/env"
	"github.com/aipadad/aipa/bpl"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/db"
	"github.com/aipadad/aipa/db/platform/codedb"
	"github.com/aipadad/aipa/producer"
	"github.com/aipadad/aipa/role"
	"github.com/aipadad/aipa/version"
)

// ProducerActor is to define actor for producer
type ProducerActor struct {
	roleIntf         role.RoleInterface
	ins              producer.ReporterRepo
	db               *db.DBService
	pendingTxSession *codedb.UndoSession
}

// NewProducerActor is to create actor for producer
func NewProducerActor(env *env.ActorEnv) *actor.PID {

	ins := producer.New(env.Chain, env.RoleIntf, env.Protocol)
	props := actor.FromProducer(func() actor.Actor {
		return &ProducerActor{
			env.RoleIntf,
			ins,
			env.Db,
			env.PendingTxSession,
		}
	})

	pid, err := actor.SpawnNamed(props, "ProducerActor")

	if err != nil {
		return nil
	}

	return pid
}

func (p *ProducerActor) handleSystemMsg(context actor.Context) bool {
	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Infof("ProducerActor received started msg: %s", msg)

		context.SetReceiveTimeout(time.Duration(config.PRODUCER_TIME_OUT) * time.Millisecond)

	case *actor.ReceiveTimeout:
		elapse := p.working()
		context.SetReceiveTimeout(time.Duration(elapse) * time.Millisecond)

	case *actor.Stopping:
		log.Info("ProducerActor received stopping msg")

	case *actor.Restart:
		log.Info("ProducerActor received restart msg")

	case *actor.Restarting:
		log.Info("ProducerActor received restarting msg")

	case *actor.Stop:
		log.Info("ProducerActor received Stop msg")

	case *actor.Stopped:
		log.Info("ProducerActor received Stopped msg")

	default:
		return false
	}

	return true
}

// Receive is to receive and handle message
func (p *ProducerActor) Receive(context actor.Context) {

	if p.handleSystemMsg(context) {
		return
	}

	log.Error("ProducerActor received Unknown msg")
}
func (p *ProducerActor) working() uint32 {

	if p.ins.IsReady() {
		p.db.Lock()
		defer p.db.UnLock()
		log.Infof("begin to producer block ")
		p.pendingTxSession = p.db.GetSession()
		if p.pendingTxSession != nil {
			log.Infof("p.pendingTxSession need to reset ")
			p.db.ResetSession()
		}
		log.Infof("begin session......... ")
		p.pendingTxSession = p.db.BeginUndo(config.PRIMARY_TRX_SESSION)

		trxs := GetAllPendingTrx()
		log.Info("get trx times", common.Elapsed(start))
		block := &types.Block{}
		data, _ := bpl.Marshal(block)
		pendingBlockSize := uint32(len(data))
		coreStat, err := p.roleIntf.GetCoreState()
		if err != nil {
			log.Info("GetGlobalPropertyRole failed")
			p.db.ResetSession()
			return config.PRODUCER_TIME_OUT
		}
		chainStat, err := p.roleIntf.GetChainState()
		if err != nil {
			log.Error("PRODUCER GetChainState failed,begin rollback", err)
			p.db.ResetSession()
			return config.PRODUCER_TIME_OUT
		}
		version := version.GetVersionNumByBlockNum(chainStat.LastBlockNum + 1)
		var pendingBlockTrx = []*types.BlockTransaction{}
		var removeTrx = []*types.Transaction{}

		err = p.roleIntf.OnBlock(p.roleIntf.GetMySelf())
		if err != nil {
			return config.PRODUCER_TIME_OUT
		}

		for _, trx := range trxs {

			if uint64(common.Elapsed(start)) > config.DEFAULT_BLOCK_TIME_LIMIT {
				p.db.ResetSubSession()
				log.Error("PRODUCER block reach elapse time ", common.Elapsed(start))
				break
			}

			p.db.BeginUndo(config.SUB_TRX_SESSION)
			applyStart := common.MeasureStart()
			pass, _, resReceipt := verifyTransactions(trx, version)
			if pass == false {
				log.Info("ApplyTransaction failed")
				p.db.ResetSubSession()
				removeTrx = append(removeTrx, trx)
				continue
			}
			log.Info("apply start elapse", common.Elapsed(applyStart))
			data, _ := bpl.Marshal(trx)
			pendingBlockSize += uint32(unsafe.Sizeof(data))
			log.Info("pendingBlockSize ", pendingBlockSize)

			if pendingBlockSize > coreStat.Config.MaxBlockSize {
				p.db.ResetSubSession()
				log.Info("Warning pending block size reach MaxBlockSize")
				pendingTrx = append(pendingTrx, dtag)
				continue
			}
			p.db.Squash()

			pendingBlockTrx = append(pendingBlockTrx, dtag)
			log.Info("pack apply elapse", common.Elapsed(applyStart))
		}

		removeTransaction(removeTrx)
		block = p.ins.Woker(pendingBlockTrx)
		p.db.ResetSession()
		trxs = nil
		if block != nil {
			log.Infof("Generate block: hash: %x, delegate: %s, number:%v, trxn:%v,blockTime:%s\n", block.Hash(), block.Header.Delegate, block.GetNumber(), len(block.Transactions), time.Unix(int64(block.Header.Timestamp), 0))
			ApplyBlock(block)
		}
		return p.ins.CalcNextReportTime(block)
	}
	return config.PRODUCER_TIME_OUT
}
