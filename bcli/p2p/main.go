package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aipadad/aipa/bcli/p2p/stub"
	"github.com/aipadad/aipa/common/types"
	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/context"
	"github.com/aipadad/aipa/protocol"
	log "github.com/cihub/seelog"
)

type p2pConfig struct {
	ServAddr string
	ServPort string
	PeerLst  []string
	ChainId  string
	Producer bool
}

type chainConfig struct {
	LibNumber   uint32
	BlockNumber uint32
	Blocks      []types.Block
}

func main() {

	pc := readP2PConfig("p2pconfig.json")
	if pc == nil {
		return
	}
	port, _ := strconv.Atoi(pc.ServPort)
	param := config.P2PConfig{P2PServAddr: pc.ServAddr,
		P2PPort:  port,
		PeerList: pc.PeerLst,
	}

	bc := stub.MakeBlockChainStub()

	chain := readChainConfig("chainconfig.json")
	if chain != nil {
		if chain.BlockNumber > uint32(len(chain.Blocks)) {
			fmt.Printf("chain cfg number error")
			return
		}
		bc.SetBlocks(chain.Blocks[0:chain.BlockNumber])
		bc.SetLibNumber(chain.LibNumber)
	}

	log.Info("blocknumber:", chain.BlockNumber)
	p := protocol.MakeProtocol(&param, bc)

	actor := stub.NewDumActor(bc)

	p.SetChainActor(actor)

	go p.Start()

	if pc.Producer {
		go newBlockTimer(bc, p)
	}
	/*new block timer*/

	select {}
}

func newBlockTimer(bc *stub.BlockChainStub, p context.ProtocolInstance) {
	time.Sleep(1 * time.Minute)

	blockTimer := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-blockTimer.C:
			newBlock(bc, p)
			blockTimer.Reset(2 * time.Second)
		}
	}
}

func newBlock(bc *stub.BlockChainStub, p context.ProtocolInstance) {
	if p.GetBlockSyncState() {
		msg := bc.NewBlockMsg()
		p.SendNewBlock(msg)
	}
}

//ReadFile parse json configuration
func readP2PConfig(filename string) *p2pConfig {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("read p2p config error:%s", err)
		return nil
	}

	str := string(bytes)

	var pc p2pConfig
	if err := json.Unmarshal([]byte(str), &pc); err != nil {
		fmt.Printf("p2p config unmarshall error:%s", err)
		return nil
	}

	return &pc
}

//ReadFile parse json configuration
func readChainConfig(filename string) *chainConfig {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("read chain config error:%s", err)
		return nil
	}

	str := string(bytes)

	var pc chainConfig
	if err := json.Unmarshal([]byte(str), &pc); err != nil {
		fmt.Printf("chain config unmarshall error:%s", err)
		return nil
	}

	return &pc
}
