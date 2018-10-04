In cmd/main.go, I added a new flag --proof, -p for pow (default) or pos.
In blockchain/blockchain.go, I put the functionality of PoW and PoS in the Blockchain, Block structures and related functions.

You can still use example command in the original Course instruction to run, 
e.g. To start a peer to run PoW,
go run main.go -c chain -s lzhx_ -l 8080 -a 15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi

To start a peer to run PoS,
go run main.go -c chain -s lzhx_ -l 8080 -a 15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi -p pos

Check my pow_running_screenshot.txt and pos_running_screenshot.txt what it looked like when I ran 2 peers.
