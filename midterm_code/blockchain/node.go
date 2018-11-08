package blockchain

import (
        "bufio"
        "context"
        "crypto/rand"
        "fmt"
        "io"
        "log"
	"strings"
	"time"
        mrand "math/rand"

	cid "github.com/ipfs/go-cid"
	iaddr "github.com/ipfs/go-ipfs-addr"
        libp2p "github.com/libp2p/go-libp2p"
        crypto "github.com/libp2p/go-libp2p-crypto"
        host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
        net "github.com/libp2p/go-libp2p-net"
	//peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
        ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
)

// IPFS bootstrap nodes. Used to find other peers in the network.
var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
}

type Node struct {
	Address	string	`json:"address"`
	MinerOnlyNode	bool	`json:"miner_only_node"`
	TrustedNodes	[]*Node	`json:"trusted_nodes"`
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func MakeBasicHost(listenPort int, secio bool, randseed int64, initAccount string) (host.Host, error) {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
        // deterministic randomness source to make generated keys stay the same
        // across multiple runs
        var r io.Reader
        if randseed == 0 {
                r = rand.Reader
        } else {
                r = mrand.New(mrand.NewSource(randseed))
        }

        // Generate a key pair for this host. We will use it
        // to obtain a valid host ID.
        priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
        if err != nil {
                return nil, err
        }

        opts := []libp2p.Option{
                libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
                libp2p.Identity(priv),
        }
        basicHost, err := libp2p.New(context.Background(), opts...)
        if err != nil {
                return nil, err
        }

        // Build host multiaddress
        hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

        // Now we can build a full multiaddress to reach this host
        // by encapsulating both addresses:
        addr := basicHost.Addrs()[0]
        fullAddr := addr.Encapsulate(hostAddr)
        log.Printf("I am %s\n", fullAddr)
        if secio {
                if initAccount != "" {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -a %s -p %s -secio\x1b[0m\" on a different terminal\n", listenPort+2, initAccount, BlockchainInstance.Proof)
                } else {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -p %s -secio\x1b[0m\" on a different terminal\n", listenPort+2, BlockchainInstance.Proof)
                }
        } else {
		 log.Printf("Now run \"go run main.go -c chain -l %d -p %s\" on a different terminal\n", listenPort+2, BlockchainInstance.Proof)
                if initAccount != "" {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -a %s -p %s\x1b[0m\" on a different terminal\n", listenPort+2, initAccount, BlockchainInstance.Proof)
                } else {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -p %s\x1b[0m\" on a different terminal\n", listenPort+2, BlockchainInstance.Proof)
                }
        }

        return basicHost, nil
}

func HandleStream(s net.Stream) {

        log.Println("Got a new stream!")

        // Create a buffer stream for non blocking read and write.
        rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

        go ReadData(rw)
        go WriteData(rw)

        // stream 's' will stay open until you close it (or the other side closes it).
}

func MakeHostAndConnect(target string, listenPort int, secio bool, randseed int64, initAccount string, rendezvous string) {
 	// Make a host that listens on the given multiaddress
        ha, err := MakeBasicHost(listenPort, secio, randseed, initAccount)
        if err != nil {
                log.Fatal(err)
        }      

        // if target == "" {  // bootstrap node
        //        log.Println("listening for connections")
                // Set a stream handler on host A. /p2p/1.0.0 is
                // a user-defined protocol name.
        //        ha.SetStreamHandler("/p2p/1.0.0", HandleStream)

        //        select {} // hang forever
                /**** This is where the listener code ends ****/
        //} else {
                ha.SetStreamHandler("/p2p/1.0.0", HandleStream)

		ctx := context.Background()
		kadDht, err := dht.New(ctx, ha)
                if err != nil {
                        panic(err)
                }

		log.Println("opening stream to bootstrap peer")
		// Let's connect to the bootstrap nodes first. They will tell us about the other nodes in the network.
		for _, peerAddr := range bootstrapPeers {
			addr, _ := iaddr.ParseString(peerAddr)
			peerinfo, _ := pstore.InfoFromP2pAddr(addr.Multiaddr())

			if err := ha.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Connection established with bootstrap node: ", *peerinfo)
			}
		}

		v1b := cid.V1Builder{Codec: cid.Raw, MhType: mh.SHA2_256}
		rendezvousPoint, _ := v1b.Sum([]byte(rendezvous))

		fmt.Println("announcing ourselves...")
		tctx, cancel := context.WithTimeout(ctx, time.Second*10000)
		defer cancel()
		if err := kadDht.Provide(tctx, rendezvousPoint, true); err != nil {
			log.Println("KadDht.Provide error")
		}

		// Now, look for others who have announced
		// This is like your friend telling you the location to meet you.
		// 'FindProviders' will return 'PeerInfo' of all the peers which
		// have 'Provide' or announced themselves previously.
		fmt.Println("searching for other peers...")
		tctx, cancel = context.WithTimeout(ctx, time.Second*10000)
		defer cancel()
		peers, err := kadDht.FindProviders(tctx, rendezvousPoint)
		if err != nil {
			log.Println("kdaDht.FindProviders error")
		}
		fmt.Printf("Found %d peers!\n", len(peers))

                log.Println("opening stream")
               	for _, p := range peers {
			if p.ID == ha.ID() || len(p.Addrs) == 0 || !strings.Contains(p.Addrs[0].String(), "127.0.0.1") {
				// No sense connecting to ourselves, or if addrs are not available, or to public peers.
				continue
			}
		 	// make a new stream from host B to host A
			// it should be handled on host A by the handler we set above because
                	// we use the same /p2p/1.0.0 protocol
                	s, err := ha.NewStream(ctx, p.ID, "/p2p/1.0.0")
                	if err != nil {
                        	log.Fatalln(err)
                	} else {       
                		// Create a buffered stream so that read and writes are non blocking.
                		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
                
                		// Create a thread to read and write data.
                		go WriteData(rw)
                		go ReadData(rw)
			}
		}
                
                select {} // hang forever
	//}
}
