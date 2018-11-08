# Midterm

This public chain is based on my homework 2 to homework 5. Modifications are on Genaro course code.

*   Transactions and special transactions were implemented in homework 3-4.

*   My concensus in previous homeworks is PoS + PoW, using PoS to choose small number of winners to do PoW.
I call it pox in the code.

*   The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, slash for bad behaviors, and dynamically increase the number of PoW winners as number of miner peers increase significantly.

*   I added dynamic peer discovery using some known bootstrap peers. It worked, but it took several minutes to hook up just the test peers running on my machine. You will see mixed blockchain messages and many libp2p logs.

Test commands that I used on different terminals are
go run main.go -c chain -s lzhx_ -l 8080 -a 15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi -p pox

go run main.go  -c chain -s lzx -l 8082 -a 1Hn94smEVwEd3kPfvF39ozhqCQKGqce5qc -p pox

// Run this only after the peers running on local machine are connected and blockchain has >= 1-2 blocks.
curl -i --request POST --header 'Content-Type: application/json' --data '{"From":"1Hn94smEVwEd3kPvF39ozhqCQKGqce5qc","To":"15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi","Value":100,"Data":"message2"}' http://127.0.0.1:8083/txpool 
