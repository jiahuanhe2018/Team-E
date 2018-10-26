Homework 5

1.  Blockchain governance

*   Because my homework 4 idea was to combine PoS and PoW for concensus (see below), my proposed governance needs to be
based on PoS and PoW as well, and I will use on chain governance.

*   Miner nodes will receive rewards for mining new blocks and packaging transactions, and receive penalties in the same way as bitcoin and ethererum.

*   To avoid rewards always go to the miners with largest stakes and largest computing powers, I propose to add time and number of blocks mined decay to the stakes held by miners. If a miner always gets new blocks, her stake will decay faster than the miner who gets less nor none blocks, so she has less chance to be one of the winning validators who get to do PoW.

*   Still I understand that if an attacker group can schedule attackers to appear in different time, one by one, the attacker might still be able to control this chain. 

*   To mitigate such attack, I could consider adding credit, so credit + stake + computing power decides who generates new blocks. I don't yet know how to assign credit though, because I think credit has something to do with off chain governance.

2.  Completeing account system

*   My implementation is modification based on code of previous several homeworks. Still use -p pox flag, where pox means PoS + PoWcombined concensus idea in homework 4.
     
*   I modified wallets, blockchain, rpc and main code. I let each node's wallets have a DefaultMinerAccount to receive miner's reward, and I calculated and assigned miner's reward.

*   You can check peer1_screenshot.txt and peer2_screenshot.txt for what it looks like when I run two peers and two accounts. I also http posted a real transaction between the two accounts.
-------------------------------------------------------------------------------------------------------------------------------

Homework 4 idea:

My idea is to combine PoS and PoW.

*   PoS is more efficient, but does not have reliable time, and easier to have fork chains.

*   PoW has more reliable time compared to PoS, but it wastes a lot of resources by letting all competing peers to calculate a difficult hash.

*   So I want to use PoS to choose multiple winner validators, and then run PoW only among these validators. Whoever calculates the difficult hash first produces the next block. In this way, I try to use the reliable time in PoW to help mitigate the time and fork issue of PoS, and make it more efficient by letting small number of peers to do PoW.

*   My other consideration is to use well validated way of proof. It is difficult and less sound in theory to think of a proof myself in a short time.

Comments are welcome directly in this repository. Thanks!
