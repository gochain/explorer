// +build client

package client_test

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/gochain-io/explorer/client"
)

const (
	testnetTestAddr = "0x78cb510135787f42a23aD46998eB16B756111559"
	mainnetTestAddr = "0x3ad14430951aBa12068A8167CeBE3ddd57614432"
)

var (
	url  = flag.String("url", "", "url")
	test = flag.Bool("test", false, "testnet")
	addr = flag.String("addr", "", "address")

	testAddr string
	c        *client.Client
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *url != "" {
		log.Println("Using custom url:", *url)
		c = client.NewClient(*url)
	} else {
		if *test {
			log.Println("Using testnet")
			c = client.Testnet
		} else {
			log.Println("Using mainnet")
			c = client.Mainnet
		}
	}
	if *addr != "" {
		testAddr = *addr
	} else {
		if *test {
			testAddr = testnetTestAddr
		} else {
			testAddr = mainnetTestAddr
		}
	}

	os.Exit(m.Run())
}

func TestClient_TotalSupplyWei(t *testing.T) {
	if supply, err := c.TotalSupplyWei(); err != nil {
		t.Error("Failed to get total supply:", err)
	} else {
		t.Log("Total supply (wei):", supply)
	}
}

func TestClient_CirculatingSupplyWei(t *testing.T) {
	if supply, err := c.CirculatingSupplyWei(); err != nil {
		t.Error("Failed to get circulating supply:", err)
	} else {
		t.Log("Circulating supply (wei):", supply)
	}
}

func TestClient_RichList(t *testing.T) {
	if richlist, err := c.RichList(client.NewSkipLimit().Skip(5).Limit(10)); err != nil {
		t.Error("Failed to get rich list:", err)
	} else {
		b, err := json.Marshal(&richlist)
		if err != nil {
			t.Error("Failed to marshal rich list:", err)
		} else {
			t.Logf("RichList: %s\n", b)
		}
	}
}

func TestClient_Address(t *testing.T) {
	if address, err := c.Address(testAddr); err != nil {
		t.Error("Failed to get address")
	} else {
		b, err := json.Marshal(&address)
		if err != nil {
			t.Error("Failed to marshal address:", err)
		} else {
			t.Logf("Address: %s\n", b)
		}
	}
}

func TestClient_AddressTransactions(t *testing.T) {
	if txs, err := c.AddressTransactions(testAddr, client.NewTxParams().Skip(1).Limit(5)); err != nil {
		t.Error("Failed to get transactions:", err)
	} else {
		b, err := json.Marshal(&txs)
		if err != nil {
			t.Error("Failed to marshal transactions:", err)
		} else {
			t.Logf("Transactions: %s\n", b)
		}
	}
}

func TestClient_AddressHolders(t *testing.T) {
	if txs, err := c.AddressHolders(testAddr, client.NewSkipLimit().Skip(1).Limit(5)); err != nil {
		t.Error("Failed to get holders:", err)
	} else {
		b, err := json.Marshal(&txs)
		if err != nil {
			t.Error("Failed to marshal holders:", err)
		} else {
			t.Logf("Address Holders: %s\n", b)
		}
	}
}

func TestClient_AddressInternalTransactions(t *testing.T) {
	if txs, err := c.AddressInternalTransactions(testAddr, client.NewSkipLimit().Skip(1).Limit(5)); err != nil {
		t.Error("Failed to get holders:", err)
	} else {
		b, err := json.Marshal(&txs)
		if err != nil {
			t.Error("Failed to marshal holders:", err)
		} else {
			t.Logf("Address Holders: %s\n", b)
		}
	}
}
