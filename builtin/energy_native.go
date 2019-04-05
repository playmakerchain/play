// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmakerchain/powerplay/powerplay"
	"github.com/playmakerchain/powerplay/xenv"
)

func init() {
	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_totalSupply", func(env *xenv.Environment) []interface{} {
			env.UseGas(powerplay.SloadGas)
			supply := Energy.Native(env.State(), env.BlockContext().Time).TotalSupply()
			return []interface{}{supply}
		}},
		{"native_totalBurned", func(env *xenv.Environment) []interface{} {
			env.UseGas(powerplay.SloadGas)
			burned := Energy.Native(env.State(), env.BlockContext().Time).TotalBurned()
			return []interface{}{burned}
		}},
		{"native_get", func(env *xenv.Environment) []interface{} {
			var addr common.Address
			env.ParseArgs(&addr)

			env.UseGas(powerplay.GetBalanceGas)
			bal := Energy.Native(env.State(), env.BlockContext().Time).Get(powerplay.Address(addr))
			return []interface{}{bal}
		}},
		{"native_add", func(env *xenv.Environment) []interface{} {
			var args struct {
				Addr   common.Address
				Amount *big.Int
			}
			env.ParseArgs(&args)
			if args.Amount.Sign() == 0 {
				return nil
			}

			env.UseGas(powerplay.GetBalanceGas)
			if env.State().Exists(powerplay.Address(args.Addr)) {
				env.UseGas(powerplay.SstoreResetGas)
			} else {
				env.UseGas(powerplay.SstoreSetGas)
			}
			Energy.Native(env.State(), env.BlockContext().Time).Add(powerplay.Address(args.Addr), args.Amount)
			return nil
		}},
		{"native_sub", func(env *xenv.Environment) []interface{} {
			var args struct {
				Addr   common.Address
				Amount *big.Int
			}
			env.ParseArgs(&args)
			if args.Amount.Sign() == 0 {
				return []interface{}{true}
			}

			env.UseGas(powerplay.GetBalanceGas)
			ok := Energy.Native(env.State(), env.BlockContext().Time).Sub(powerplay.Address(args.Addr), args.Amount)
			if ok {
				env.UseGas(powerplay.SstoreResetGas)
			}
			return []interface{}{ok}
		}},
		{"native_master", func(env *xenv.Environment) []interface{} {
			var addr common.Address
			env.ParseArgs(&addr)

			env.UseGas(powerplay.GetBalanceGas)
			master := env.State().GetMaster(powerplay.Address(addr))
			return []interface{}{master}
		}},
	}
	abi := Energy.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Energy.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
