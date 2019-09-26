package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gochain-io/gochain/v3/common"
	"github.com/gorilla/schema"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var backendInstance *backend.Backend
var wwwRoot string
var wei = big.NewInt(1000000000000000000)
var reCaptchaSecret string
var rpcUrl string
var logger *zap.Logger

const defaultFetchLimit = 500

func parseTime(r *http.Request) (time.Time, time.Time, error) {
	var err error
	fromTime := time.Unix(0, 0)
	toTime := time.Now()
	fromTimeStr := r.URL.Query().Get("from_time")
	toTimeStr := r.URL.Query().Get("to_time")
	if fromTimeStr != "" {
		fromTime, err = time.Parse(time.RFC3339, fromTimeStr)
		if err != nil {
			return fromTime, toTime, err
		}
	}
	if toTimeStr != "" {
		toTime, err = time.Parse(time.RFC3339, toTimeStr)
		if err != nil {
			return fromTime, toTime, err
		}
	}
	return fromTime, toTime, nil
}

func parseBool(r *http.Request, name string) (boolean *bool) {
	valStr := r.URL.Query().Get(name)
	var val bool
	if valStr != "" {
		if valStr == "true" {
			val = true
		} else {
			val = false
		}
		return &val
	}
	return nil
}

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

	if limit <= 0 {
		limit = defaultFetchLimit
	}
	return skip, limit
}

func parseGetParam(r *http.Request, result interface{}) error {
	return schema.NewDecoder().Decode(result, r.URL.Query())
}

func parseBlockNumber(r *http.Request) (int, error) {
	bnumS := chi.URLParam(r, "num")
	bnum, err := strconv.Atoi(bnumS)
	if err != nil {
		logger.Error("Error converting bnumS to num", zap.Error(err), zap.String("bnumS", bnumS))
		return 0, err
	}
	return bnum, nil
}

func main() {
	var mongoUrl string
	var dbName string
	var signersFile string

	app := cli.NewApp()
	app.Usage = "Server serves the explorer web interface, backed by a mongo database."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rpc-url, u",
			Value:       "https://rpc.gochain.io",
			Usage:       "rpc api url",
			EnvVar:      "RPC_URL",
			Destination: &rpcUrl,
		},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url",
			EnvVar:      "MONGO_URL",
			Destination: &mongoUrl,
		},
		cli.StringFlag{
			Name:        "mongo-dbname, db",
			Value:       "blocks",
			Usage:       "mongo database name",
			EnvVar:      "MONGO_DBNAME",
			Destination: &dbName,
		},
		cli.StringFlag{
			Name:        "signers-file, signers",
			Usage:       "signers file name",
			EnvVar:      "SIGNERS_FILE",
			Value:       "",
			Destination: &signersFile,
		},
		cli.StringFlag{
			Name:        "dist, d",
			Value:       "../dist/explorer/",
			Usage:       "folder that should be served",
			EnvVar:      "DIST",
			Destination: &wwwRoot,
		},
		cli.StringFlag{
			Name:        "recaptcha, r",
			Value:       "",
			Usage:       "secret key for google recaptcha v3",
			EnvVar:      "RECAPTCHA",
			Destination: &reCaptchaSecret,
		},
		cli.StringSliceFlag{
			Name:  "locked-accounts",
			Usage: "accounts with locked funds to exclude from rich list and circulating supply",
		},
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancelFn()
		}
	}()

	app.Action = func(c *cli.Context) error {
		var err error
		cfg := zapdriver.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		logger, err = cfg.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
			os.Exit(1)
		}

		lockedAccounts := c.StringSlice("locked-accounts")
		for i, l := range lockedAccounts {
			if !common.IsHexAddress(l) {
				return fmt.Errorf("invalid hex address: %s", l)
			}
			// Ensure canonical form, since queries are case-sensitive.
			lockedAccounts[i] = common.HexToAddress(l).Hex()
		}

		var signers = make(map[common.Address]models.Signer)

		if signersFile != "" {
			data, err := ioutil.ReadFile(signersFile)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &signers)
			if err != nil {
				return err
			}
		}

		backendInstance, err = backend.NewBackend(ctx, mongoUrl, rpcUrl, dbName, lockedAccounts, signers, logger)
		if err != nil {
			return fmt.Errorf("failed to create backend: %v", err)
		}
		r := chi.NewRouter()
		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		cors2 := cors.New(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
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

		r.Route("/", func(r chi.Router) {
			r.Head("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			r.Get("/totalSupply", getTotalSupply)
			r.Get("/circulatingSupply", getCirculating)
			r.Get("/*", staticHandler)

			r.Route("/api", func(r chi.Router) {
				r.Head("/", pingDB)
				r.Get("/", pingDB)
				r.Post("/verify", verifyContract)
				r.Get("/compiler", getCompilerVersion)
				r.Get("/rpc_provider", getRpcProvider)
				r.Get("/stats", getCurrentStats)
				r.Get("/richlist", getRichlist)

				r.Route("/signers", func(r chi.Router) {
					r.Get("/stats", getSignersStats)
					r.Get("/list", getSignersList)
				})

				r.Route("/blocks", func(r chi.Router) {
					r.Get("/", getListBlocks)
					r.Get("/{num}", getBlock)
					r.Head("/{hash}", checkBlockExist)
					r.Get("/{num}/transactions", getBlockTransactions)
				})

				r.Route("/address", func(r chi.Router) {
					r.Get("/{address}", getAddress)
					r.Get("/{address}/transactions", getAddressTransactions)
					r.Get("/{address}/holders", getTokenHolders)
					r.Get("/{address}/owned_tokens", getOwnedTokens)
					r.Get("/{address}/internal_transactions", getInternalTransactions)
					r.Get("/{address}/contract", getContract)
					r.Get("/{address}/qr", getQr)
				})

				r.Route("/transaction", func(r chi.Router) {
					r.Head("/{hash}", checkTransactionExist)
					r.Get("/{hash}", getTransaction)
				})

				r.Get("/contracts", getContractsList)
			})
		})
		server := &http.Server{Addr: ":8080", Handler: r}
		go func() {
			select {
			case <-ctx.Done():
				server.Close()
			}
		}()
		return server.ListenAndServe()

	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal("Run", zap.Error(err))
	}

}

