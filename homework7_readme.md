# Homework 7

## Background

Below summarizes the previous choices I've made for my blockchain.
*   Concensus: PoS + PoW, and using PoS to choose small number of winners to do PoW.

*   Governance: Typical governance methods used by PoS, PoW system, with reward for mining new blocks and packaging transactions, and slash for bad behaviors.

*   Scaling: Two layers solver challenger protocol, similar to those of TrueBit and Plasma.

## Optimization Design

I referenced the idea of Loom Network + Plasma. I choose side chain as the way of optimization, and use two layers solver challenger protocol for committing data from side chain to the main (root) chain.

*   Side chain can reuse the same Block data structure as the main chain. Side chain can just use PoS or DPoS as concensus algorithm. A peer can only participate into either side chain or main chain, not both at the same time.

*   One side chain handles certain types or certain domains of transactions. Routing to an appropriate side chain or main chain can be handled in the http service of peers and underlying p2p communication. 

*   Side chain verifies and records those transactions, and proposes soft modifications on involved accounts. "soft" means those may be rejected when commiting to the main chain. If that happens, it will cause cascading rollback of subsequent related transactions. This can be implemented as traditional db transactions with reject as db transaction abort.

*   Periodically, or determined by number of transactions in a batch, or complexities of transactions in a batch, side chain enters solver challenger protocol to commit changes (transaction history, accounts updates) to the main chain. This is a typical task of miner peer on side chain, meaning miner peers either mine a new block, or tries to be a solver to commit changes to the main chain. Both tasks, if successful, should get reward. Commit changes will be new code change in existing blockchain code. 

*   When commiting changes to the main chain, merge will be needed if there are multiple side chains, and if main chain is also directly handling some transactions. To make it easier to merge account states, we can even have different accounts per chain, e.g. child accounts on side chain, and main accounts on the main chain, then merge does not need to all happen synchronously. Merging changes on the main chain can be implemented as one or multiple special transactions on the main chain.
