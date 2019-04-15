// Copyright (c) 2018 The VeChainThor developers
// COpyright (c) 2019 The PlayMaker developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package api

import (
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/playmakerchain/thor/api/accounts"
	"github.com/playmakerchain/thor/api/blocks"
	"github.com/playmakerchain/thor/api/debug"
	"github.com/playmakerchain/thor/api/doc"
	"github.com/playmakerchain/thor/api/events"
	"github.com/playmakerchain/thor/api/eventslegacy"
	"github.com/playmakerchain/thor/api/node"
	"github.com/playmakerchain/thor/api/subscriptions"
	"github.com/playmakerchain/thor/api/transactions"
	"github.com/playmakerchain/thor/api/transfers"
	"github.com/playmakerchain/thor/api/transferslegacy"
	"github.com/playmakerchain/thor/chain"
	"github.com/playmakerchain/thor/logdb"
	"github.com/playmakerchain/thor/state"
	"github.com/playmakerchain/thor/txpool"
)

//New return api router
func New(chain *chain.Chain, stateCreator *state.Creator, txPool *txpool.TxPool, logDB *logdb.LogDB, nw node.Network, allowedOrigins string, backtraceLimit uint32, callGasLimit uint64) (http.HandlerFunc, func()) {
	origins := strings.Split(strings.TrimSpace(allowedOrigins), ",")
	for i, o := range origins {
		origins[i] = strings.ToLower(strings.TrimSpace(o))
	}

	router := mux.NewRouter()

	// to serve api doc and swagger-ui
	router.PathPrefix("/doc").Handler(
		http.StripPrefix("/doc/", http.FileServer(
			&assetfs.AssetFS{
				Asset:     doc.Asset,
				AssetDir:  doc.AssetDir,
				AssetInfo: doc.AssetInfo})))

	// redirect swagger-ui
	router.Path("/").HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "doc/swagger-ui/", http.StatusTemporaryRedirect)
		})

	accounts.New(chain, stateCreator, callGasLimit).
		Mount(router, "/accounts")
	eventslegacy.New(logDB).
		Mount(router, "/events")
	transferslegacy.New(logDB).
		Mount(router, "/transfers")
	eventslegacy.New(logDB).
		Mount(router, "/logs/events")
	events.New(logDB).
		Mount(router, "/logs/event")
	transferslegacy.New(logDB).
		Mount(router, "/logs/transfers")
	transfers.New(logDB).
		Mount(router, "/logs/transfer")
	blocks.New(chain).
		Mount(router, "/blocks")
	transactions.New(chain, txPool).
		Mount(router, "/transactions")
	debug.New(chain, stateCreator).
		Mount(router, "/debug")
	node.New(nw).
		Mount(router, "/node")
	subs := subscriptions.New(chain, origins, backtraceLimit)
	subs.Mount(router, "/subscriptions")

	handler := handlers.CompressHandler(router)
	handler = handlers.CORS(
		handlers.AllowedOrigins(origins),
		handlers.AllowedHeaders([]string{"content-type"}))(handler)
	return handler.ServeHTTP,
		subs.Close // subscriptions handles hijacked conns, which need to be closed
}
