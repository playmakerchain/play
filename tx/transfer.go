// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package tx

import (
	"math/big"

	"github.com/playmakerchain/powerplay/powerplay"
)

// Transfer token transfer log.
type Transfer struct {
	Sender    powerplay.Address
	Recipient powerplay.Address
	Amount    *big.Int
}

// Transfers slisce of transfer logs.
type Transfers []*Transfer
