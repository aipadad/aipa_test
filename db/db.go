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

	"github.com/aipadad/aipa/db/platform/codedb"
	"github.com/aipadad/aipa/db/platform/kvdb"
	"github.com/aipadad/aipa/db/platform/optiondb"
	log "github.com/cihub/seelog"
)

//DBService is db service struct for other package
type DBService struct {
	kvRepo    kvdb.KvDBRepo
	codeRepo  codedb.CodeDbRepo
}

type OptionDBService struct {
	optDbRepo optiondb.OptionDbRepo
}

//NewDbService is to create a new db service with kv databse, codedb, and option db for optionally
func NewDbService(path string, codedbPath string) *DBService {
	kv, err := kvdb.NewKVDatabase(path)
	if err != nil {
		log.Errorf("DB load key value database failed ", path)
		return nil
	}
	db, err := codedb.NewMultindexDB(codedbPath)
	if err != nil {
		log.Errorf("DB load code database failed ", codedbPath)
		return nil
	}
	return &DBService{kvRepo: kv, codeRepo: db}

}

func NewOptionDbService(optPath string) *OptionDBService {
	if optPath == "" {
		return nil
	}
	optiondb := optiondb.NewOptionDbRepository(optPath)
	return &OptionDBService{optDbRepo: optiondb}

}

//Close is to close db.
func (d *DBService) Close() {
	log.Error("DB close all")
	d.codeRepo.CallUndoFlush()
	d.kvRepo.CallClose()
	d.codeRepo.CallClose()
}

func (d *OptionDBService) Close() {
	if d.optDbRepo != nil {
		log.Info("mongodb is not connect do not need close")
		d.optDbRepo.CallClose()
	}
}

//DBApi is listing all the interface that DBService provides.
type DBApi interface {
	Lock()
	UnLock()
	//kv database interface
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Flush() error
	Seek(prefixKey []byte) ([]string, error)
	NewBatch()
	BatchPut(key []byte, value []byte)
	BatchDelete(key []byte)
	BatchCommit() error
	//code db interface can rollback
	StartUndoSession()
	CreatObjectIndex(objectName string, indexName string, indexJson string) error
	CreatObjectMultiIndex(objectName string, indexName string, indexJson string, secKey string) error
	SetObject(objectName string, objectValue interface{}) error
	SetObjectByIndex(objectName string, indexName string, indexValue interface{}, objectValue interface{}) error
	SetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}, objectValue interface{}) error
	GetObject(objectName string, key string) (interface{}, error)
	GetObjectByIndex(objectName string, indexName string, indexValue string) (interface{}, error)
	GetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}) (interface{}, error)
	GetAllObjectKeys(objectName string) ([]string, error)
	GetAllObjects(keyName string) ([]string, error)
	GetAllObjectsSortByIndex(indexName string) ([]string, error)
	GetObjectsWithinRangeByIndex(indexName string, lessOrEqual string, greater string) ([]string, error)
	DeleteObject(objectName string, key string) (string, error)
	Commit()
	Rollback()
	Reset() //TODO
}

type OptionDBApi interface {
	IsOpDbConfigured() bool
	Insert(collection string, value interface{}) error
	Find(collection string, key string, value interface{}) (interface{}, error)
	Update(collection string, key string, value interface{}, updatekey string, updatevalue interface{}) error
}
