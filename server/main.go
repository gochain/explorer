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

	"github.com/blendle/zapdriver"
	"github.com/dgraph-io/ristretto"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain/gochain/v4/common"
	"github.com/gorilla/schema"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var backendInstance *backend.Backend
var wwwRoot string
var wei = big.NewInt(1000000000000000000)
var reCaptchaSecret string
var rpcUrl string
var logger *zap.Logger

type Params interface {
	Sanitize()
}

func parseGetParam(r *http.Request, w http.ResponseWriter, params Params) bool {
	if err := schema.NewDecoder().Decode(params, r.URL.Query()); err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return false
	}
	params.Sanitize()
	return true
}

func main() {
	cfg := zapdriver.NewProductionConfig()
	var err error
	logger, err = zapdriver.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Fatal panic", zap.String("panic", fmt.Sprintf("%+v", r)))
		}
	}()

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
		cli.StringFlag{
			Name:  "log-level",
			Usage: "Minimum log level to include. Lower levels will be discarded. (debug, info, warn, error, dpanic, panic, fatal)",
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
		if c.IsSet("log-level") {
			var lvl zapcore.Level
			s := c.String("log-level")
			if err := lvl.Set(s); err != nil {
				return fmt.Errorf("invalid log-level %q: %v", s, err)
			}
			cfg.Level.SetLevel(lvl)
		}
		lockedAccounts := c.StringSlice("locked-accounts")
		initialAllocation := new(big.Int)
		for i, l := range lockedAccounts {
			if !common.IsHexAddress(l) {
				// Local accounts with total supply provided
				s := strings.Split(l, ":")
				if len(s) != 2 {
					return fmt.Errorf("locked accounts wrong format: %s", l)
				}
				if !common.IsHexAddress(s[0]) {
					return fmt.Errorf("invalid hex address: %s", s[0])
				}
				lockedAccounts[i] = s[0]
				v := new(big.Int)
				v.SetString(s[1], 16)
				initialAllocation = initialAllocation.Add(initialAllocation, v)
			} else {
				// Ensure canonical form, since queries are case-sensitive.
				lockedAccounts[i] = common.HexToAddress(l).Hex()
			}
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

		// init memory cache
		cache, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e6,    // number of keys to track frequency of.
			MaxCost:     100000, // maximum cost of cache.
			BufferItems: 64,     // number of keys per Get buffer.
		})
		if err != nil {
			panic(err)
		}

		backendInstance, err = backend.NewBackend(ctx, mongoUrl, rpcUrl, dbName, lockedAccounts, signers, initialAllocation, logger, cache)
		if err != nil {
			return fmt.Errorf("failed to create backend: %v", err)
		}
		r := chi.NewRouter()
		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		// r.Use(middleware.RequestLogger(&zapLogFormatter{logger}))
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
				r.Get("/supply", getSupplyStats)
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
					r.Get("/{address}/tx/{nonce}", getAddressTxByNonce)
					r.Get("/{address}/tx/{nonce}/hash", getAddressTxHashByNonce)
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
				if err := server.Close(); err != nil {
					logger.Warn("Error closing server", zap.Error(err))
				}
			}
		}()
		switch err := server.ListenAndServe(); err {
		case nil, http.ErrServerClosed:
			return nil
		default:
			return err
		}

	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("Fatal error", zap.Error(err))
	}
	logger.Info("Stopping")
}

func getTotalSupply(w http.ResponseWriter, r *http.Request) {
	totalSupply, _, err := backendInstance.TotalSupply(r.Context())
	if err == nil {
		total := new(big.Rat).SetFrac(totalSupply, wei) // return in GO instead of wei
		w.Write([]byte(total.FloatString(18)))
	} else {
		logger.Error("Failed to get total supply", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, err)
	}
}

func getCirculating(w http.ResponseWriter, r *http.Request) {
	circulatingSupply, err := backendInstance.CirculatingSupply(r.Context())
	if err == nil {
		circulating := new(big.Rat).SetFrac(circulatingSupply, wei) // return in GO instead of wei
		w.Write([]byte(circulating.FloatString(18)))
	} else {
		logger.Error("Failed to get circulating supply", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, err)
	}
}

