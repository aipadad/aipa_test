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

package p2p

import (
	"errors"
	log "github.com/cihub/seelog"
	"sync"
)

type collection struct {
	peers map[uint16]*Peer

	lock sync.RWMutex
}

func createCollection() *collection {
	c := &collection{
		peers: make(map[uint16]*Peer),
	}

	return c
}

func (c *collection) addPeer(peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("p2p collection add peer index: %d, id: %s , add: %s, port: %s",
		peer.Index, peer.Info.Id, peer.Info.Addr, peer.Info.Port)

	if peer.Info.IsIncomplete() {
		log.Info("p2p peer info error")
		return errors.New("peer info error")
	}

	for _, p := range c.peers {
		if p.Info.Equal(peer.Info) {
			if p.isconn {
				log.Info("p2p peer is already exist")
				return errors.New("peer is already exist")
			}
		}
	}

	c.peers[peer.Index] = peer
	return nil
}

func (c *collection) getPeer(index uint16) *PeerInfo {
	c.lock.Lock()
	defer c.lock.Unlock()

	var info PeerInfo
	peer, ok := c.peers[index]
	if ok {
		info.Addr = peer.Info.Addr
		info.Port = peer.Info.Port
		return &info
	}

	return nil
}

func (c *collection) delPeer(index uint16) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("p2p collection delete peer index: %d", index)
	peer, ok := c.peers[index]
	if ok {
		log.Debugf("p2p delete peer index: %d, %s:%s", peer.Index, peer.Info.Addr, peer.Info.Port)
		if peer.isconn {
			log.Error("p2p peer is connected , don't delete")
			return false
		}

		log.Error("p2p peer is disconnected , delete")
		peer.Stop()
		delete(c.peers, index)
		return true
	}

	log.Error("p2p peer not exist")
	return false
}

func (c *collection) isPeerExist(index uint16) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.peers[index]
	return ok
}

func (c *collection) isPeerInfoExist(info PeerInfo) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, value := range c.peers {
		if value.Info.Equal(info) {
			return true
		}
	}

	return false
}

func (c *collection) getPeers() []PeerInfo {
	c.lock.Lock()
	defer c.lock.Unlock()

	var peers []PeerInfo
	for _, p := range c.peers {
		peers = append(peers, p.Info)
	}

	return peers
}

func (c *collection) getPeersData() PeerDataSet {
	c.lock.Lock()
	defer c.lock.Unlock()

	var peers PeerDataSet
	for _, p := range c.peers {
		peers = append(peers, PeerData{Id: p.Info.Id, Index: p.Index})
	}

	return peers
}

func (c *collection) send(msg *UniMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	peer, ok := c.peers[msg.Index]
	if !ok {
		log.Errorf("p2p peer not exist %s", msg.Index)
		return
	}

	if !peer.isconn {
		return
	}

	go peer.Send(msg.P)

}

func (c *collection) sendBroadcast(msg *BcastMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for id, peer := range c.peers {
		if !peer.isconn {
			continue
		}

		if len(msg.Indexs) == 0 {
			go peer.Send(msg.P)
			continue
		}

		//filter index by msg index set
		bsend := true
		for _, eid := range msg.Indexs {
			if id == eid {
				bsend = false
				break
			}
		}

		if bsend {
			go peer.Send(msg.P)
		}

	}
}
