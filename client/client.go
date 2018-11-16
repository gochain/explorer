package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

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

func (c *Client) AddressTransactions(addr string, skip, limit int) (*models.TransactionList, error) {
	vals := make(url.Values)
	vals.Add("skip", strconv.Itoa(skip))
	vals.Add("limit", strconv.Itoa(limit))
	var data models.TransactionList
	err := c.get(fmt.Sprintf("/api/address/%s/transactions", addr), vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) AddressHolders(addr string, skip, limit int) (*models.TokenHolderList, error) {
	vals := make(url.Values)
	vals.Add("skip", strconv.Itoa(skip))
	vals.Add("limit", strconv.Itoa(limit))
	var data models.TokenHolderList
	err := c.get(fmt.Sprintf("/api/address/%s/holders", addr), vals, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) AddressInternalTransactions(addr string, skip, limit int) (*models.TokenHolderList, error) {
	vals := make(url.Values)
	vals.Add("skip", strconv.Itoa(skip))
	vals.Add("limit", strconv.Itoa(limit))
	var data models.TokenHolderList
	err := c.get(fmt.Sprintf("/api/address/%s/internal_transactions", addr), vals, &data)
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

func (c *Client) RichList(skip, limit int) (*models.Richlist, error) {
	vals := make(url.Values)
	vals.Add("skip", strconv.Itoa(skip))
	vals.Add("limit", strconv.Itoa(limit))
	var data models.Richlist
	err := c.get("/api/richlist", vals, &data)
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
