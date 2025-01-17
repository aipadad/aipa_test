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
 * file description: database for key-value
 * @Author: 
 * @Date:   2018-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package kvdb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

//KvDBRepo is interface for leveldb
type KvDBRepo interface {
	CallPut(key []byte, value []byte) error
	CallGet(key []byte) ([]byte, error)
	CallDelete(key []byte) error
	CallClose()
	CallFlush() error
	CallSeek(prefixKey []byte) ([]string, error)

	CallNewIterator() iterator.Iterator
}
