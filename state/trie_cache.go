// Copyright (c) 2018 The VeChainThor developers
// Copyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package state

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/playmakerchain/powerplay/kv"
	"github.com/playmakerchain/powerplay/powerplay"
	"github.com/playmakerchain/powerplay/trie"
)

var trCache = newTrieCache()

type trieCache struct {
	cache *lru.Cache
}

type trieCacheEntry struct {
	trie *trie.SecureTrie
	kv   kv.GetPutter
}

func newTrieCache() *trieCache {
	cache, _ := lru.New(256)
	return &trieCache{cache: cache}
}

// to get a trie for writing, copy should be set to true
func (tc *trieCache) Get(root powerplay.Bytes32, kv kv.GetPutter, copy bool) (*trie.SecureTrie, error) {

	if v, ok := tc.cache.Get(root); ok {
		entry := v.(*trieCacheEntry)
		if entry.kv == kv {
			if copy {
				return entry.trie.Copy(), nil
			}
			return entry.trie, nil
		}
	}
	tr, err := trie.NewSecure(root, kv, 16)
	if err != nil {
		return nil, err
	}
	tc.cache.Add(root, &trieCacheEntry{tr, kv})
	if copy {
		return tr.Copy(), nil
	}
	return tr, nil
}

func (tc *trieCache) Add(root powerplay.Bytes32, trie *trie.SecureTrie, kv kv.GetPutter) {
	tc.cache.Add(root, &trieCacheEntry{trie.Copy(), kv})
}
