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
 * file description: option database interface
 * @Author: 
 * @Date:   2018-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package optiondb

import (
	"errors"
	//"fmt"

	"github.com/aipadad/aipa/config"
	log "github.com/cihub/seelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//OptionDbRepo is the interface for plugin db for user queries
type OptionDbRepo interface {
	InsertOptionDb(collection string, value interface{}) error
	OptionDbFind(collection string, key string, value interface{}) (interface{}, error)
	OptionDbUpdate(collection string, key string, value interface{}, updatekey string, updatevalue interface{}) error
	CallClose()
}

var insertSession *MongoContext
var getSession *MongoContext

//InsertOptionDb is to insert record in option db
func (r *OptionDbRepository) InsertOptionDb(collection string, value interface{}) error {
	var err error
	
	if r.is_optiondb_offline == true {
		return nil
	}

	if insertSession == nil || insertSession.mgoSession == nil {
		insertSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			r.is_optiondb_offline = true
			log.Error("collection cccccccccccc", insertSession, collection, err)
			return errors.New("Get session faild:" + r.mgoEndpoint)
		}
	}

	insert := func(c *mgo.Collection) error {
		return c.Insert(value)
	}
	insertSession.SetCollection(collection, insert)
	//insertSession.Close()

	return nil
}

//RemoveAllOptionDb is to remove all records in option db
func (r *OptionDbRepository) RemoveAllOptionDb(collection string) error {
	var err error
	
	if r.is_optiondb_offline == true {
		return nil
	}

	if getSession == nil || getSession.mgoSession == nil {
		getSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			r.is_optiondb_offline = true
			log.Error("collection ", getSession, collection, err)
			return errors.New("Get session faild" + r.mgoEndpoint)
		}
	}
	removeAll := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(nil)
		return err
	}
	getSession.SetCollection(collection, removeAll)

	return nil
}

//OptionDbFind is to find record in option db
func (r *OptionDbRepository) OptionDbFind(collection string, key string, value interface{}) (interface{}, error) {
	var err error

	if r.is_optiondb_offline == true {
		return nil, errors.New("Optiondb is out of connection")
	}
	
	if getSession == nil || getSession.mgoSession == nil {
		getSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			r.is_optiondb_offline = true
			log.Error("collection ", getSession, collection, err)
			return nil, errors.New("Get session faild" + r.mgoEndpoint)
		}
	}
	var mesgs interface{}
	getSession.GetCollection(config.DEFAULT_OPTIONDB_NAME, collection).Find(bson.M{"$or": []bson.M{{key: value}}}).One(&mesgs)
	if mesgs == nil {
		return nil, errors.New("No record is found.")
	}

	return mesgs, nil
}

//OptionDbUpdate is to update record in option db
func (r *OptionDbRepository) OptionDbUpdate(collection string, key string, value interface{}, updatekey string, updatevalue interface{}) error {
	var err error
	
	if r.is_optiondb_offline == true {
		return errors.New("Optiondb is out of connection")
	}
	
	selector := bson.M{key: value}
	data := bson.M{"$set": bson.M{updatekey: updatevalue}}

	if getSession == nil || getSession.mgoSession == nil {
		getSession, err = GetSession(r.mgoEndpoint)
		if err != nil {
			r.is_optiondb_offline = true
			log.Error("collection ", getSession, collection, err)
			return errors.New("Get session faild" + r.mgoEndpoint)
		}
	}

	_, err = getSession.mgoSession.DB(config.DEFAULT_OPTIONDB_NAME).C(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME).UpdateAll(selector, data)

	return err
}
func (r *OptionDbRepository) CallClose() {
	insertSession.Close()
	getSession.Close()
}
