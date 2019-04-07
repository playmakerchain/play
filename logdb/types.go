// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package logdb

import (
	"math/big"

	"github.com/playmakerchain//block"
	"github.com/playmakerchain//"
	"github.com/playmakerchain//tx"
)

//Event represents tx.Event that can be stored in db.
type Event struct {
	BlockID     .Bytes32
	Index       uint32
	BlockNumber uint32
	BlockTime   uint64
	TxID        .Bytes32
	TxOrigin    .Address //contract caller
	Address     .Address // always a contract address
	Topics      [5]*.Bytes32
	Data        []byte
}

//newEvent converts tx.Event to Event.
func newEvent(header *block.Header, index uint32, txID .Bytes32, txOrigin .Address, txEvent *tx.Event) *Event {
	ev := &Event{
		BlockID:     header.ID(),
		Index:       index,
		BlockNumber: header.Number(),
		BlockTime:   header.Timestamp(),
		TxID:        txID,
		TxOrigin:    txOrigin,
		Address:     txEvent.Address, // always a contract address
		Data:        txEvent.Data,
	}
	for i := 0; i < len(txEvent.Topics) && i < len(ev.Topics); i++ {
		ev.Topics[i] = &txEvent.Topics[i]
	}
	return ev
}

//Transfer represents tx.Transfer that can be stored in db.
type Transfer struct {
	BlockID     .Bytes32
	Index       uint32
	BlockNumber uint32
	BlockTime   uint64
	TxID        .Bytes32
	TxOrigin    .Address
	Sender      .Address
	Recipient   .Address
	Amount      *big.Int
}

//newTransfer converts tx.Transfer to Transfer.
func newTransfer(header *block.Header, index uint32, txID .Bytes32, txOrigin .Address, transfer *tx.Transfer) *Transfer {
	return &Transfer{
		BlockID:     header.ID(),
		Index:       index,
		BlockNumber: header.Number(),
		BlockTime:   header.Timestamp(),
		TxID:        txID,
		TxOrigin:    txOrigin,
		Sender:      transfer.Sender,
		Recipient:   transfer.Recipient,
		Amount:      transfer.Amount,
	}
}

type RangeType string

const (
	Block RangeType = "block"
	Time  RangeType = "time"
)

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type Range struct {
	Unit RangeType
	From uint64
	To   uint64
}

type Options struct {
	Offset uint64
	Limit  uint64
}

type EventCriteria struct {
	Address *.Address // always a contract address
	Topics  [5]*.Bytes32
}

//EventFilter filter
type EventFilter struct {
	CriteriaSet []*EventCriteria
	Range       *Range
	Options     *Options
	Order       Order //default asc
}

type TransferCriteria struct {
	TxOrigin  *.Address //who send transaction
	Sender    *.Address //who transferred tokens
	Recipient *.Address //who recieved tokens
}

type TransferFilter struct {
	TxID        *.Bytes32
	CriteriaSet []*TransferCriteria
	Range       *Range
	Options     *Options
	Order       Order //default asc
}
