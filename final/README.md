# Final

This public chain is based on class homework 2 to homework 5. Modifications are on Genaro course code.

*   Transactions and special transactions were implemented in previous class homework 3-4.

*   The concensus is PoS + PoW, using PoS to choose small number of winners to do PoW.
I call it pox and use -p flag for concensus in the code.

*   The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, slash for bad behaviors (slash not implemented), and dynamically increase the number of PoW winners as number of miner peers increase significantly.

*   I added -i flag to specify a non 127.0.0.1 ip for the peer. If empty, the code will try to find the outbound ip for the peer. Currently it seems to find 192.168... ip instead of actual outbound ip.

*   Use --demo flag for manual input result in each round of mining and watching block mining process more slowly. Default is false, and does not require manul input result.    

Test commands that I used on different terminals are

```
// Create local wallet and accounts. Each peer should have its own wallet (replace "desiree" in createwallet command) and accounts.
// Run this command multiple times to create multiple accounts for one wallet.
// Later, -s and -a flags in the chain mode commands should use those wallet and accounts addresses.
go run main.go -c account desiree createwallet

// Bootstrap peer with a wallet (-s flag) and a miner account (-a flag).
go run main.go -c chain -s desiree -l 8080 -a 1FghRtifoTLuMsFRacRBpBYD2VwLmGoAhW -p pox -i ""

// Another peer directly dialing to the bootstrap peer (-d flag) with a different wallet (-s flag) and another miner account (-a flag).
// Each peer can have its own wallet and accounts. -s and -a flags should be stable after having wallet and accounts created.
// The address in -d flag could be a bit dynamic, but depends on the address print out by the bootstrap peer.
//
go run main.go  -c chain -s lzx -l 8082 -a 1Hn94smEVwEd3kPfvF39ozhqCQKGqce5qc -d /ip4/192.168.1.6/tcp/8080/ipfs/QmaHAkUhArtD2UwW4PMRGzZxCSVV6SAssy2C6XJRjonUWR -p pox

// Run this only after the peers running are connected and blockchain has >= 1-2 blocks.
// Use the accounts generated at the beginning.
curl -i --request POST --header 'Content-Type: application/json' --data '{"From":"1FghRtifoTLuMsFRacRBpBYD2VwLmGoAhW","To":"1BvT54va6zRhos2rVkT4DDSMexTCtT4q6J","Value":100,"Data":"message2"}' http://127.0.0.1:8081/txpool 
```

## Bugs to Fix and Features Wish List (in decending priority order)

*   Transaction currently only works best for accounts in the same wallet, but not for accounts in different wallets, probably due to account state not being properly updated. I'm (Desiree) working on fixing it, expected by 11/16 15pm Beijing time.

*   Unclear if actual outbound ip would work in p2p communication. I only tested on my single machine. This needs to be tested.

*   We may need a test script, given a list of all accounts and number of transactions, generate http post transactions and send them to the block chain.

*   This is better to be integrated with two layer optimization, such as TrueBit or Plasma. We can implement our own simplified version, e.g. in the code that handles http post transaction, and submit multiple transactions in batch to the blockchain.

*   States are not stored in leveldb yet, mainly in memory.

*   Occasionally dialing the bootstrap peer would give an error, but trying again usually works. I tried DHT peer discovery before, but it didn't work, but current peer communication is purely by dialing -- this limits p2p communication.
 
