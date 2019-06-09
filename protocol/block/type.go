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

package block

import "github.com/aipadad/aipa/common/types"

//DO NOT EDIT
const (
	BLOCK_REQ = 1
	//BLOCK_INFO update or response
	BLOCK_UPDATE = 2

	LAST_BLOCK_NUMBER_REQ = 3
	LAST_BLOCK_NUMBER_RSP = 4

	BLOCK_HEADER_REQ = 5
	BLOCK_HEADER_RSP = 6

	BLOCK_HEADER_UPDATE = 7

	BLOCK_CATCH_REQUEST  = 8
	BLOCK_CATCH_RESPONSE = 9
)

type chainNumber struct {
	LibNumber   uint64
	BlockNumber uint64
}

type headerReq struct {
	index uint16
	req   *blockHeaderReq
}

type blockHeaderReq struct {
	Begin uint64
	End   uint64
}

type blockHeaderRsp struct {
	set []types.Header
}

type blockUpdate struct {
	index uint16
	block *types.Block
}

type headerUpdate struct {
	index  uint16
	header *types.Header
}
