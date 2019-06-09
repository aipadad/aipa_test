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
 * file description: code database interface
 * @Author: 
 * @Date:   2018-12-05
 * @Last Modified by:
 * @Last Modified time:
 */

package db

//StartUndoSession is to start undo session
func (d *DBService) StartUndoSession() {
	d.codeRepo.CallStartUndoSession(true)
}

//CreatObjectIndex is to create object index when one object's attributed is going to
//search or sort.
func (d *DBService) CreatObjectIndex(objectName string, indexName string, indexJson string) error {
	return d.codeRepo.CallCreatObjectIndex(objectName, indexName, indexJson)
}

//CreatObjectUintIndex is to create object index when one object's attributed is going to
//search or sort.
func (d *DBService) CreatObjectMultiIndex(objectName string, indexName string, indexJson string, secKey string) error {
	return d.codeRepo.CallCreatObjectMultiIndex(objectName, indexName, indexJson, secKey)
}

//SetObject is to set object by key and value. which should have key has create index.
func (d *DBService) SetObject(objectName string, key string, objectValue string) error {
	return d.codeRepo.CallSetObject(objectName, key, objectValue)
}

//GetObject is to get object by key and return value. which should have key has create index.
func (d *DBService) GetObject(objectName string, key string) (string, error) {
	return d.codeRepo.CallGetObject(objectName, key)
}

//GetObjectByIndex is to get object by index name and index value,
//which index name has create object index in the db
func (d *DBService) GetObjectByIndex(objectName string, indexName string, indexValue string) (string, error) {
	return d.codeRepo.CallGetObjectByIndex(objectName, indexName, indexValue)
}

//GetAllObjectKeys is to get all objects by object name
func (d *DBService) GetAllObjectKeys(objectName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjectKeys(objectName)
}

//GetAllObjects is to get all objects by keyName
func (d *DBService) GetAllObjects(keyName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjects(keyName)
}

//GetAllObjectsSortByIndex is to get all objects by sorted index
func (d *DBService) GetAllObjectsSortByIndex(indexName string) ([]string, error) {
	return d.codeRepo.CallGetAllObjectsSortByIndex(indexName)
}

//DeleteObject is to delete object by object and key
func (d *DBService) DeleteObject(objectName string, key string) (string, error) {
	return d.codeRepo.CallDeleteObject(objectName, key)
}

//Commit is to commit object
func (d *DBService) Commit() error {
	return d.codeRepo.CallCommit()
}

//Rollback is to rollback object
func (d *DBService) Rollback() error {
	return d.codeRepo.CallRollback()
}
