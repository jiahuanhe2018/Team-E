# Midterm

This public chain is based on my homework 2 to homework 5. Modifications are on Genaro course code.

*   Transactions and special transactions were implemented in homework 3-4.

*   My concensus in previous homeworks is PoS + PoW, using PoS to choose small number of winners to do PoW.
I call it pox and use -p flag for concensus in the code.

*   The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, slash for bad behaviors, and dynamically increase the number of PoW winners as number of miner peers increase significantly.

*   I added -i flag to specify a non 127.0.0.1 ip for the peer. If empty, the code will try to find the outbound ip for the peer. Currently it seems to find 192.168... ip instead of actual outbound ip. 

Test commands that I used on different terminals are

```
// Bootstrap peer with a wallet and an account.
go run main.go -c chain -s lzhx_ -l 8080 -a 15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi -p pox -i ""

// Another peer directly dialing to bootstrap peer with a different wallet and another account.
go run main.go  -c chain -s lzx -l 8082 -a 1Hn94smEVwEd3kPfvF39ozhqCQKGqce5qc -d /ip4/192.168.1.6/tcp/8080/ipfs/QmaHAkUhArtD2UwW4PMRGzZxCSVV6SAssy2C6XJRjonUWR -p pox

// Run this only after the peers running are connected and blockchain has >= 1-2 blocks.
curl -i --request POST --header 'Content-Type: application/json' --data '{"To":"1GBF4VK4Ys2R5Yiz3K6ZovXpYarGH2qsEc","From":"15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi","Value":100,"Data":"message2"}' http://127.0.0.1:8081/txpool 
```
