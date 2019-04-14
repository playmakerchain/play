// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package subscriptions

import (
	"github.com/playmakerchain/powerplay/chain"
	"github.com/playmakerchain/powerplay/powerplay"
)

type blockReader struct {
	chain       *chain.Chain
	blockReader chain.BlockReader
}

func newBlockReader(chain *chain.Chain, position powerplay.Bytes32) *blockReader {
	return &blockReader{
		chain:       chain,
		blockReader: chain.NewBlockReader(position),
	}
}

func (br *blockReader) Read() ([]interface{}, bool, error) {
	blocks, err := br.blockReader.Read()
	if err != nil {
		return nil, false, err
	}
	var msgs []interface{}
	for _, block := range blocks {
		msg, err := convertBlock(block)
		if err != nil {
			return nil, false, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, len(blocks) > 0, nil
}
