Homework 4 idea:

My idea is to combine PoS and PoW.

*   PoS is more efficient, but does not have reliable time, and easier to have fork chains.

*   PoW has more reliable time compared to PoS, but it wastes a lot of resources by letting all competing peers to calculate a difficult hash.

*   So I want to use PoS to choose multiple winner validators, and then run PoW only among these validators. Whoever calculates the difficult hash first produces the next block. In this way, I try to use the reliable time in PoW to help mitigate the time and fork issue of PoS, and make it more efficient by letting small number of peers to do PoW.

*   My other consideration is to use well validated way of proof. It is difficult and less sound in theory to think of a proof myself in a short time.


My implementation is modification on the code of homework 3. I added -p pox option, where pox means the above PoS + PoW idea. Also check the 3 peers screenshots I captured when running my code.

Comments are welcome directly in this repository. Thanks!
