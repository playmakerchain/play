// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package authority

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/playmakerchain//lvldb"
	"github.com/playmakerchain//state"
	"github.com/playmakerchain//"
)

func M(a ...interface{}) []interface{} {
	return a
}

func TestAuthority(t *testing.T) {
	kv, _ := lvldb.NewMem()
	st, _ := state.New(.Bytes32{}, kv)

	p1 := .BytesToAddress([]byte("p1"))
	p2 := .BytesToAddress([]byte("p2"))
	p3 := .BytesToAddress([]byte("p3"))

	st.SetBalance(p1, big.NewInt(10))
	st.SetBalance(p2, big.NewInt(20))
	st.SetBalance(p3, big.NewInt(30))

	aut := New(.BytesToAddress([]byte("aut")), st)
	tests := []struct {
		ret      interface{}
		expected interface{}
	}{
		{aut.Add(p1, p1, .Bytes32{}), true},
		{M(aut.Get(p1)), []interface{}{true, p1, .Bytes32{}, true}},
		{aut.Add(p2, p2, .Bytes32{}), true},
		{aut.Add(p3, p3, .Bytes32{}), true},
		{M(aut.Candidates(big.NewInt(10), .MaxBlockProposers)), []interface{}{
			[]*Candidate{{p1, p1, .Bytes32{}, true}, {p2, p2, .Bytes32{}, true}, {p3, p3, .Bytes32{}, true}},
		}},
		{M(aut.Candidates(big.NewInt(20), .MaxBlockProposers)), []interface{}{
			[]*Candidate{{p2, p2, .Bytes32{}, true}, {p3, p3, .Bytes32{}, true}},
		}},
		{M(aut.Candidates(big.NewInt(30), .MaxBlockProposers)), []interface{}{
			[]*Candidate{{p3, p3, .Bytes32{}, true}},
		}},
		{M(aut.Candidates(big.NewInt(10), 2)), []interface{}{
			[]*Candidate{{p1, p1, .Bytes32{}, true}, {p2, p2, .Bytes32{}, true}},
		}},
		{M(aut.Get(p1)), []interface{}{true, p1, .Bytes32{}, true}},
		{aut.Update(p1, false), true},
		{M(aut.Get(p1)), []interface{}{true, p1, .Bytes32{}, false}},
		{aut.Update(p1, true), true},
		{M(aut.Get(p1)), []interface{}{true, p1, .Bytes32{}, true}},
		{aut.Revoke(p1), true},
		{M(aut.Get(p1)), []interface{}{false, p1, .Bytes32{}, false}},
		{M(aut.Candidates(&big.Int{}, .MaxBlockProposers)), []interface{}{
			[]*Candidate{{p2, p2, .Bytes32{}, true}, {p3, p3, .Bytes32{}, true}},
		}},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.ret)
	}

	assert.Nil(t, st.Err())

}
