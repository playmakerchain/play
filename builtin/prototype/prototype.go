// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package prototype

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/playmakerchain/powerplay/state"
	"github.com/playmakerchain/powerplay/powerplay"
)

type Prototype struct {
	addr  powerplay.Address
	state *state.State
}

func New(addr powerplay.Address, state *state.State) *Prototype {
	return &Prototype{addr, state}
}

func (p *Prototype) Bind(self powerplay.Address) *Binding {
	return &Binding{p.addr, p.state, self}
}

type Binding struct {
	addr  powerplay.Address
	state *state.State
	self  powerplay.Address
}

func (b *Binding) userKey(user powerplay.Address) powerplay.Bytes32 {
	return powerplay.Blake2b(b.self.Bytes(), user.Bytes(), []byte("user"))
}

func (b *Binding) creditPlanKey() powerplay.Bytes32 {
	return powerplay.Blake2b(b.self.Bytes(), []byte("credit-plan"))
}

func (b *Binding) sponsorKey(sponsor powerplay.Address) powerplay.Bytes32 {
	return powerplay.Blake2b(b.self.Bytes(), sponsor.Bytes(), []byte("sponsor"))
}

func (b *Binding) curSponsorKey() powerplay.Bytes32 {
	return powerplay.Blake2b(b.self.Bytes(), []byte("cur-sponsor"))
}

func (b *Binding) getUserObject(user powerplay.Address) *userObject {
	var uo userObject
	b.state.DecodeStorage(b.addr, b.userKey(user), func(raw []byte) error {
		if len(raw) == 0 {
			uo = userObject{&big.Int{}, 0}
			return nil
		}
		return rlp.DecodeBytes(raw, &uo)
	})
	return &uo
}

func (b *Binding) setUserObject(user powerplay.Address, uo *userObject) {
	b.state.EncodeStorage(b.addr, b.userKey(user), func() ([]byte, error) {
		if uo.IsEmpty() {
			return nil, nil
		}
		return rlp.EncodeToBytes(uo)
	})
}

func (b *Binding) getCreditPlan() *creditPlan {
	var cp creditPlan
	b.state.DecodeStorage(b.addr, b.creditPlanKey(), func(raw []byte) error {
		if len(raw) == 0 {
			cp = creditPlan{&big.Int{}, &big.Int{}}
			return nil
		}
		return rlp.DecodeBytes(raw, &cp)
	})
	return &cp
}

func (b *Binding) setCreditPlan(cp *creditPlan) {
	b.state.EncodeStorage(b.addr, b.creditPlanKey(), func() ([]byte, error) {
		if cp.IsEmpty() {
			return nil, nil
		}
		return rlp.EncodeToBytes(cp)
	})
}

func (b *Binding) IsUser(user powerplay.Address) bool {
	return !b.getUserObject(user).IsEmpty()
}

func (b *Binding) AddUser(user powerplay.Address, blockTime uint64) {
	b.setUserObject(user, &userObject{&big.Int{}, blockTime})
}

func (b *Binding) RemoveUser(user powerplay.Address) {
	// set to empty
	b.setUserObject(user, &userObject{&big.Int{}, 0})
}

func (b *Binding) UserCredit(user powerplay.Address, blockTime uint64) *big.Int {
	uo := b.getUserObject(user)
	if uo.IsEmpty() {
		return &big.Int{}
	}
	return uo.Credit(b.getCreditPlan(), blockTime)
}

func (b *Binding) SetUserCredit(user powerplay.Address, credit *big.Int, blockTime uint64) {
	up := b.getCreditPlan()
	used := new(big.Int).Sub(up.Credit, credit)
	if used.Sign() < 0 {
		used = &big.Int{}
	}
	b.setUserObject(user, &userObject{used, blockTime})
}

func (b *Binding) CreditPlan() (credit, recoveryRate *big.Int) {
	cp := b.getCreditPlan()
	return cp.Credit, cp.RecoveryRate
}

func (b *Binding) SetCreditPlan(credit, recoveryRate *big.Int) {
	b.setCreditPlan(&creditPlan{credit, recoveryRate})
}

func (b *Binding) Sponsor(sponsor powerplay.Address, flag bool) {
	b.state.EncodeStorage(b.addr, b.sponsorKey(sponsor), func() ([]byte, error) {
		if !flag {
			return nil, nil
		}
		return rlp.EncodeToBytes(&flag)
	})
}

func (b *Binding) IsSponsor(sponsor powerplay.Address) (flag bool) {
	b.state.DecodeStorage(b.addr, b.sponsorKey(sponsor), func(raw []byte) error {
		if len(raw) == 0 {
			return nil
		}
		return rlp.DecodeBytes(raw, &flag)
	})
	return
}

func (b *Binding) SelectSponsor(sponsor powerplay.Address) {
	b.state.SetStorage(b.addr, b.curSponsorKey(), powerplay.BytesToBytes32(sponsor.Bytes()))
}

func (b *Binding) CurrentSponsor() powerplay.Address {
	return powerplay.BytesToAddress(b.state.GetStorage(b.addr, b.curSponsorKey()).Bytes())
}
