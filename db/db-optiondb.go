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

import (
	"errors"
)

//Insert is to insert record to option db
func (d *OptionDBService) Insert(collection string, value interface{}) error {
	if d.optDbRepo == nil {
		return nil
	}
	return d.optDbRepo.InsertOptionDb(collection, value)
}

//Find is to find record in option db
func (d *OptionDBService) Find(collection string, key string, value interface{}) (interface{}, error) {
	if d.optDbRepo == nil {
		return nil, errors.New("error optiondb is not init")
	}
	return d.optDbRepo.OptionDbFind(collection, key, value)
}

//Update is to update record in option db
func (d *OptionDBService) Update(collection string, key string, value interface{}, updatekey string, updatevalue interface{}) error {
	if d.optDbRepo == nil {
		return errors.New("error optiondb is not init")
	}

	return d.optDbRepo.OptionDbUpdate(collection, key, value, updatekey, updatevalue)
}

func (d *OptionDBService) IsOpDbConfigured() bool {
	if d.optDbRepo == nil {
		return false
	}
	return true
}
