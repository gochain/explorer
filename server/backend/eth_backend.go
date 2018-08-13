package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/gochain-io/gochain/common"
	"github.com/rs/zerolog/log"
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
}

func NewEthClient(url string) *EthRPC {
	rpc := &EthRPC{
		url:    url,
		client: http.DefaultClient,
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
	log.Info().Str("checking balance", address).Msg("response from eth_getBalance")
	if err := rpc.call("eth_getBalance", &response, address, block); err != nil {
		return new(big.Int), err
	}
	log.Info().Str("checking balance response", response).Msg("response from eth_getBalance")
	balance, err := parseBigInt(response)
	return balance, err
}

func (rpc *EthRPC) ethTotalSupply() (*big.Int, error) {
	var response string
	if err := rpc.call("eth_totalSupply", &response, "latest"); err != nil {
		return new(big.Int), err
	}
	totalSupply, _ := parseBigInt(response)
	log.Info().Str("totalSupply", totalSupply.String()).Msg("response from EthTotalSupply")
	return totalSupply, nil
}

func (rpc *EthRPC) ethGenesisAlloc() (map[common.Address]GenesisAccount, error) {
	var response map[common.Address]GenesisAccount
	if err := rpc.call("eth_genesisAlloc", &response); err != nil {
		log.Info().Err(err).Msg("failed response from eth_genesisAlloc")
		return nil, err
	}
	log.Info().Interface("supply", response).Msg("response from eth_genesisAlloc")
	return response, nil
}

func (rpc *EthRPC) genesisAlloc() (*big.Int, error) {
	data, err := rpc.ethGenesisAlloc()
	genesisAlloc := new(big.Int)
	if err != nil {
		log.Info().Err(err).Msg("failed response from GenesisAlloc")
		return genesisAlloc, err
	}
	for _, val := range data {
		bal, _ := parseBigInt(val.Balance)
		genesisAlloc = new(big.Int).Add(genesisAlloc, bal)
	}
	log.Info().Str("GenesisAlloc", genesisAlloc.String()).Msg("response from GenesisAlloc")
	return genesisAlloc, nil
}

func (rpc *EthRPC) circulatingSupply() (*big.Int, error) {
	genesisAllocated, err := rpc.genesisAlloc()
	totalSupply, err2 := rpc.ethTotalSupply()
	if err != nil || err2 != nil {
		log.Info().Err(err).Err(err2).Msg("failed parsing CirculatingSupply")
		return new(big.Int), err
	}
	circulatingSupply := new(big.Int).Sub(totalSupply, genesisAllocated)
	return circulatingSupply, nil
}

func parseBigInt(value string) (*big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(value, &i)

	return &i, err
}
