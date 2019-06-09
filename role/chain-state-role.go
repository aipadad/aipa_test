package role

import (
	"encoding/json"

	"github.com/aipadad/aipa/common"
	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/db"
)

// ChainState is definition of chain state
type ChainState struct {
	LastBlockNum          uint64      `json:"last_block_num"`
	LastBlockHash         common.Hash `json:"last_block_hash"`
	LastBlockTime         uint64      `json:"last_block_time"`
	LastConsensusBlockNum uint64      `json:"last_consensus_block_num"`
	CurrentDelegate       string      `json:"current_delegate"`
	CurrentAbsoluteSlot   uint64      `json:"current_absolute_slot"`
	RecentSlotFilled      uint64      `json:"recent_slot_filled"`
}

const (
	// ChainStateName is definition of chain state name
	ChainStateName string = "chain_state"
	// ChainStateDefaultKey is definition of chain stake default key
	ChainStateDefaultKey string = "chain_state_defkey"
)

func getGenesisTime() uint64 {
	t := config.Genesis.GenesisTime
	return uint64(t)
}

// CreateChainStateRole is to save init chain state
func CreateChainStateRole(ldb *db.DBService) error {
	_, err := ldb.GetObject(ChainStateName, ChainStateDefaultKey)
	if err != nil {
		object := &ChainState{
			LastBlockTime:    getGenesisTime(),
			RecentSlotFilled: ^uint64(0),
		}
		return SetChainStateRole(ldb, object)
	}
	return nil
}

// SetChainStateRole is to save chain state
func SetChainStateRole(ldb *db.DBService, value *ChainState) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(ChainStateName, ChainStateDefaultKey, string(jsonvalue))
}

//GetChainStateRole is to get chain state
func GetChainStateRole(ldb *db.DBService) (*ChainState, error) {
	value, err := ldb.GetObject(ChainStateName, ChainStateDefaultKey)
	if err != nil {
		return nil, err
	}

	res := &ChainState{}
	json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
