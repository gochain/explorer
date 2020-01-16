package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gochain-io/explorer/server/models"
)

const (
	MainnetURL = "https://explorer.gochain.io"
	TestnetURL = "https://testnet-explorer.gochain.io"
)

var (
	Mainnet = NewClient(MainnetURL)
	Testnet = NewClient(TestnetURL)
	wei, _  = new(big.Int).SetString("1000000000000000000", 10)
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url: url}
}

func (c *Client) Address(addr string) (*models.Address, error) {
	var data models.Address
	err := c.get("/api/address/"+addr, nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) AddressTransactions(addr string, txParams *TxParams) (*models.TransactionList, error) {
	var data models.TransactionList
	err := c.get(fmt.Sprintf("/api/address/%s/transactions", addr), txParams.sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) AddressHolders(addr string, sl *SkipLimit) (*models.TokenHolderList, error) {
	var data models.TokenHolderList
	err := c.get(fmt.Sprintf("/api/address/%s/holders", addr), sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) AddressInternalTransactions(addr string, sl *SkipLimit) (*models.TokenHolderList, error) {
	var data models.TokenHolderList
	err := c.get(fmt.Sprintf("/api/address/%s/internal_transactions", addr), sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) CirculatingSupply() (string, error) {
	return c.getStr("/circulatingSupply")
}

func (c *Client) CirculatingSupplyWei() (*big.Int, error) {
	s, err := c.CirculatingSupply()
	if err != nil {
		return nil, err
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("failed to set string: %q", s)
	}
	r = r.Mul(r, new(big.Rat).SetInt(wei))
	if !r.IsInt() {
		return nil, fmt.Errorf("not an int: %s", r)
	}
	return r.Num(), nil
}

func (c *Client) TotalSupply() (string, error) {
	return c.getStr("/totalSupply")
}

func (c *Client) TotalSupplyWei() (*big.Int, error) {
	s, err := c.TotalSupply()
	if err != nil {
		return nil, err
	}
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("failed to set string: %q", s)
	}
	r = r.Mul(r, new(big.Rat).SetInt(wei))
	if !r.IsInt() {
		return nil, fmt.Errorf("not an int: %s", r)
	}
	return r.Num(), nil
}

func (c *Client) RichList(sl *SkipLimit) (*models.Richlist, error) {
	var data models.Richlist
	err := c.get("/api/richlist", sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) Stats() (*models.Stats, error) {
	var data models.Stats
	err := c.get("/api/stats", nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) Blocks(sl *SkipLimit) (*models.BlockList, error) {
	var data models.BlockList
	err := c.get("/api/blocks", sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) Block(number uint64) (*models.Block, error) {
	var data models.Block
	err := c.get(fmt.Sprintf("/api/blocks/%d", number), nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) BlockTransactions(number uint64, sl *SkipLimit) (*models.TransactionList, error) {
	var data models.TransactionList
	err := c.get(fmt.Sprintf("/api/blocks/%d/transactions", number), sl.vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) Transaction(hash string) (*models.Transaction, error) {
	var data models.Transaction
	err := c.get(fmt.Sprintf("/api/transaction/%s", hash), nil, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) get(servicePath string, vals url.Values, data interface{}) error {
	resp, err := http.Get(c.url + fmt.Sprintf("%s?%s", servicePath, vals.Encode()))
	if err != nil {
		return err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) getStr(servicePath string) (string, error) {
	resp, err := http.Get(c.url + servicePath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

type SkipLimit struct {
	vals url.Values
}

func NewSkipLimit() *SkipLimit {
	return &SkipLimit{vals: make(url.Values)}
}

func (sl *SkipLimit) Skip(s int) *SkipLimit {
	sl.vals.Add("skip", strconv.Itoa(s))
	return sl
}

func (sl *SkipLimit) Limit(l int) *SkipLimit {
	sl.vals.Add("limit", strconv.Itoa(l))
	return sl
}

type TxParams struct {
	sl SkipLimit
}

func NewTxParams() *TxParams {
	tx := &TxParams{}
	tx.sl.vals = make(url.Values)
	return tx
}

func (tx *TxParams) Skip(s int) *TxParams {
	tx.sl.Skip(s)
	return tx
}

func (tx *TxParams) Limit(l int) *TxParams {
	tx.sl.Limit(l)
	return tx
}

func (tx *TxParams) FromTime(from time.Time) *TxParams {
	tx.sl.vals.Add("from_time", from.Format(time.RFC3339))
	return tx
}

func (tx *TxParams) ToTime(to time.Time) *TxParams {
	tx.sl.vals.Add("to_time", to.Format(time.RFC3339))
	return tx
}