func getTotalSupply(w http.ResponseWriter, r *http.Request) {
	totalSupply, err := backendInstance.TotalSupply(r.Context())
	if err == nil {
		total := new(big.Rat).SetFrac(totalSupply, wei) // return in GO instead of wei
		w.Write([]byte(total.FloatString(18)))
	} else {
		writeJSON(w, http.StatusInternalServerError, err)
	}
}

func getCirculating(w http.ResponseWriter, r *http.Request) {
	circulatingSupply, err := backendInstance.CirculatingSupply(r.Context())
	if err == nil {
		circulating := new(big.Rat).SetFrac(circulatingSupply, wei) // return in GO instead of wei
		w.Write([]byte(circulating.FloatString(18)))
	} else {
		writeJSON(w, http.StatusInternalServerError, err)
	}

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

func getSignersStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, backendInstance.GetSignersStats())
}

func getSignersList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, backendInstance.GetSignersList())
}

func getRichlist(w http.ResponseWriter, r *http.Request) {
	totalSupply, err := backendInstance.TotalSupply(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	skip, limit := parseSkipLimit(r)
	circulatingSupply, err := backendInstance.CirculatingSupply(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	bl := &models.Richlist{
		Rankings:          []*models.Address{},
		TotalSupply:       new(big.Rat).SetFrac(totalSupply, wei).FloatString(18),
		CirculatingSupply: new(big.Rat).SetFrac(circulatingSupply, wei).FloatString(18),
	}
	bl.Rankings = backendInstance.GetRichlist(skip, limit)
	writeJSON(w, http.StatusOK, bl)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	addressHash := chi.URLParam(r, "address")
	logger.Info("looking up address", zap.String("address", addressHash))
	address, err := backendInstance.GetAddressByHash(r.Context(), addressHash)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, address)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	transactionHash := chi.URLParam(r, "hash")
	logger.Info("looking up transaction", zap.String("transaction", transactionHash))
	transaction := backendInstance.GetTransactionByHash(r.Context(), transactionHash)
	writeJSON(w, http.StatusOK, transaction)
}

func checkTransactionExist(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	tx := backendInstance.GetTransactionByHash(r.Context(), hash)
	if tx != nil {
		writeJSON(w, http.StatusOK, nil)
	} else {
		writeJSON(w, http.StatusNotFound, nil)
	}
}

func getAddressTransactions(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	inputDataEmpty := parseBool(r, "input_data_empty")
	skip, limit := parseSkipLimit(r)
	fromTime, toTime, err := parseTime(r)
	if err == nil {
		transactions := &models.TransactionList{
			Transactions: []*models.Transaction{},
		}
		transactions.Transactions = backendInstance.GetTransactionList(address, skip, limit, fromTime, toTime, inputDataEmpty)
		writeJSON(w, http.StatusOK, transactions)
	} else {
		logger.Info("getAddressTransactions", zap.Error(err))
		errorResponse(w, http.StatusBadRequest, err)
	}
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

func getOwnedTokens(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	skip, limit := parseSkipLimit(r)
	tokens := &models.OwnedTokenList{
		OwnedTokens: []*models.TokenHolder{},
	}
	tokens.OwnedTokens = backendInstance.GetOwnedTokensList(contractAddress, skip, limit)
	writeJSON(w, http.StatusOK, tokens)
}

func getInternalTransactions(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	tokenTransactions := false
	token_transactions_param := r.URL.Query().Get("token_transactions")
	if token_transactions_param != "" && token_transactions_param != "false" {
		tokenTransactions = true
	}
	skip, limit := parseSkipLimit(r)
	internalTransactions := &models.InternalTransactionsList{
		Transactions: []*models.InternalTransaction{},
	}
	internalTransactions.Transactions = backendInstance.GetInternalTransactionsList(contractAddress, tokenTransactions, skip, limit)
	writeJSON(w, http.StatusOK, internalTransactions)
}

func getContract(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	contract := backendInstance.GetContract(contractAddress)
	writeJSON(w, http.StatusOK, contract)
}

func getQr(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	var png []byte
	png, err := qrcode.Encode(contractAddress, qrcode.Medium, 256)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	writeFile(w, http.StatusOK, "image/png", png)
}

func verifyContract(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contractData *models.Contract
	err := decoder.Decode(&contractData)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	/*if contractData.RecaptchaToken == "" {
		err := errors.New("recaptcha token is empty")
		errorResponse(w, http.StatusBadRequest, err)
		return
	}*/
	if contractData.Address == "" || contractData.ContractName == "" || contractData.SourceCode == "" || contractData.CompilerVersion == "" {
		err := errors.New("required field is empty")
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	if len(contractData.Address) != 42 {
		err := errors.New("contract address is wrong")
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	compilerVersions, err := backendInstance.GetCompilerVersion()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	compilerOk := false
	for _, compiler := range compilerVersions {
		if contractData.CompilerVersion == compiler {
			compilerOk = true
		}
	}

	if compilerOk != true {
		err := errors.New("wrong compiler version")
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	/*err = verifyReCaptcha(contractData.RecaptchaToken, reCaptchaSecret, "contractVerification", r.RemoteAddr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}*/
	result, err := backendInstance.VerifyContract(r.Context(), contractData)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}

func getCompilerVersion(w http.ResponseWriter, r *http.Request) {
	result, err := backendInstance.GetCompilerVersion()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func getRpcProvider(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, rpcUrl)
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
	var block *models.Block
	if err != nil {
		hash := chi.URLParam(r, "num")
		logger.Info("failed to parse number of the block so assuming it's hash", zap.String("hash", hash))
		block = backendInstance.GetBlockByHash(hash)
	} else {
		logger.Info("looking up block", zap.Int("bnum", bnum))
		block = backendInstance.GetBlockByNumber(r.Context(), int64(bnum))
	}
	writeJSON(w, http.StatusOK, block)
}

func checkBlockExist(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	block := backendInstance.GetBlockByHash(hash)
	if block != nil {
		writeJSON(w, http.StatusOK, nil)
	} else {
		writeJSON(w, http.StatusNotFound, nil)
	}
}

func pingDB(w http.ResponseWriter, r *http.Request) {
	err := backendInstance.PingDB()
	if err != nil {
		logger.Error("Cannot ping DB", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func getContractsList(w http.ResponseWriter, r *http.Request) {
	filter := new(models.ContractsFilter)
	if err := parseGetParam(r, filter); err != nil {
		logger.Info("failed to parse get params")
		errorResponse(w, http.StatusBadRequest, errors.New("invalid params"))
	}
	addresses := backendInstance.GetContracts(filter)
	writeJSON(w, http.StatusOK, addresses)
}
