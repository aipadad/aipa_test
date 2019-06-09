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
 * file description:  support uint128
 * @Author: 
 * @Date:   2018-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"math/big"
)

//MaxUint128 return max value of 128 bit
func MaxUint128() *big.Int {
	maxValue := bigpow(2, 128)
	return maxValue
}

//MaxUint256 return max value of 128 bit
func MaxUint256() *big.Int {
	maxValue := bigpow(2, 256)
	return maxValue
}

//bigpow returns a ** b as a big integer.
func bigpow(a int64, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}
