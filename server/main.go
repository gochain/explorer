package main

import (
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"

	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var ethClient *ethclient.Client
var backendInstance *backend.Backend
var wwwRoot string
var wei = big.NewInt(1000000000000000000)

func parseSkipLimit(r *http.Request) (int, int) {
	limitS := r.URL.Query().Get("limit")
	skipS := r.URL.Query().Get("skip")
	skip := 0
	limit := 100
	if skipS != "" {
		skip, _ = strconv.Atoi(skipS)
	}
	if limitS != "" {
		limit, _ = strconv.Atoi(limitS)
	}

	if skip < 0 {
		skip = 0
	}

	if limit > 500 || limit <= 0 {
		limit = 500
	}
	return skip, limit
}

func parseBlockNumber(r *http.Request) (int, error) {
	bnumS := chi.URLParam(r, "num")
	bnum, err := strconv.Atoi(bnumS)
	if err != nil {
		log.Error().Err(err).Str("bnumS", bnumS).Msg("Error converting bnumS to num")
		return 0, err
	}
	return bnum, nil
}

func main() {
	var rpcUrl string
	var mongoUrl string
	var loglevel string
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rpc-url, u",
			Value:       "https://rpc.gochain.io",
			Usage:       "rpc api url, 'https://rpc.gochain.io'",
			Destination: &rpcUrl,
		},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url, '127.0.0.1:27017'",
			Destination: &mongoUrl,
		},
		cli.StringFlag{
			Name:        "log, l",
			Value:       "info",
			Usage:       "loglevel debug/info/warn/fatal, default is Info",
			Destination: &loglevel,
		},
		cli.StringFlag{
			Name:        "dist, d",
			Value:       "../dist/explorer/",
			Usage:       "folder that should be served",
			Destination: &wwwRoot,
		},
	}

	app.Action = func(c *cli.Context) error {
		level, _ := zerolog.ParseLevel(loglevel)
		zerolog.SetGlobalLevel(level)
		backendInstance = backend.NewBackend(mongoUrl, rpcUrl)
		r := chi.NewRouter()
		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		cors2 := cors.New(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		r.Use(cors2.Handler)
		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		r.Use(middleware.Timeout(60 * time.Second))

		r.Route("/api/stats", func(r chi.Router) {
			r.Get("/", getCurrentStats)
		})
		r.Route("/api/blocks", func(r chi.Router) {
			r.Get("/", getListBlocks)
			r.Get("/{num}", getBlock)
			r.Get("/{num}/transactions", getBlockTransactions)
		})
		r.Route("/api/address", func(r chi.Router) {
			r.Get("/{address}", getAddress)
			r.Get("/{address}/transactions", getAddressTransactions)
			r.Get("/{address}/holders", getTokenHolders)
			r.Get("/{address}/internal_transactions", getInternalTransactions)
		})
		r.Route("/api/transaction", func(r chi.Router) {
			r.Get("/{hash}", getTransaction)
		})

		r.Route("/api/richlist", func(r chi.Router) {
			r.Get("/", getRichlist)
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/totalSupply", getTotalSupply)
			r.Get("/circulatingSupply", getCirculating)

			r.Get("/*", staticHandler)
		})

		http.ListenAndServe(":8080", r)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Run")
	}

}

func getTotalSupply(w http.ResponseWriter, r *http.Request) {
	totalSupply, _ := backendInstance.TotalSupply()
	total := new(big.Rat).SetFrac(totalSupply, wei) // return in GO instead of wei
	w.Write([]byte(total.FloatString(18)))
}

func getCirculating(w http.ResponseWriter, r *http.Request) {
	circulatingSupply, _ := backendInstance.CirculatingSupply()
	circulating := new(big.Rat).SetFrac(circulatingSupply, wei) // return in GO instead of wei
	w.Write([]byte(circulating.FloatString(18)))
}
func staticHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	fileSystemPath := wwwRoot + r.URL.Path
	endURIPath := strings.Split(requestPath, "/")[len(strings.Split(requestPath, "/"))-1]
	splitPath := strings.Split(endURIPath, ".")
	if len(splitPath) > 1 {
		if f, err := os.Stat(fileSystemPath); err == nil && !f.IsDir() {
			http.ServeFile(w, r, fileSystemPath)
			return
		}
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, wwwRoot+"index.html")
}

func getCurrentStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, backendInstance.GetStats())
}

func getRichlist(w http.ResponseWriter, r *http.Request) {
	totalSupply, _ := backendInstance.TotalSupply()
	skip, limit := parseSkipLimit(r)
	circulatingSupply, _ := backendInstance.CirculatingSupply()
	bl := &models.Richlist{
		Rankings:          []*models.Address{},
		TotalSupply:       totalSupply.String(),
		CirculatingSupply: circulatingSupply.String(),
	}
	bl.Rankings = backendInstance.GetRichlist(skip, limit)
	writeJSON(w, http.StatusOK, bl)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	addressHash := chi.URLParam(r, "address")
	log.Info().Str("address", addressHash).Msg("looking up address")
	address := backendInstance.GetAddressByHash(addressHash)
	balance, err := backendInstance.BalanceAt(addressHash, "pending")
	if err == nil {
		address.Balance = balance.String() //to make sure that we are showing most recent balance even if db is outdated
	}
	writeJSON(w, http.StatusOK, address)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	transactionHash := chi.URLParam(r, "hash")
	log.Info().Str("transaction", transactionHash).Msg("looking up transaction")
	transaction := backendInstance.GetTransactionByHash(transactionHash)
	writeJSON(w, http.StatusOK, transaction)
}

func getAddressTransactions(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	skip, limit := parseSkipLimit(r)
	transactions := &models.TransactionList{
		Transactions: []*models.Transaction{},
	}
	transactions.Transactions = backendInstance.GetTransactionList(address, skip, limit)
	writeJSON(w, http.StatusOK, transactions)
}

func getTokenHolders(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	skip, limit := parseSkipLimit(r)
	tokenHolders := &models.TokenHolderList{
		Holders: []*models.TokenHolder{},
	}
	tokenHolders.Holders = backendInstance.GetTokenHoldersList(contractAddress, skip, limit)
	writeJSON(w, http.StatusOK, tokenHolders)
}

func getInternalTransactions(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	skip, limit := parseSkipLimit(r)
	internalTransactions := &models.InternalTransactionsList{
		Transactions: []*models.InternalTransaction{},
	}
	internalTransactions.Transactions = backendInstance.GetInternalTransactionsList(contractAddress, skip, limit)
	writeJSON(w, http.StatusOK, internalTransactions)
}

func getListBlocks(w http.ResponseWriter, r *http.Request) {
	bl := &models.LightBlockList{
		Blocks: []*models.LightBlock{},
	}
	skip, limit := parseSkipLimit(r)
	bl.Blocks = backendInstance.GetLatestsBlocks(skip, limit)
	writeJSON(w, http.StatusOK, bl)
}
func getBlockTransactions(w http.ResponseWriter, r *http.Request) {
	bnum, err := parseBlockNumber(r)
	if err != nil {
		return
	}
	skip, limit := parseSkipLimit(r)
	transactions := &models.TransactionList{
		Transactions: []*models.Transaction{},
	}
	transactions.Transactions = backendInstance.GetBlockTransactionsByNumber(int64(bnum), skip, limit)
	writeJSON(w, http.StatusOK, transactions)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	bnum, err := parseBlockNumber(r)
	if err != nil {
		return
	}
	log.Info().Int("bnum", bnum).Msg("looking up block")
	block := backendInstance.GetBlockByNumber(int64(bnum))
	writeJSON(w, http.StatusOK, block)
}
