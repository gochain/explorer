package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/gochain-io/gochain/v3/common/hexutil"
	"go.uber.org/zap"
)

type httpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EthError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type GenesisAccount struct {
	Balance string `json:"balance"`
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type EthRPC struct {
	url    string
	client httpClient
	Lgr    *zap.Logger
}

func NewEthClient(url string, lgr *zap.Logger) *EthRPC {
	rpc := &EthRPC{
		url:    url,
		client: http.DefaultClient,
		Lgr:    lgr,
	}

	return rpc
}

func (rpc *EthRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}
func (rpc *EthRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}
func (rpc *EthRPC) ethGetBalance(address, block string) (*big.Int, error) {
	var response string
	rpc.Lgr.Debug("response from eth_getBalance", zap.String("checking balance", address))
	if err := rpc.call("eth_getBalance", &response, address, block); err != nil {
		return new(big.Int), err
	}
	rpc.Lgr.Debug("response from eth_getBalance", zap.String("checking balance response", response))
	balance, err := parseBigInt(response)
	return balance, err
}

// func (rpc *EthRPC) ethGetBlockByNumber(number int64, withTransactions bool) (types.Block, error) {
// 	return rpc.getBlock("eth_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
// }

func (rpc *EthRPC) ethBlockNumber() (int64, error) {
	var response string
	if err := rpc.call("eth_blockNumber", &response); err != nil {
		return 0, err
	}
	return parseInt(response)
}

func (rpc *EthRPC) codeAt(address, block string) ([]byte, error) {
	var result hexutil.Bytes
	err := rpc.call("eth_getCode", &result, address, block)
	return result, err
}

func (rpc *EthRPC) ethTotalSupply() (*big.Int, error) {
	var response string
	if err := rpc.call("eth_totalSupply", &response, "latest"); err != nil {
		return new(big.Int), err
	}
	totalSupply, _ := parseBigInt(response)
	rpc.Lgr.Info("response from EthTotalSupply", zap.String("totalSupply", totalSupply.String()))
	return totalSupply, nil
}

func parseBigInt(value string) (*big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(value, &i)

	return &i, err
}

func parseInt(value string) (int64, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(value, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}
