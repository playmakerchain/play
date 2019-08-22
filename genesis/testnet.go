// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package genesis

import (
	"math/big"

	"github.com/playmakerchain/powerplay/builtin"
	"github.com/playmakerchain/powerplay/state"
	"github.com/playmakerchain/powerplay/powerplay"
	"github.com/playmakerchain/powerplay/tx"
	"github.com/playmakerchain/powerplay/vm"
)

// NewTestnet create genesis for testnet.
func NewTestnet() *Genesis {
	launchTime := uint64(1565869157) // 'Tue Jun 26 2018 20:00:00 GMT+0800 (CST)'

	// use this address as executor instead of builtin one, for test purpose
	executor, _ := powerplay.ParseAddress("0xB5A34b62b63A6f1EE99DFD30b133B657859f8d79")
	acccount0, _ := powerplay.ParseAddress("0xe59D475Abe695c7f67a8a2321f33A856B0B4c71d")

	master0, _ := powerplay.ParseAddress("0x25AE0ef84dA4a76D5a1DFE80D3789C2c46FeE30a")
	endorser0, _ := powerplay.ParseAddress("0xb4094c25f86d628fdD571Afc4077f0d0196afB48")

	builder := new(Builder).
		Timestamp(launchTime).
		GasLimit(powerplay.InitialGasLimit).
		State(func(state *state.State) error {
			tokenSupply := new(big.Int)

			// alloc precompiled contracts
			for addr := range vm.PrecompiledContractsByzantium {
				state.SetCode(powerplay.Address(addr), emptyRuntimeBytecode)
			}

			// setup builtin contracts
			state.SetCode(builtin.Authority.Address, builtin.Authority.RuntimeBytecodes())
			state.SetCode(builtin.Energy.Address, builtin.Energy.RuntimeBytecodes())
			state.SetCode(builtin.Params.Address, builtin.Params.RuntimeBytecodes())
			state.SetCode(builtin.Prototype.Address, builtin.Prototype.RuntimeBytecodes())
			state.SetCode(builtin.Extension.Address, builtin.Extension.RuntimeBytecodes())

			// 50 billion for account0
			amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(50*1000*1000*1000))
			state.SetBalance(acccount0, amount)
			state.SetEnergy(acccount0, &big.Int{}, launchTime)
			tokenSupply.Add(tokenSupply, amount)

			// 25 million for endorser0
			amount = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(25*1000*1000))
			state.SetBalance(endorser0, amount)
			state.SetEnergy(endorser0, &big.Int{}, launchTime)
			tokenSupply.Add(tokenSupply, amount)

			builtin.Energy.Native(state, launchTime).SetInitialSupply(tokenSupply, &big.Int{})
			return nil
		}).
		// set initial params
		// use an external account as executor to manage testnet easily
		Call(
			tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", powerplay.KeyExecutorAddress, new(big.Int).SetBytes(executor[:]))),
			powerplay.Address{}).
		Call(
			tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", powerplay.KeyRewardRatio, powerplay.InitialRewardRatio)),
			executor).
		Call(
			tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", powerplay.KeyBaseGasPrice, powerplay.InitialBaseGasPrice)),
			executor).
		Call(
			tx.NewClause(&builtin.Params.Address).WithData(mustEncodeInput(builtin.Params.ABI, "set", powerplay.KeyProposerEndorsement, powerplay.InitialProposerEndorsement)),
			executor).
		// add master0 as the initial block proposer
		Call(tx.NewClause(&builtin.Authority.Address).WithData(mustEncodeInput(builtin.Authority.ABI, "add", master0, endorser0, powerplay.BytesToBytes32([]byte("master0")))),
			executor)

	id, err := builder.ComputeID()
	if err != nil {
		panic(err)
	}
	return &Genesis{builder, id, "testnet"}
}
