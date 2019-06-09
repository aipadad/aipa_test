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
 * file description:  chain actor
 * @Author:
 * @Date:   2018-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package chainactor

import (
	"fmt"

	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/aipadad/aipa/action/env"
	"github.com/aipadad/aipa/action/message"
	"github.com/aipadad/aipa/chain"
	"github.com/aipadad/aipa/common/types"
)

//ChainActorPid is chain actor pid
var ChainActorPid *actor.PID
var actorEnv *env.ActorEnv
var trxPoolActorPid *actor.PID
var NetActorPid *actor.PID

//ChainActor is actor props
type ChainActor struct {
	props *actor.Props
}

//ContructChainActor new and actor
func ContructChainActor() *ChainActor {
	return &ChainActor{}
}

//SetTrxActorPid set trx actor pid
func SetTrxPoolActorPid(tpid *actor.PID) {
	trxPoolActorPid = tpid
}

//SetNetActorPid set trx actor pid
func SetNetActorPid(pid *actor.PID) {
	NetActorPid = pid
}

//NewChainActor spawn a named actor
func NewChainActor(env *env.ActorEnv) *actor.PID {
	var err error

	props := actor.FromProducer(func() actor.Actor { return ContructChainActor() })

	ChainActorPid, err = actor.SpawnNamed(props, "ChainActor")
	actorEnv = env

	if err != nil {
		panic(log.Errorf("ChainActor SpawnNamed error: %v", err))
	} else {
		return ChainActorPid
	}
}

func handleSystemMsg(context actor.Context) bool {

	switch context.Message().(type) {
	case *actor.Started:
		log.Info("BlockActor received started msg")
	case *actor.Stopping:
		log.Info("BlockActor received stopping msg")
	case *actor.Restart:
		log.Info("BlockActor received restart msg")
	case *actor.Restarting:
		log.Info("BlockActor received restarting msg")
	case *actor.Stop:
		log.Info("BlockActor received Stop msg")
	case *actor.Stopped:
		log.Info("BlockActor received Stopped msg")
	default:
		return false
	}

	return true
}

//Receive process chain msg
func (c *ChainActor) Receive(context actor.Context) {

	if handleSystemMsg(context) {
		return
	}

	switch msg := context.Message().(type) {
	case *message.InsertBlockReq:
		c.HandleNewProducedBlock(context, msg)
	case *message.ReceiveBlock:
		c.HandleReceiveBlock(context, msg)
	case *message.QueryTrxReq:
		c.HandleQueryTrxReq(context, msg)
	case *message.QueryBlockReq:
		c.HandleQueryBlockReq(context, msg)
	case *message.QueryChainInfoReq:
		c.HandleQueryChainInfoReq(context, msg)
	default:
		log.Error("BlockActor received Unknown msg")
	}
}

//HandleNewProducedBlock new block msg
func (c *ChainActor) HandleNewProducedBlock(ctx actor.Context, req *message.InsertBlockReq) {
	errcode := actorEnv.Chain.InsertBlock(req.Block)
	if ctx.Sender() != nil {
		resp := &message.InsertBlockRsp{
			Hash:  req.Block.Hash(),
			Error: fmt.Errorf("Insert block error: %v", errcode),
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
	if errcode == chain.InsertBlockSuccess {
		r := &message.RemovePendingTrxsReq{Trxs: req.Block.Transactions}
		trxPoolActorPid.Tell(r)

		BroadCastBlock(req.Block)
		log.Infof("Broadcast block: block num:%v, trxn:%v, delegate: %s, hash: %x\n", req.Block.GetNumber(), len(req.Block.Transactions), req.Block.Header.Delegate, req.Block.Hash())
	}
}

//HandleReceiveBlock receive block
func (c *ChainActor) HandleReceiveBlock(ctx actor.Context, req *message.ReceiveBlock) {
	errcode := actorEnv.Chain.InsertBlock(req.Block)

	if ctx.Sender() != nil {
		resp := &message.ReceiveBlockResp{
			BlockNum: req.Block.GetNumber(),
			ErrorNo:  errcode,
		}
		ctx.Sender().Request(resp, ctx.Self())
	}

	if errcode == chain.InsertBlockSuccess {
		req := &message.RemovePendingTrxsReq{Trxs: req.Block.Transactions}
		trxPoolActorPid.Tell(req)
	}
}

//HandleQueryTrxReq query trx
func (c *ChainActor) HandleQueryTrxReq(ctx actor.Context, req *message.QueryTrxReq) {
	tx := actorEnv.TxStore.GetTransaction(req.TrxHash)
	if ctx.Sender() != nil {
		resp := &message.QueryTrxResp{}
		if tx == nil {
			resp.Error = log.Errorf("Transaction not found")
		} else {
			resp.Trx = tx
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
}

//HandleQueryBlockReq query block
func (c *ChainActor) HandleQueryBlockReq(ctx actor.Context, req *message.QueryBlockReq) {
	block := actorEnv.Chain.GetBlockByHash(req.BlockHash)
	if block == nil {
		block = actorEnv.Chain.GetBlockByNumber(req.BlockNumber)
	}
	if ctx.Sender() != nil {
		resp := &message.QueryBlockResp{}
		if block == nil {
			resp.Error = log.Errorf("Block not found")
		} else {
			resp.Block = block
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
}

//HandleQueryChainInfoReq query chain info
func (c *ChainActor) HandleQueryChainInfoReq(ctx actor.Context, req *message.QueryChainInfoReq) {
	if ctx.Sender() != nil {
		resp := &message.QueryChainInfoResp{}
		resp.HeadBlockNum = actorEnv.Chain.HeadBlockNum()
		resp.HeadBlockHash = actorEnv.Chain.HeadBlockHash()
		resp.HeadBlockTime = actorEnv.Chain.HeadBlockTime()
		resp.HeadBlockDelegate = actorEnv.Chain.HeadBlockDelegate()
		resp.LastConsensusBlockNum = actorEnv.Chain.LastConsensusBlockNum()
		ctx.Sender().Request(resp, ctx.Self())
	}
}

func BroadCastBlock(block *types.Block) {
	broadCastBlock := &message.NotifyBlock{block}
	NetActorPid.Tell(broadCastBlock)
	return
}
