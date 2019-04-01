// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package chain

import (
	"encoding/binary"
	"testing"

	"github.com/vechain/thor/lvldb"
	"github.com/vechain/thor/thor"
)

func BenchmarkGet(b *testing.B) {
	kv, _ := lvldb.NewMem()
	at := newAncestorTrie(kv)

	const maxBN = 1000
	for bn := uint32(0); bn < maxBN; bn++ {
		var id, parentID thor.Bytes32
		binary.BigEndian.PutUint32(id[:], bn)
		binary.BigEndian.PutUint32(parentID[:], bn-1)
		if err := at.Update(kv, id, parentID); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bn := uint32(i) % maxBN
		if bn == 0 {
			bn = maxBN / 2
		}
		var id thor.Bytes32
		binary.BigEndian.PutUint32(id[:], bn)
		if _, err := at.GetAncestor(id, bn-1); err != nil {
			b.Fatal(err)
		}
	}
}
