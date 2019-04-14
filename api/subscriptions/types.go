// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package subscriptions

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/playmakerchain/play/block"
	"github.com/playmakerchain/play/chain"
	"github.com/playmakerchain/play/play"
	"github.com/playmakerchain/play/tx"
)

//BlockMessage block piped by websocket
type BlockMessage struct {
	Number       uint32         `json:"number"`
	ID           play.Bytes32   `json:"id"`
	Size         uint32         `json:"size"`
	ParentID     play.Bytes32   `json:"parentID"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	Beneficiary  play.Address   `json:"beneficiary"`
	GasUsed      uint64         `json:"gasUsed"`
	TotalScore   uint64         `json:"totalScore"`
	TxsRoot      play.Bytes32   `json:"txsRoot"`
	StateRoot    play.Bytes32   `json:"stateRoot"`
	ReceiptsRoot play.Bytes32   `json:"receiptsRoot"`
	Signer       play.Address   `json:"signer"`
	Transactions []play.Bytes32 `json:"transactions"`
	Obsolete     bool           `json:"obsolete"`
}

func convertBlock(b *chain.Block) (*BlockMessage, error) {
	header := b.Header()
	signer, err := header.Signer()
	if err != nil {
		return nil, err
	}

	txs := b.Transactions()
	txIds := make([]play.Bytes32, len(txs))
	for i, tx := range txs {
		txIds[i] = tx.ID()
	}
	return &BlockMessage{
		Number:       header.Number(),
		ID:           header.ID(),
		ParentID:     header.ParentID(),
		Timestamp:    header.Timestamp(),
		TotalScore:   header.TotalScore(),
		GasLimit:     header.GasLimit(),
		GasUsed:      header.GasUsed(),
		Beneficiary:  header.Beneficiary(),
		Signer:       signer,
		Size:         uint32(b.Size()),
		StateRoot:    header.StateRoot(),
		ReceiptsRoot: header.ReceiptsRoot(),
		TxsRoot:      header.TxsRoot(),
		Transactions: txIds,
		Obsolete:     b.Obsolete,
	}, nil
}

type LogMeta struct {
	BlockID        play.Bytes32 `json:"blockID"`
	BlockNumber    uint32       `json:"blockNumber"`
	BlockTimestamp uint64       `json:"blockTimestamp"`
	TxID           play.Bytes32 `json:"txID"`
	TxOrigin       play.Address `json:"txOrigin"`
}

//TransferMessage transfer piped by websocket
type TransferMessage struct {
	Sender    play.Address          `json:"sender"`
	Recipient play.Address          `json:"recipient"`
	Amount    *math.HexOrDecimal256 `json:"amount"`
	Meta      LogMeta               `json:"meta"`
	Obsolete  bool                  `json:"obsolete"`
}

func convertTransfer(header *block.Header, tx *tx.Transaction, transfer *tx.Transfer, obsolete bool) (*TransferMessage, error) {
	signer, err := tx.Signer()
	if err != nil {
		return nil, err
	}

	return &TransferMessage{
		Sender:    transfer.Sender,
		Recipient: transfer.Recipient,
		Amount:    (*math.HexOrDecimal256)(transfer.Amount),
		Meta: LogMeta{
			BlockID:        header.ID(),
			BlockNumber:    header.Number(),
			BlockTimestamp: header.Timestamp(),
			TxID:           tx.ID(),
			TxOrigin:       signer,
		},
		Obsolete: obsolete,
	}, nil
}

//EventMessage event piped by websocket
type EventMessage struct {
	Address  play.Address   `json:"address"`
	Topics   []play.Bytes32 `json:"topics"`
	Data     string         `json:"data"`
	Meta     LogMeta        `json:"meta"`
	Obsolete bool           `json:"obsolete"`
}

func convertEvent(header *block.Header, tx *tx.Transaction, event *tx.Event, obsolete bool) (*EventMessage, error) {
	signer, err := tx.Signer()
	if err != nil {
		return nil, err
	}
	return &EventMessage{
		Address: event.Address,
		Data:    hexutil.Encode(event.Data),
		Meta: LogMeta{
			BlockID:        header.ID(),
			BlockNumber:    header.Number(),
			BlockTimestamp: header.Timestamp(),
			TxID:           tx.ID(),
			TxOrigin:       signer,
		},
		Topics:   event.Topics,
		Obsolete: obsolete,
	}, nil
}

// EventFilter contains options for contract event filtering.
type EventFilter struct {
	Address *play.Address // restricts matches to events created by specific contracts
	Topic0  *play.Bytes32
	Topic1  *play.Bytes32
	Topic2  *play.Bytes32
	Topic3  *play.Bytes32
	Topic4  *play.Bytes32
}

// Match returs whether event matches filter
func (ef *EventFilter) Match(event *tx.Event) bool {
	if (ef.Address != nil) && (*ef.Address != event.Address) {
		return false
	}

	matchTopic := func(topic *play.Bytes32, index int) bool {
		if topic != nil {
			if len(event.Topics) <= index {
				return false
			}

			if *topic != event.Topics[index] {
				return false
			}
		}
		return true
	}

	return matchTopic(ef.Topic0, 0) &&
		matchTopic(ef.Topic1, 1) &&
		matchTopic(ef.Topic2, 2) &&
		matchTopic(ef.Topic3, 3) &&
		matchTopic(ef.Topic4, 4)
}

// TransferFilter contains options for contract transfer filtering.
type TransferFilter struct {
	TxOrigin  *play.Address // who send transaction
	Sender    *play.Address // who transferred tokens
	Recipient *play.Address // who received tokens
}

// Match returs whether transfer matches filter
func (tf *TransferFilter) Match(transfer *tx.Transfer, origin play.Address) bool {
	if (tf.TxOrigin != nil) && (*tf.TxOrigin != origin) {
		return false
	}

	if (tf.Sender != nil) && (*tf.Sender != transfer.Sender) {
		return false
	}

	if (tf.Recipient != nil) && (*tf.Recipient != transfer.Recipient) {
		return false
	}
	return true
}

type BeatMessage struct {
	Number    uint32       `json:"number"`
	ID        play.Bytes32 `json:"id"`
	ParentID  play.Bytes32 `json:"parentID"`
	Timestamp uint64       `json:"timestamp"`
	Bloom     string       `json:"bloom"`
	K         uint32       `json:"k"`
	Obsolete  bool         `json:"obsolete"`
}
