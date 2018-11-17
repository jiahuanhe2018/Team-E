package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/GenaroSanada/Course/wallet"
)

func main() {
	wallets_str := flag.String("wallets", "", "wallet suffixes, separated by comma")
	num_txns := flag.Int("n", 1, "number of http post transactions to generate")
	post_addr := flag.String("address", "http://127.0.0.1:8081/txpool", "http post address, must correspond to a blockchain peer")
	flag.Parse()

	wallets_array := strings.Split(*wallets_str, ",")
	accounts := make([]string, 0)
	
	for _, wallet_str := range wallets_array {
		wallets, err := wallet.NewWallets(wallet_str)
		if err != nil {
			log.Panic(err)
		}
		addresses := wallets.GetAddresses()
		for _, address := range addresses {
			log.Printf("Account %s\n", address)
			accounts = append(accounts, address)
		}
	}
	
	f, err := os.Create("test_http_post_txns")
	if err != nil {
        	panic(err)
    	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < *num_txns; i++ {
		account1 := rand.Intn(len(accounts))
		account2 := rand.Intn(len(accounts))
		for account2 == account1 {
			account2 = rand.Intn(len(accounts))
		}
		w.WriteString(fmt.Sprintf("curl -i --request POST --header 'Content-Type: application/json' --data '{\"From\":\"%s\",\"To\":\"%s\",\"Value\":100,\"Data\":\"message\"}' %s\n", accounts[account1], accounts[account2], *post_addr))
		w.Flush()
	}	
}
