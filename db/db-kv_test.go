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
 * file description: database interface
 * @Author: 
 * @Date:   2018-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package db

import (
	//"fmt"
	"testing"

	log "github.com/cihub/seelog"
)

func TestDBService_Callput(t *testing.T) {
	log.Info("abc")
	ins := NewDbService("./db")
	ins.Put([]byte("abc"), []byte("123"))
	res, _ := ins.Get([]byte("abc"))
	log.Info(res)
	ins.Close()
}
func TestDBService_CallGet(t *testing.T) {
	log.Info("abc")
	ins := NewDbService("./db")
	ins.Put([]byte("abc"), []byte("123"))
	res, _ := ins.Get([]byte("abc"))
	log.Info(res)
}
func TestDBService_CallFlush(t *testing.T) {
	log.Info("abc")
	ins := NewDbService("./db")
	ins.Put([]byte("abc"), []byte("123"))
	ins.Flush()
}
