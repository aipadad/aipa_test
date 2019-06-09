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
 * file description:  producer actor
 * @Author: 
 * @Date:   2018-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package consensus

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/aipadad/aipa/action/message"
	"github.com/aipadad/aipa/bpl"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/p2p"
	pcommon "github.com/aipadad/aipa/protocol/common"
	log "github.com/cihub/seelog"
)

type Consensus struct {
	actor *actor.PID
}

func MakeConsensus() *Consensus {
	return &Consensus{}
}

func (c *Consensus) SetActor(tid *actor.PID) {
	c.actor = tid
}

func (c *Consensus) Start() {

}

func (c *Consensus) Dispatch(index uint16, p *p2p.Packet) {
	switch p.H.PacketType {
	case ConsensusPreVote:
		c.processPrevote(index, p.Data)
	case ConsensusPreCommit:
		c.processPrecommit(index, p.Data)
	case ConsensusCommit:
		c.processCommit(index, p.Data)
	}
}

func (c *Consensus) SendPrevote(notify *message.SendPrevote) {
	buf, err := bpl.Marshal(notify.BlockState)
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusPreVote,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)

}

func (c *Consensus) SendPrecommit(notify *message.SendPrecommit) {
	buf, err := bpl.Marshal(notify.BlockState)
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusPreCommit,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)

}

func (c *Consensus) SendCommit(notify *message.SendCommit) {
	buf, err := bpl.Marshal(notify.BftHeaderState)
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusCommit,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)
}

func (c *Consensus) processPrevote(index uint16, data []byte) {
	var block types.ConsensusBlockState
	err := bpl.Unmarshal(data, &block)
	if err != nil {
		log.Errorf("PROTOCOL consensus block Unmarshal error:%s, blockId%x", err, block.HeaderState.BlockId)
		return
	}

	prevote := &message.RcvPrevoteReq{BlockState: &block}
	c.actor.Tell(prevote)
}

func (c *Consensus) processPrecommit(index uint16, data []byte) {
	var block types.ConsensusBlockState
	err := bpl.Unmarshal(data, &block)
	if err != nil {
		log.Errorf("PROTOCOL consensus head Unmarshal error:%s,blockId%x", err, block.HeaderState.BlockId)
		return
	}

	precommit := &message.RcvPrecommitReq{BlockState: &block}
	c.actor.Tell(precommit)
}

func (c *Consensus) processCommit(index uint16, data []byte) {
	var head types.ConsensusHeaderState
	err := bpl.Unmarshal(data, &head)
	if err != nil {
		log.Errorf("PROTOCOL consensus head Unmarshal error:%s,blockId%x", err, head.BlockId)
		return
	}

	commit := &message.RcvCommitReq{BftHeaderState: &head}
	c.actor.Tell(commit)
}
