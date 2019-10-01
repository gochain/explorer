package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/gochain-io/gochain/v3/common/hexutil"
	"go.uber.org/zap"
)

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
	url string
	Lgr *zap.Logger
}

func NewEthClient(url string, lgr *zap.Logger) *EthRPC {
	rpc := &EthRPC{
		url: url,
		Lgr: lgr,
	}

	return rpc
}

func (rpc *EthRPC) Call(ctx context.Context, method string, params ...interface{}) (json.RawMessage, error) {
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

	req, err := http.NewRequest("POST", rpc.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req.WithContext(ctx))
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
func (rpc *EthRPC) call(ctx context.Context, method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(ctx, method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}
func (rpc *EthRPC) ethGetBalance(ctx context.Context, address, block string) (*big.Int, error) {
	lgr := rpc.Lgr.With(zap.String("address", address))
	lgr.Debug("Checking balance")
	var response hexutil.Big
	if err := rpc.call(ctx, "eth_getBalance", &response, address, block); err != nil {
		return nil, err
	}
	lgr.Debug("Got balance", zap.Stringer("balance", &response))
	return response.ToInt(), nil
}

// func (rpc *EthRPC) ethGetBlockByNumber(number int64, withTransactions bool) (types.Block, error) {
// 	return rpc.getBlock("eth_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
// }

func (rpc *EthRPC) ethBlockNumber(ctx context.Context) (int64, error) {
	var response hexutil.Big
	if err := rpc.call(ctx, "eth_blockNumber", &response); err != nil {
		return 0, err
	}
	return response.ToInt().Int64(), nil
}

func (rpc *EthRPC) codeAt(ctx context.Context, address, block string) ([]byte, error) {
	var result hexutil.Bytes
	err := rpc.call(ctx, "eth_getCode", &result, address, block)
	return result, err
}

func (rpc *EthRPC) ethTotalSupply(ctx context.Context) (*big.Int, error) {
	var response hexutil.Big
	if err := rpc.call(ctx, "eth_totalSupply", &response, "latest"); err != nil {
		return nil, fmt.Errorf("failed to get total supply: %v", err)
	}
	return response.ToInt(), nil
}
