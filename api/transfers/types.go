// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package transfers

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/playmakerchain/play/api/transactions"
	"github.com/playmakerchain/play/logdb"
	"github.com/playmakerchain/play/play"
)

type FilteredTransfer struct {
	Sender    play.Address          `json:"sender"`
	Recipient play.Address          `json:"recipient"`
	Amount    *math.HexOrDecimal256 `json:"amount"`
	Meta      transactions.LogMeta  `json:"meta"`
}

func convertTransfer(transfer *logdb.Transfer) *FilteredTransfer {
	v := math.HexOrDecimal256(*transfer.Amount)
	return &FilteredTransfer{
		Sender:    transfer.Sender,
		Recipient: transfer.Recipient,
		Amount:    &v,
		Meta: transactions.LogMeta{
			BlockID:        transfer.BlockID,
			BlockNumber:    transfer.BlockNumber,
			BlockTimestamp: transfer.BlockTime,
			TxID:           transfer.TxID,
			TxOrigin:       transfer.TxOrigin,
		},
	}
}
