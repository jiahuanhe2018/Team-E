# Homework 6

My concensus in previous homeworks is PoS + PoW, using PoS to choose small number of winners to do PoW.
The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, and slash for bad behaviors.

So I use two layers for optimization and scaling. The idea is very similar to TrueBit and Plasma.

*   Root chain: 

    *   Contains finalized blocks.
   
    *   Still runs PoS + PoW as before.

    *   Serves queries on account balances.

*   Child chain:

    *   Contains to be confirmed transactions and to be confirmed accounts updates.

    *   May still use PoS in a small group to find the solver.

    *   Serves new trasaction processing. 

*   Each transaction posted to a child chain should have a deposit. This transaction deposit is used as reward for solver or challenger later.

*   Solver packaging transactions and calculating accounts updates, or calculating smart contract also puts a deposit. Solver deposit < transaction deposit. Solver finishes one or multiple transactions, and is ready for being challenged and submitting to root chain. Challenging has a timeout, during which if no one challenges, solver submits transactions and related accounts updates successfully to root chain.

*   Challanger also needs to put a deposit every time it participates. Challenger deposit < transaction deposit. 

    *   If challenger verifies and decides not to challenge, its deposit goes back, and it gets a small percent of transaction deposit, while the majority of transaction deposit still goes to the solver. 

    *   If challenger challenges but is wrong, its deposit along with transaction deposit goes to the solver.

    *   There could be several challengers who do not agree with each other. If majority challengers agree with the solver, rest of challengers are penalized, and do not go to dispute. Otherwise, go to dispute. 

*   If a transaction needs to go to dispute, it will be posted to the root chain as a new transaction. The dispute result will notify the child chain when this transaction is packaged in a new block.

*   A peer can choose to participate in root chain or child chain, but it can only be one chain at a time. For example, a peer cannot be the miner on the root chain, and the solver or challenger on the child chain at the same time, to avoid conflict of interest.    
