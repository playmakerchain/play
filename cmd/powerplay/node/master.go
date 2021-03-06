// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package node

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/playmakerchain/powerplay/powerplay"
)

type Master struct {
	PrivateKey  *ecdsa.PrivateKey
	Beneficiary *powerplay.Address
}

func (m *Master) Address() powerplay.Address {
	return powerplay.Address(crypto.PubkeyToAddress(m.PrivateKey.PublicKey))
}
