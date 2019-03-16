// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package block_test

import (
	"math"
	"testing"

	"github.com/playmakerchain/play/block"
	"github.com/playmakerchain/play/play"
	"github.com/stretchr/testify/assert"
)

func TestGasLimit_IsValid(t *testing.T) {

	tests := []struct {
		gl       uint64
		parentGL uint64
		want     bool
	}{
		{powerplay.MinGasLimit, powerplay.MinGasLimit, true},
		{powerplay.MinGasLimit - 1, powerplay.MinGasLimit, false},
		{powerplay.MinGasLimit, powerplay.MinGasLimit * 2, false},
		{powerplay.MinGasLimit * 2, powerplay.MinGasLimit, false},
		{powerplay.MinGasLimit + powerplay.MinGasLimit/powerplay.GasLimitBoundDivisor, powerplay.MinGasLimit, true},
		{powerplay.MinGasLimit*2 + powerplay.MinGasLimit/powerplay.GasLimitBoundDivisor, powerplay.MinGasLimit * 2, true},
		{powerplay.MinGasLimit*2 - powerplay.MinGasLimit/powerplay.GasLimitBoundDivisor, powerplay.MinGasLimit * 2, true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).IsValid(tt.parentGL))
	}
}

func TestGasLimit_Adjust(t *testing.T) {

	tests := []struct {
		gl    uint64
		delta int64
		want  uint64
	}{
		{play.MinGasLimit, 1, play.MinGasLimit + 1},
		{play.MinGasLimit, -1, play.MinGasLimit},
		{math.MaxUint64, 1, math.MaxUint64},
		{play.MinGasLimit, int64(play.MinGasLimit), play.MinGasLimit + play.MinGasLimit/play.GasLimitBoundDivisor},
		{play.MinGasLimit * 2, -int64(play.MinGasLimit), play.MinGasLimit*2 - (play.MinGasLimit*2)/play.GasLimitBoundDivisor},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).Adjust(tt.delta))
	}
}

func TestGasLimit_Qualify(t *testing.T) {
	tests := []struct {
		gl       uint64
		parentGL uint64
		want     uint64
	}{
		{play.MinGasLimit, play.MinGasLimit, play.MinGasLimit},
		{play.MinGasLimit - 1, play.MinGasLimit, play.MinGasLimit},
		{play.MinGasLimit, play.MinGasLimit * 2, play.MinGasLimit*2 - (play.MinGasLimit*2)/play.GasLimitBoundDivisor},
		{play.MinGasLimit * 2, play.MinGasLimit, play.MinGasLimit + play.MinGasLimit/play.GasLimitBoundDivisor},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).Qualify(tt.parentGL))
	}
}
