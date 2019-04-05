// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin_test

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/playmakerchain//builtin"
	"github.com/playmakerchain//chain"
	"github.com/playmakerchain//lvldb"
	"github.com/playmakerchain//runtime"
	"github.com/playmakerchain//state"
	"github.com/playmakerchain//"
	"github.com/playmakerchain//tx"
	"github.com/playmakerchain//xenv"
)

func approverEvent(approver .Address, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("Approver")
	var b32 .Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  [].Bytes32{ev.ID(), .BytesToBytes32(approver.Bytes())},
		Data:    data,
	}
}
func proposalEvent(id .Bytes32, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("Proposal")
	var b32 .Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  [].Bytes32{ev.ID(), id},
		Data:    data,
	}
}

func votingContractEvent(addr .Address, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("VotingContract")
	var b32 .Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  [].Bytes32{ev.ID(), .BytesToBytes32(addr.Bytes())},
		Data:    data,
	}
}

func initExectorTest() *ctest {
	kv, _ := lvldb.NewMem()
	b0 := buildGenesis(kv, func(state *state.State) error {
		state.SetCode(builtin.Prototype.Address, builtin.Prototype.RuntimeBytecodes())
		state.SetCode(builtin.Executor.Address, builtin.Executor.RuntimeBytecodes())
		state.SetCode(builtin.Params.Address, builtin.Params.RuntimeBytecodes())
		builtin.Params.Native(state).Set(.KeyExecutorAddress, new(big.Int).SetBytes(builtin.Executor.Address[:]))
		return nil
	})

	c, _ := chain.New(kv, b0)
	st, _ := state.New(b0.Header().StateRoot(), kv)
	seeker := c.NewSeeker(b0.Header().ID())

	rt := runtime.New(seeker, st, &xenv.BlockContext{Time: uint64(time.Now().Unix())})

	return &ctest{
		rt:  rt,
		abi: builtin.Executor.ABI,
		to:  builtin.Executor.Address,
	}
}

func TestExecutorApprover(t *testing.T) {
	test := initExectorTest()
	var approvers [].Address
	for i := 0; i < 7; i++ {
		approvers = append(approvers, .BytesToAddress([]byte(fmt.Sprintf("approver%d", i))))
	}

	for _, a := range approvers {
		// zero identity
		test.Case("addApprover", a, .Bytes32{}).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("addApprover", a, .BytesToBytes32(a.Bytes())).
			Caller(.BytesToAddress([]byte("other"))).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("addApprover", a, .BytesToBytes32(a.Bytes())).
			Caller(builtin.Executor.Address).
			ShouldLog(approverEvent(a, "added")).
			Assert(t)
		assert.True(t, builtin.Prototype.Native(test.rt.State()).Bind(test.to).IsUser(a))
	}

	test.Case("approverCount").
		ShouldOutput(uint8(len(approvers))).
		Assert(t)

	test.Case("addApprover", approvers[0], .BytesToBytes32(approvers[0].Bytes())).
		Caller(builtin.Executor.Address).
		ShouldVMError(errReverted).
		Assert(t)

	for _, a := range approvers {
		test.Case("approvers", a).
			ShouldOutput(.BytesToBytes32(a.Bytes()), true).
			Assert(t)
	}

	for _, a := range approvers {
		test.Case("revokeApprover", a).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("revokeApprover", a).
			Caller(builtin.Executor.Address).
			ShouldLog(approverEvent(a, "revoked")).
			Assert(t)
		assert.False(t, builtin.Prototype.Native(test.rt.State()).Bind(test.to).IsUser(a))
	}
	test.Case("approverCount").
		ShouldOutput(uint8(0)).
		Assert(t)
}

func TestExecutorVotingContract(t *testing.T) {

	test := initExectorTest()
	voting := .BytesToAddress([]byte("voting"))
	test.Case("attachVotingContract", voting).
		ShouldVMError(errReverted).
		Assert(t)
	test.Case("votingContracts", voting).
		ShouldOutput(false).
		Assert(t)
	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "attached")).
		Assert(t)

	test.Case("votingContracts", voting).
		ShouldOutput(true).
		Assert(t)

	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("detachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "detached")).
		Assert(t)

	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "attached")).
		Assert(t)
}

func TestExecutorProposal(t *testing.T) {
	test := initExectorTest()

	target := builtin.Params.Address
	setParam, _ := builtin.Params.ABI.MethodByName("set")
	data, _ := setParam.EncodeInput(.BytesToBytes32([]byte("paramKey")), big.NewInt(12345))
	test.Case("propose", target, data).
		ShouldVMError(errReverted).
		Assert(t)

	approver := .BytesToAddress([]byte("approver"))
	test.Case("addApprover", approver, .BytesToBytes32(approver.Bytes())).
		Caller(builtin.Executor.Address).
		Assert(t)

	proposalID := func() .Bytes32 {
		var b8 [8]byte
		binary.BigEndian.PutUint64(b8[:], test.rt.Context().Time)
		return .Bytes32(crypto.Keccak256Hash(b8[:], approver[:]))
	}()
	test.Case("propose", target, data).
		Caller(approver).
		ShouldOutput(proposalID).
		ShouldLog(proposalEvent(proposalID, "proposed")).
		Assert(t)

	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(0),
			false,
			target,
			data).
		Assert(t)

	test.Case("approve", proposalID).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("approve", proposalID).
		Caller(approver).
		ShouldLog(proposalEvent(proposalID, "approved")).
		Assert(t)
	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(1),
			false,
			target,
			data).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldLog(proposalEvent(proposalID, "executed")).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldVMError(errReverted).
		Assert(t)
	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(1),
			true,
			target,
			data).
		Assert(t)

	assert.Equal(t, builtin.Params.Native(test.rt.State()).Get(.BytesToBytes32([]byte("paramKey"))), big.NewInt(12345))
}
