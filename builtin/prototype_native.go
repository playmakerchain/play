// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmakerchain/powerplay/abi"
	"github.com/playmakerchain/powerplay/powerplay"
	"github.com/playmakerchain/powerplay/xenv"
)

func init() {

	events := Prototype.Events()

	mustEventByName := func(name string) *abi.Event {
		if event, found := events.EventByName(name); found {
			return event
		}
		panic("event not found")
	}

	masterEvent := mustEventByName("$Master")
	creditPlanEvent := mustEventByName("$CreditPlan")
	userEvent := mustEventByName("$User")
	sponsorEvent := mustEventByName("$Sponsor")

	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_master", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)

			env.UseGas(powerplay.GetBalanceGas)
			master := env.State().GetMaster(powerplay.Address(self))

			return []interface{}{master}
		}},
		{"native_setMaster", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self      common.Address
				NewMaster common.Address
			}
			env.ParseArgs(&args)

			env.UseGas(powerplay.SstoreResetGas)
			env.State().SetMaster(powerplay.Address(args.Self), powerplay.Address(args.NewMaster))

			env.Log(masterEvent, powerplay.Address(args.Self), nil, args.NewMaster)
			return nil
		}},
		{"native_balanceAtBlock", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self        common.Address
				BlockNumber uint32
			}
			env.ParseArgs(&args)
			ctx := env.BlockContext()

			if args.BlockNumber > ctx.Number {
				return []interface{}{&big.Int{}}
			}

			if ctx.Number-args.BlockNumber > powerplay.MaxBackTrackingBlockNumber {
				return []interface{}{&big.Int{}}
			}

			if args.BlockNumber == ctx.Number {
				env.UseGas(powerplay.GetBalanceGas)
				val := env.State().GetBalance(powerplay.Address(args.Self))
				return []interface{}{val}
			}

			env.UseGas(powerplay.SloadGas)
			blockID := env.Seeker().GetID(args.BlockNumber)

			env.UseGas(powerplay.SloadGas)
			header := env.Seeker().GetHeader(blockID)

			env.UseGas(powerplay.SloadGas)
			state := env.State().Spawn(header.StateRoot())

			env.UseGas(powerplay.GetBalanceGas)
			val := state.GetBalance(powerplay.Address(args.Self))

			return []interface{}{val}
		}},
		{"native_energyAtBlock", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self        common.Address
				BlockNumber uint32
			}
			env.ParseArgs(&args)
			ctx := env.BlockContext()
			if args.BlockNumber > ctx.Number {
				return []interface{}{&big.Int{}}
			}

			if ctx.Number-args.BlockNumber > powerplay.MaxBackTrackingBlockNumber {
				return []interface{}{&big.Int{}}
			}

			if args.BlockNumber == ctx.Number {
				env.UseGas(powerplay.GetBalanceGas)
				val := env.State().GetEnergy(powerplay.Address(args.Self), ctx.Time)
				return []interface{}{val}
			}

			env.UseGas(powerplay.SloadGas)
			blockID := env.Seeker().GetID(args.BlockNumber)

			env.UseGas(powerplay.SloadGas)
			header := env.Seeker().GetHeader(blockID)

			env.UseGas(powerplay.SloadGas)
			state := env.State().Spawn(header.StateRoot())

			env.UseGas(powerplay.GetBalanceGas)
			val := state.GetEnergy(powerplay.Address(args.Self), header.Timestamp())

			return []interface{}{val}
		}},
		{"native_hasCode", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)

			env.UseGas(powerplay.GetBalanceGas)
			hasCode := !env.State().GetCodeHash(powerplay.Address(self)).IsZero()

			return []interface{}{hasCode}
		}},
		{"native_storageFor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				Key  powerplay.Bytes32
			}
			env.ParseArgs(&args)

			env.UseGas(powerplay.SloadGas)
			storage := env.State().GetStorage(powerplay.Address(args.Self), args.Key)
			return []interface{}{storage}
		}},
		{"native_creditPlan", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(self))

			env.UseGas(powerplay.SloadGas)
			credit, rate := binding.CreditPlan()

			return []interface{}{credit, rate}
		}},
		{"native_setCreditPlan", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self         common.Address
				Credit       *big.Int
				RecoveryRate *big.Int
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SstoreSetGas)
			binding.SetCreditPlan(args.Credit, args.RecoveryRate)
			env.Log(creditPlanEvent, powerplay.Address(args.Self), nil, args.Credit, args.RecoveryRate)
			return nil
		}},
		{"native_isUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			isUser := binding.IsUser(powerplay.Address(args.User))

			return []interface{}{isUser}
		}},
		{"native_userCredit", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(2 * powerplay.SloadGas)
			credit := binding.UserCredit(powerplay.Address(args.User), env.BlockContext().Time)

			return []interface{}{credit}
		}},
		{"native_addUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			if binding.IsUser(powerplay.Address(args.User)) {
				return []interface{}{false}
			}

			env.UseGas(powerplay.SstoreSetGas)
			binding.AddUser(powerplay.Address(args.User), env.BlockContext().Time)

			var action powerplay.Bytes32
			copy(action[:], "added")
			env.Log(userEvent, powerplay.Address(args.Self), []powerplay.Bytes32{powerplay.BytesToBytes32(args.User[:])}, action)
			return []interface{}{true}
		}},
		{"native_removeUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			if !binding.IsUser(powerplay.Address(args.User)) {
				return []interface{}{false}
			}

			env.UseGas(powerplay.SstoreResetGas)
			binding.RemoveUser(powerplay.Address(args.User))

			var action powerplay.Bytes32
			copy(action[:], "removed")
			env.Log(userEvent, powerplay.Address(args.Self), []powerplay.Bytes32{powerplay.BytesToBytes32(args.User[:])}, action)
			return []interface{}{true}
		}},
		{"native_sponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			if binding.IsSponsor(powerplay.Address(args.Sponsor)) {
				return []interface{}{false}
			}

			env.UseGas(powerplay.SstoreSetGas)
			binding.Sponsor(powerplay.Address(args.Sponsor), true)

			var action powerplay.Bytes32
			copy(action[:], "sponsored")
			env.Log(sponsorEvent, powerplay.Address(args.Self), []powerplay.Bytes32{powerplay.BytesToBytes32(args.Sponsor.Bytes())}, action)
			return []interface{}{true}
		}},
		{"native_unsponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			if !binding.IsSponsor(powerplay.Address(args.Sponsor)) {
				return []interface{}{false}
			}

			env.UseGas(powerplay.SstoreResetGas)
			binding.Sponsor(powerplay.Address(args.Sponsor), false)

			var action powerplay.Bytes32
			copy(action[:], "unsponsored")
			env.Log(sponsorEvent, powerplay.Address(args.Self), []powerplay.Bytes32{powerplay.BytesToBytes32(args.Sponsor.Bytes())}, action)
			return []interface{}{true}
		}},
		{"native_isSponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			isSponsor := binding.IsSponsor(powerplay.Address(args.Sponsor))

			return []interface{}{isSponsor}
		}},
		{"native_selectSponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(args.Self))

			env.UseGas(powerplay.SloadGas)
			if !binding.IsSponsor(powerplay.Address(args.Sponsor)) {
				return []interface{}{false}
			}

			env.UseGas(powerplay.SstoreResetGas)
			binding.SelectSponsor(powerplay.Address(args.Sponsor))

			var action powerplay.Bytes32
			copy(action[:], "selected")
			env.Log(sponsorEvent, powerplay.Address(args.Self), []powerplay.Bytes32{powerplay.BytesToBytes32(args.Sponsor.Bytes())}, action)

			return []interface{}{true}
		}},
		{"native_currentSponsor", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)
			binding := Prototype.Native(env.State()).Bind(powerplay.Address(self))

			env.UseGas(powerplay.SloadGas)
			addr := binding.CurrentSponsor()

			return []interface{}{addr}
		}},
	}
	abi := Prototype.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Prototype.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
