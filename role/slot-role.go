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
 * file description: producer
 * @Author: 
 * @Date:   2018-12-11
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"github.com/aipadad/aipa/config"
)

//GetSlotAtTime is geting slot at time
func (r *Role) GetSlotAtTime(current uint64) uint64 {
	firstSlotTime := r.GetSlotTime(1)

	if current < firstSlotTime {
		return 0
	}
	return (current-firstSlotTime)/uint64(config.DEFAULT_BLOCK_INTERVAL) + 1
}

//GetSlotTime is geting slot time
func (r *Role) GetSlotTime(slotNum uint64) uint64 {

	if slotNum == 0 {
		return 0
	}
	interval := config.DEFAULT_BLOCK_INTERVAL

	object, err := r.GetChainState()
	if err != nil {
		return 0
	}

	if object.LastBlockNum == 0 {
		genesisTime := config.Genesis.GenesisTime
		return genesisTime + slotNum*uint64(interval)
	}

	headBlockAbsSlot := object.LastBlockTime / uint64(interval)
	headSlotTime := headBlockAbsSlot * uint64(interval)
	return headSlotTime + slotNum*uint64(interval)
}