func getSupplyStats(w http.ResponseWriter, r *http.Request) {
	supplyStats, err := backendInstance.SupplyStats(r.Context())
	if err != nil {
		logger.Error("Failed to get supply stats", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, supplyStats)
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

func getCurrentStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := backendInstance.GetStats()
	if err != nil {
		logger.Error("Failed to get stats", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func getSignersStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := backendInstance.GetSignersStats()
	if err != nil {
		logger.Error("Failed to get signer stats", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func getSignersList(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, backendInstance.GetSignersList())
}

func getRichlist(w http.ResponseWriter, r *http.Request) {
	filter := new(models.PaginationFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	totalSupply, _, err := backendInstance.TotalSupply(r.Context())
	if err != nil {
		logger.Error("Failed to get total supply", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	circulatingSupply, err := backendInstance.CirculatingSupply(r.Context())
	if err != nil {
		logger.Error("Failed to get circulating supply", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	bl := &models.Richlist{
		Rankings:          []*models.Address{},
		TotalSupply:       new(big.Rat).SetFrac(totalSupply, wei).FloatString(18),
		CirculatingSupply: new(big.Rat).SetFrac(circulatingSupply, wei).FloatString(18),
	}
	bl.Rankings, err = backendInstance.GetRichlist(filter)
	if err != nil {
		logger.Error("Failed to get rich list", zap.Int("skip", filter.Skip), zap.Int("limit", filter.Limit), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, bl)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
	addressHash := chi.URLParam(r, "address")
	if !common.IsHexAddress(addressHash) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", addressHash))
		return
	}
	address, err := backendInstance.GetAddressByHash(r.Context(), addressHash)
	if err != nil {
		logger.Error("Failed to get address", zap.String("address", addressHash), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, address)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	transactionHash := chi.URLParam(r, "hash")
	transaction, err := backendInstance.GetTransactionByHash(r.Context(), transactionHash)
	if err != nil {
		logger.Error("Failed to get tx", zap.String("address", transactionHash), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if transaction == nil {
		writeJSON(w, http.StatusNotFound, nil)
	}
	writeJSON(w, http.StatusOK, transaction)
}

func checkTransactionExist(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	tx, err := backendInstance.GetTransactionByHash(r.Context(), hash)
	if err != nil {
		logger.Error("Failed to get tx", zap.String("address", hash), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if tx != nil {
		writeJSON(w, http.StatusOK, nil)
	} else {
		writeJSON(w, http.StatusNotFound, nil)
	}
}

func getAddressTransactions(w http.ResponseWriter, r *http.Request) {
	var err error
	address := chi.URLParam(r, "address")
	if !common.IsHexAddress(address) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", address))
		return
	}
	filter := new(models.TxsFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	transactions := &models.TransactionList{}
	transactions.Transactions, err = backendInstance.GetTransactionList(address, filter)
	if err != nil {
		logger.Error("Failed to get txs", zap.String("address", address), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func getAddressTxByNonce(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	if !common.IsHexAddress(address) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", address))
		return
	}
	nonceS := chi.URLParam(r, "nonce")
	nonce, err := strconv.ParseInt(nonceS, 0, 64)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'nonce' %q: %v", nonce, err))
		return
	}
	tx, err := backendInstance.GetTxByAddressAndNonce(r.Context(), address, nonce)
	if err != nil {
		logger.Error("Failed to get tx", zap.String("address", address), zap.Int64("nonce", nonce), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if tx == nil {
		writeJSON(w, http.StatusNotFound, nil)
	}
	writeJSON(w, http.StatusOK, tx)
}

func getAddressTxHashByNonce(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	if !common.IsHexAddress(address) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", address))
		return
	}
	nonceS := chi.URLParam(r, "nonce")
	nonce, err := strconv.ParseInt(nonceS, 0, 64)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'nonce' %q: %v", nonce, err))
		return
	}
	tx, err := backendInstance.GetTxByAddressAndNonce(r.Context(), address, nonce)
	if err != nil {
		logger.Error("Failed to get tx", zap.String("address", address), zap.Int64("nonce", nonce), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if tx == nil {
		writeJSON(w, http.StatusNotFound, nil)
	}
	writeJSON(w, http.StatusOK, tx.TxHash)
}

func getTokenHolders(w http.ResponseWriter, r *http.Request) {
	var err error
	contractAddress := chi.URLParam(r, "address")
	if !common.IsHexAddress(contractAddress) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", contractAddress))
		return
	}
	filter := new(models.PaginationFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	skipLimit := 5000
	if filter.Skip > skipLimit {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("skip cannot be higher than %d for this endpoint", skipLimit))
		return
	}
	tokenHolders := &models.TokenHolderList{}
	tokenHolders.Holders, err = backendInstance.GetTokenHoldersList(contractAddress, filter)
	if err != nil {
		logger.Error("Failed to get token holders", zap.String("address", contractAddress), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, errors.New("Failed to get token holders"))
		return
	}
	writeJSON(w, http.StatusOK, tokenHolders)
}

func getOwnedTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	contractAddress := chi.URLParam(r, "address")
	if !common.IsHexAddress(contractAddress) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", contractAddress))
		return
	}
	filter := new(models.PaginationFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	tokens := &models.OwnedTokenList{}
	tokens.OwnedTokens, err = backendInstance.GetOwnedTokensList(ctx, contractAddress, filter)
	if err != nil {
		logger.Error("Failed to get owned tokens", zap.String("address", contractAddress), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func getInternalTransactions(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	if !common.IsHexAddress(contractAddress) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", contractAddress))
		return
	}
	filter := new(models.InternalTxFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	tokenTransfers := &models.TokenTransfers{}
	var err error
	if filter.TokenTransactions {
		tokenTransfers.Transfers, err = backendInstance.GetHeldTokenTransfers(contractAddress, &filter.PaginationFilter)
		if err != nil {
			logger.Error("Failed to get contract's held token transfers", zap.String("address", contractAddress), zap.Error(err))
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		if filter.InternalAddress != "" && filter.TokenID != "" {
			errorResponse(w, http.StatusBadRequest, fmt.Errorf("only one of internal_address and token_id may be used"))
			return
		}
		tokenTransfers.Transfers, err = backendInstance.GetInternalTokenTransfers(contractAddress, filter)
		if err != nil {
			logger.Error("Failed to get contract's internal token transfers", zap.String("address", contractAddress), zap.Error(err))
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, tokenTransfers)
}

func getContract(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	if !common.IsHexAddress(contractAddress) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", contractAddress))
		return
	}
	contract, err := backendInstance.GetContract(contractAddress)
	if err != nil {
		logger.Error("Failed to get contract", zap.String("address", contractAddress), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if contract == nil {
		writeJSON(w, http.StatusNotFound, nil)
		return
	}
	writeJSON(w, http.StatusOK, contract)
}

func getQr(w http.ResponseWriter, r *http.Request) {
	contractAddress := chi.URLParam(r, "address")
	if !common.IsHexAddress(contractAddress) {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid hex 'address': %s", contractAddress))
		return
	}
	var png []byte
	png, err := qrcode.Encode(contractAddress, qrcode.Medium, 256)
	if err != nil {
		logger.Error("Failed to encode qrcode", zap.String("address", contractAddress), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
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
	var err error
	filter := new(models.PaginationFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	bl := &models.LightBlockList{}
	bl.Blocks, err = backendInstance.GetLatestsBlocks(filter)
	if err != nil {
		logger.Error(
			"Failed to get latest blocks",
			zap.Int("skip", filter.Skip),
			zap.Int("limit", filter.Limit), zap.Error(err),
		)
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, bl)
}
func getBlockTransactions(w http.ResponseWriter, r *http.Request) {
	numS := chi.URLParam(r, "num")
	bnum, err := strconv.ParseInt(numS, 10, 0)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'num' parameter %q: %s", numS, err))
		return
	}
	filter := new(models.PaginationFilter)
	if !parseGetParam(r, w, filter) {
		return
	}
	transactions := &models.TransactionList{}
	transactions.Transactions, err = backendInstance.GetBlockTransactionsByNumber(bnum, filter)
	if err != nil {
		logger.Error(
			"Failed to get latest blocks",
			zap.Int64("block", bnum),
			zap.Int("skip", filter.Skip),
			zap.Int("limit", filter.Limit),
			zap.Error(err),
		)
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "num")
	bnum, err := strconv.ParseInt(param, 10, 0)
	var block *models.Block
	if err != nil {
		block, err = backendInstance.GetBlockByHash(r.Context(), param)
	} else {
		block, err = backendInstance.GetBlockByNumber(r.Context(), bnum, false)
	}
	if err != nil {
		logger.Error("Failed to get block", zap.String("num", param), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	} else if block == nil {
		errorResponse(w, http.StatusNotFound, nil)
		return
	}
	writeJSON(w, http.StatusOK, block)
}

func checkBlockExist(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	block, err := backendInstance.GetBlockByHash(r.Context(), hash)
	if err != nil {
		logger.Error("Failed to get block", zap.String("hash", hash), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
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
	if !parseGetParam(r, w, filter) {
		return
	}
	addresses, err := backendInstance.GetContracts(filter)
	if err != nil {
		logger.Error("Failed to get contracts list", zap.Reflect("filter", filter), zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, addresses)
}

// var _ middleware.LogFormatter = &zapLogFormatter{}

// type zapLogFormatter struct {
// 	lgr *zap.Logger
// }

// NewLogEntry gathers information from the request, logs 'Request started'
// and returns a log entry which stores the info to log again later with
// the additional response info.
// func (z *zapLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
// 	h := NewHTTP(r, nil)
// 	lgr := z.lgr
// 	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
// 		lgr = lgr.With(zap.String("requestID", reqID))
// 	}
// 	lgr.Debug("Request started", zapdriver.HTTP(h))
// 	return &zapLogEntry{lgr: lgr, http: h}
// }

// var _ middleware.LogEntry = &zapLogEntry{}

// type zapLogEntry struct {
// 	lgr  *zap.Logger
// 	http *zapdriver.HTTPPayload
// }

// func (z *zapLogEntry) Write(status, bytes int, elapsed time.Duration) {
// 	z.http.Status = status
// 	z.http.ResponseSize = strconv.Itoa(bytes)
// 	z.http.Latency = fmt.Sprintf("%.9fs", elapsed.Seconds())
// 	z.lgr.Info("Request complete", zapdriver.HTTP(z.http))
// }

// func (z *zapLogEntry) Panic(v interface{}, stack []byte) {
// 	z.lgr = z.lgr.With(zap.String("stack", string(stack)), zap.String("panic", fmt.Sprintf("%+v", v)))
// }

// NewHTTP returns a new HTTPPayload struct, based on the passed
// in http.Request and http.Response objects. They are not modified
// in any way, unlike the zapdriver version this is based on.
func NewHTTP(req *http.Request, res *http.Response) *zapdriver.HTTPPayload {
	var p zapdriver.HTTPPayload
	if req != nil {
		p = zapdriver.HTTPPayload{
			RequestMethod: req.Method,
			UserAgent:     req.UserAgent(),
			RemoteIP:      req.RemoteAddr,
			Referer:       req.Referer(),
			Protocol:      req.Proto,
			RequestSize:   strconv.FormatInt(req.ContentLength, 10),
		}
		if req.URL != nil {
			p.RequestURL = req.URL.String()
		}
	}

	if res != nil {
		p.ResponseSize = strconv.FormatInt(res.ContentLength, 10)
		p.Status = res.StatusCode
	}

	return &p
}
