1) 在第二次作业的基础上增加交易
I modified Transaction struct, kept them as kv map in blockchain, communicated and synchronized transactions among peers.

2）增加state以及查state的特殊交易. 特殊交易：以前的交易只能用来转账，现在通过交易可以改变accounts里state状态.
I added special payload in Transaction struct, and that special payload only has one simple operation: reset all the accounts' states.
Client can still use http post to post normal transaction and this special transaction to a peer.
See my *_screenshots.txt for what it looked like when I ran two peers and a http post client.

3）（optional）将EVM添加在第二次作业的基础上
Didn't do it. Don't think I would have enough time for it.
