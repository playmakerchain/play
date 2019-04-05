// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/playmakerchain//"
	"github.com/playmakerchain//xenv"
)

func init() {
	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_executor", func(env *xenv.Environment) []interface{} {
			env.UseGas(.SloadGas)
			addr := .BytesToAddress(Params.Native(env.State()).Get(.KeyExecutorAddress).Bytes())
			return []interface{}{addr}
		}},
		{"native_add", func(env *xenv.Environment) []interface{} {
			var args struct {
				NodeMaster common.Address
				Endorsor   common.Address
				Identity   common.Hash
			}
			env.ParseArgs(&args)

			env.UseGas(.SloadGas)
			ok := Authority.Native(env.State()).Add(
				.Address(args.NodeMaster),
				.Address(args.Endorsor),
				.Bytes32(args.Identity))

			if ok {
				env.UseGas(.SstoreSetGas)
				env.UseGas(.SstoreResetGas)
			}
			return []interface{}{ok}
		}},
		{"native_revoke", func(env *xenv.Environment) []interface{} {
			var nodeMaster common.Address
			env.ParseArgs(&nodeMaster)

			env.UseGas(.SloadGas)
			ok := Authority.Native(env.State()).Revoke(.Address(nodeMaster))
			if ok {
				env.UseGas(.SstoreResetGas * 3)
			}
			return []interface{}{ok}
		}},
		{"native_get", func(env *xenv.Environment) []interface{} {
			var nodeMaster common.Address
			env.ParseArgs(&nodeMaster)

			env.UseGas(.SloadGas * 2)
			listed, endorsor, identity, active := Authority.Native(env.State()).Get(.Address(nodeMaster))

			return []interface{}{listed, endorsor, identity, active}
		}},
		{"native_first", func(env *xenv.Environment) []interface{} {
			env.UseGas(.SloadGas)
			if nodeMaster := Authority.Native(env.State()).First(); nodeMaster != nil {
				return []interface{}{*nodeMaster}
			}
			return []interface{}{.Address{}}
		}},
		{"native_next", func(env *xenv.Environment) []interface{} {
			var nodeMaster common.Address
			env.ParseArgs(&nodeMaster)

			env.UseGas(.SloadGas)
			if next := Authority.Native(env.State()).Next(.Address(nodeMaster)); next != nil {
				return []interface{}{*next}
			}
			return []interface{}{.Address{}}
		}},
		{"native_isEndorsed", func(env *xenv.Environment) []interface{} {
			var nodeMaster common.Address
			env.ParseArgs(&nodeMaster)

			env.UseGas(.SloadGas * 2)
			listed, endorsor, _, _ := Authority.Native(env.State()).Get(.Address(nodeMaster))
			if !listed {
				return []interface{}{false}
			}

			env.UseGas(.GetBalanceGas)
			bal := env.State().GetBalance(endorsor)

			env.UseGas(.SloadGas)
			endorsement := Params.Native(env.State()).Get(.KeyProposerEndorsement)
			return []interface{}{bal.Cmp(endorsement) >= 0}
		}},
	}
	abi := Authority.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Authority.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
