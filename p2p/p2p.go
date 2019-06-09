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
	"fmt"
	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/config"
	log "github.com/cihub/seelog"
	"net"
	"strconv"
)

//LocalPeerInfo ourself node info
var LocalPeerInfo PeerInfo

//Runner p2p global instance
var Runner *P2PServer

//P2PServer p2p server
type P2PServer struct {
	c      *collection
	connCb NewconnCb

	sendc  chan UniMsgPacket
	bsendc chan BcastMsgPacket
}

//SendupCb send up callback when receive a packet
type SendupCb func(index uint16, p *Packet)

//NewconnCb create a new candidate instance when accept a connection
type NewconnCb func(conn net.Conn)

//MakeP2PServer create instance
func MakeP2PServer(p *config.Parameter) *P2PServer {
	LocalPeerInfo.Addr = p.P2PServAddr
	LocalPeerInfo.Port = strconv.Itoa(p.P2PPort)
	LocalPeerInfo.ChainId = common.BytesToHex(config.GetChainID())

	id := LocalPeerInfo.Addr + LocalPeerInfo.Port
	LocalPeerInfo.Id = common.DoubleSha256([]byte(id)).ToHexString()

	Runner = &P2PServer{
		c:      createCollection(),
		sendc:  make(chan UniMsgPacket, 30),
		bsendc: make(chan BcastMsgPacket, 30),
	}

	return Runner
}

//Start start p2p
func (s *P2PServer) Start() {
	/*start listen*/
	go s.listenRoutine()
	go s.sendRoutine()
}

//SetCallback set new connection call back
func (s *P2PServer) SetCallback(conn NewconnCb) {
	s.connCb = conn
}

//SendUnicast send a  packet to a peer
func (s *P2PServer) SendUnicast(packet UniMsgPacket) {
	s.sendc <- packet
}

//SendBroadcast send a packet to some peer which is not set filter
func (s *P2PServer) SendBroadcast(packet BcastMsgPacket) {
	s.bsendc <- packet
}

//AddPeer add a peer
func (s *P2PServer) AddPeer(peer *Peer) error {
	return s.c.addPeer(peer)
}

//GetPeer get a peer by index
func (s *P2PServer) GetPeer(index uint16) *PeerInfo {
	return s.c.getPeer(index)
}

//DelPeer delete a peer by index
func (s *P2PServer) DelPeer(index uint16) bool {
	return s.c.delPeer(index)
}

//IsPeerExist judege if a peer exist or not by index
func (s *P2PServer) IsPeerExist(index uint16) bool {
	return s.c.isPeerExist(index)
}

//IsPeerInfoExist judge if a peer exist or not by peer info
func (s *P2PServer) IsPeerInfoExist(info PeerInfo) bool {
	return s.c.isPeerInfoExist(info)
}

//GetPeers get all peers
func (s *P2PServer) GetPeers() []PeerInfo {
	return s.c.getPeers()
}

//GetPeersData get a peer's info
func (s *P2PServer) GetPeersData() PeerDataSet {
	return s.c.getPeersData()
}

func (s *P2PServer) listenRoutine() {
	l, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(LocalPeerInfo.Port))
	if err != nil {
		log.Errorf("p2p start p2p server listen error: %s", err)
		panic(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("p2p NetServer::Listening() Failed to accept")
			continue
		}

		/*accpent ten new conection per second*/

		go s.connCb(conn)
	}

	return
}

func (s *P2PServer) sendRoutine() {
	for {
		select {
		case packet := <-s.bsendc:
			s.msend(&packet)
			continue
		default:
			select {
			case packet := <-s.bsendc:
				s.msend(&packet)
			case packet := <-s.sendc:
				s.send(&packet)
			}
		}
	}
}

func (s *P2PServer) send(packet *UniMsgPacket) {
	s.c.send(packet)
}

func (s *P2PServer) msend(packet *BcastMsgPacket) {
	s.c.sendBroadcast(packet)
}
