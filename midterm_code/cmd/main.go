package main

import (
	"time"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GenaroSanada/Course/blockchain"
	"github.com/GenaroSanada/Course/rpc"
	"github.com/GenaroSanada/Course/wallet"

	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
)

func main() {

	// Parse options from the command line
	command  := flag.String("c", "", "mode[ \"chain\" or \"account\"]")
	datadir := flag.String("datadir", "", "Data directory for the databases")
	listenF := flag.Int("l", 0, "wait for incoming connections[chain mode param]")
	target := flag.String("d", "", "target peer to dial[chain mode param]")
	suffix := flag.String("s", "", "wallet suffix [chain mode param]")
	initAccounts := flag.String("a", "", "init exist accounts whit value 10000")
	proof := flag.String("p", "pow", "\"pow\" or \"pos\" or \"pox\" [chain mode param]")
	minerOnly := flag.Bool("miner_only", false, "miner only write node or full node mining and serving http requests)")
	rendezvous := flag.String("sywangdesiree", "test", "Unique string to identify group of nodes during peer discovery.")
	secio := flag.Bool("secio", false, "enable secio[chain mode param]")
	seed := flag.Int64("seed", 0, "set random seed for id generation[chain mode param]")
	flag.Parse()

	if !(*proof == "pow" || *proof == "pos" || *proof == "pox") {
                flag.Usage()
		return
        }	

	if *command == "chain" {
		runblockchain(listenF, target, seed, secio, suffix, initAccounts, datadir, proof, minerOnly, rendezvous)
	}else if *command == "account" {
		cli := wallet.WalletCli{}
		cli.Run()
	}else {
		flag.Usage()
	}
}

func runblockchain(listenF *int, target *string, seed *int64, secio *bool, suffix *string, initAccounts *string, datadir *string, proof *string, minerOnly *bool, rendezvous *string) {
	if *datadir == ""{
		log.Println("data directory for this node miss，The data of the node will not be stored.")
	} else {
		if IsFile(*datadir) {
			log.Println(fmt.Sprintf("datadir[%s] is a file", *datadir))
			return
		}

		if !IsExist(*datadir) {
			log.Println(fmt.Sprintf("datadir[%s] not exist", *datadir))
			return
		}
	}

	t := time.Now()
	genesisBlock := blockchain.Block{}
	defaultAccounts := make(map[string]blockchain.Account)

	if *initAccounts != ""{
		if wallet.ValidateAddress(*initAccounts) == false {
			fmt.Println("Invalid address")
			return
		}
		newAccount := new(blockchain.Account)
		newAccount.Balance = 10000
		newAccount.State = 0
		defaultAccounts[*initAccounts] = *newAccount
	}
	genesisBlock = blockchain.Block{0, t.String(), 0, blockchain.CalculateBlockHash(genesisBlock), "", 100, blockchain.Difficulty, "", "", make(map[string]blockchain.Transaction), defaultAccounts, 0}

	var blocks []blockchain.Block
	blocks = append(blocks, genesisBlock)
	blockchain.BlockchainInstance.Blocks =  blocks
	blockchain.BlockchainInstance.DataDir = *datadir
        blockchain.BlockchainInstance.Proof = *proof
	blockchain.BlockchainInstance.NumPowWinners = 1

	blockchain.BlockchainInstance.ReadDataFromFile()

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	if *listenF == 0 {
		log.Fatal("Please provide a peer port to bind on with -l")
	}

	if *suffix == "" {
		log.Println("option param -s miss [you can't send transacion with this node]")
	} else {
		wallets, err := wallet.NewWallets(*suffix)
		if err != nil {
			log.Panic(err)
		}
		blockchain.Wallets.Wallets = wallets.Wallets
		if *initAccounts != "" {
			blockchain.Wallets.DefaultMinerAccount = *initAccounts
			wall := wallets.GetWallet(*initAccounts)
			wall.Balance = 10000
			wall.State = 0
		}
	}
	
	if !(*minerOnly) {
		go rpc.RunHttpServer(*listenF+1)
	}

	// Make a host that listens on the given multiaddress
	blockchain.MakeHostAndConnect(*target, *listenF, *secio, *seed, *initAccounts, *rendezvous)
}

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func IsExist(dir string) bool {
	fi, e := os.Stat(dir)
	if e != nil {
		return os.IsExist(e)
	}
	return fi.IsDir()
}

