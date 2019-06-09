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
 * file description:  core state role
 * @Author: 
 * @Date:   2018-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/aipadad/aipa/config"
	"github.com/aipadad/aipa/db"
)

// CoreState is definition of core state
type CoreState struct {
	Config           ChainConfig `json:"chain_config"`
	CurrentDelegates []string    `json:"current_delegates"`
}

// ChainConfig is definition of chain config
type ChainConfig struct {
	MaxBlockSize   uint32 `json:"max_block_size"`
	MaxTrxLifetime uint32 `json:"max_trx_lifetime"`
	MaxTrxRuntime  uint32 `json:"max_trx_runtime"`
	InDepthLeimit  uint32 `json:"in_depth_limit"`
}

const (
	// CoreStateName is definition of core state name
	CoreStateName string = "core_state"
	// CoreStateDefaultKey is definition of core state default key
	CoreStateDefaultKey string = "core_state_defkey"
)

// CreateCoreStateRole is to save init core state
func CreateCoreStateRole(ldb *db.DBService) error {
	_, err := ldb.GetObject(CoreStateName, CoreStateDefaultKey)
	if err != nil {
		dgp := &CoreState{
			Config: ChainConfig{
				MaxBlockSize:   5242880,
				MaxTrxLifetime: 3600,
				MaxTrxRuntime:  10000,
				InDepthLeimit:  4,
			},
			CurrentDelegates: []string{},
		}
		return SetCoreStateRole(ldb, dgp)
	}
	return nil
}

// SetCoreStateRole is to save core state
func SetCoreStateRole(ldb *db.DBService, value *CoreState) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(CoreStateName, CoreStateDefaultKey, string(jsonvalue))
}

// GetCoreStateRole is to get core state
func GetCoreStateRole(ldb *db.DBService) (*CoreState, error) {
	value, err := ldb.GetObject(CoreStateName, CoreStateDefaultKey)
	if err != nil {
		return nil, err
	}
	res := new(CoreState)
	res.CurrentDelegates = make([]string, 0, config.MAX_DELEGATE_VOTES)
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
