// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package params

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/playmakerchain//lvldb"
	"github.com/playmakerchain//state"
	"github.com/playmakerchain//"
)

func TestParamsGetSet(t *testing.T) {
	kv, _ := lvldb.NewMem()
	st, _ := state.New(.Bytes32{}, kv)
	setv := big.NewInt(10)
	key := .BytesToBytes32([]byte("key"))
	p := New(.BytesToAddress([]byte("par")), st)
	p.Set(key, setv)

	getv := p.Get(key)
	assert.Equal(t, setv, getv)

	assert.Nil(t, st.Err())
}
