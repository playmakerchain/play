// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package chain

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/playmakerchain/powerplay/block"
	"github.com/playmakerchain/powerplay/powerplay"
)

func TestRawBlock(t *testing.T) {
	b := new(block.Builder).ParentID(powerplay.Bytes32{1, 2, 3}).Build()

	priv, _ := crypto.GenerateKey()
	sig, err := crypto.Sign(b.Header().SigningHash().Bytes(), priv)
	assert.Nil(t, err)
	b = b.WithSignature(sig)

	data, _ := rlp.EncodeToBytes(b)
	raw := &rawBlock{raw: data}

	h, _ := raw.Header()
	assert.Equal(t, b.Header().ID(), h.ID())

	b1, _ := raw.Block()

	data, _ = rlp.EncodeToBytes(b1)
	assert.Equal(t, []byte(raw.raw), data)
}
