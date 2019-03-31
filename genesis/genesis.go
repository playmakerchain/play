// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package genesis

import (
	"encoding/hex"

	"github.com//thor/abi"
	"github.com//thor/block"
	"github.com//thor/state"
	"github.com//thor/thor"
	"github.com//thor/tx"
)

// Genesis to build genesis block.
type Genesis struct {
	builder *Builder
	id      thor.Bytes32
	name    string
}

// Build build the genesis block.
func (g *Genesis) Build(stateCreator *state.Creator) (blk *block.Block, events tx.Events, err error) {
	block, events, err := g.builder.Build(stateCreator)
	if err != nil {
		return nil, nil, err
	}
	if block.Header().ID() != g.id {
		panic("built genesis ID incorrect")
	}
	return block, events, nil
}

// ID returns genesis block ID.
func (g *Genesis) ID() thor.Bytes32 {
	return g.id
}

// Name returns network name.
func (g *Genesis) Name() string {
	return g.name
}

func mustEncodeInput(abi *abi.ABI, name string, args ...interface{}) []byte {
	m, found := abi.MethodByName(name)
	if !found {
		panic("method not found")
	}
	data, err := m.EncodeInput(args...)
	if err != nil {
		panic(err)
	}
	return data
}

func mustDecodeHex(str string) []byte {
	data, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return data
}

var emptyRuntimeBytecode = mustDecodeHex("6060604052600256")
