// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 THe PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package chain

import (
	"github.com/playmakerchain/thor/block"
)

// Fork describes forked chain.
type Fork struct {
	Ancestor *block.Header
	Trunk    []*block.Header
	Branch   []*block.Header
}
