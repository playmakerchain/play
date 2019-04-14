// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package events

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmakerchain/play/api/transactions"
	"github.com/playmakerchain/play/logdb"
	"github.com/playmakerchain/play/play"
)

type TopicSet struct {
	Topic0 *play.Bytes32 `json:"topic0"`
	Topic1 *play.Bytes32 `json:"topic1"`
	Topic2 *play.Bytes32 `json:"topic2"`
	Topic3 *play.Bytes32 `json:"topic3"`
	Topic4 *play.Bytes32 `json:"topic4"`
}

// FilteredEvent only comes from one contract
type FilteredEvent struct {
	Address play.Address         `json:"address"`
	Topics  []*play.Bytes32      `json:"topics"`
	Data    string               `json:"data"`
	Meta    transactions.LogMeta `json:"meta"`
}

//convert a logdb.Event into a json format Event
func convertEvent(event *logdb.Event) *FilteredEvent {
	fe := FilteredEvent{
		Address: event.Address,
		Data:    hexutil.Encode(event.Data),
		Meta: transactions.LogMeta{
			BlockID:        event.BlockID,
			BlockNumber:    event.BlockNumber,
			BlockTimestamp: event.BlockTime,
			TxID:           event.TxID,
			TxOrigin:       event.TxOrigin,
		},
	}
	fe.Topics = make([]*play.Bytes32, 0)
	for i := 0; i < 5; i++ {
		if event.Topics[i] != nil {
			fe.Topics = append(fe.Topics, event.Topics[i])
		}
	}
	return &fe
}

func (e *FilteredEvent) String() string {
	return fmt.Sprintf(`
		Event(
			address: 	   %v,
			topics:        %v,
			data:          %v,
			meta: (blockID     %v,
				blockNumber    %v,
				blockTimestamp %v),
				txID     %v,
				txOrigin %v)
			)`,
		e.Address,
		e.Topics,
		e.Data,
		e.Meta.BlockID,
		e.Meta.BlockNumber,
		e.Meta.BlockTimestamp,
		e.Meta.TxID,
		e.Meta.TxOrigin,
	)
}

type EventCriteria struct {
	Address *play.Address `json:"address"`
	TopicSet
}

type EventFilter struct {
	CriteriaSet []*EventCriteria `json:"criteriaSet"`
	Range       *logdb.Range     `json:"range"`
	Options     *logdb.Options   `json:"options"`
	Order       logdb.Order      `json:"order"`
}

func convertEventFilter(filter *EventFilter) *logdb.EventFilter {
	f := &logdb.EventFilter{
		Range:   filter.Range,
		Options: filter.Options,
		Order:   filter.Order,
	}
	if len(filter.CriteriaSet) > 0 {
		criterias := make([]*logdb.EventCriteria, len(filter.CriteriaSet))
		for i, criteria := range filter.CriteriaSet {
			var topics [5]*play.Bytes32
			topics[0] = criteria.Topic0
			topics[1] = criteria.Topic1
			topics[2] = criteria.Topic2
			topics[3] = criteria.Topic3
			topics[4] = criteria.Topic4
			criteria := &logdb.EventCriteria{
				Address: criteria.Address,
				Topics:  topics,
			}
			criterias[i] = criteria
		}
		f.CriteriaSet = criterias
	}
	return f
}
