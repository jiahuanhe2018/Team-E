package blockchain

import (
        "bufio"
        "context"
        "crypto/rand"
        "fmt"
        "io"
        "log"
	"net"
        mrand "math/rand"

	//cid "github.com/ipfs/go-cid"
	//iaddr "github.com/ipfs/go-ipfs-addr"
        libp2p "github.com/libp2p/go-libp2p"
        crypto "github.com/libp2p/go-libp2p-crypto"
        host "github.com/libp2p/go-libp2p-host"
	//dht "github.com/libp2p/go-libp2p-kad-dht"
        pnet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
        ma "github.com/multiformats/go-multiaddr"
	//mh "github.com/multiformats/go-multihash"
)

// IPFS bootstrap nodes. Used to find other peers in the network.
//var bootstrapPeers = []string{
//	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
//	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
//	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
//	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
//	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
//}

type Node struct {
	Address	string	`json:"address"`
	MinerOnlyNode	bool	`json:"miner_only_node"`
	TrustedNodes	[]*Node	`json:"trusted_nodes"`
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "127.0.0.1"
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP.String()
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func MakeBasicHost(ip string, listenPort int, secio bool, randseed int64, initAccount string) (host.Host, error) {
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

	ip_addr := ip
	if ip == "" {
		ip_addr = GetOutboundIP()
	}
        opts := []libp2p.Option{
                libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", ip_addr, listenPort)),
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
			log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -a %s -d %s -p %s -secio\x1b[0m\" on a different terminal\n", listenPort+2, initAccount, fullAddr, BlockchainInstance.Proof)
		} else {
			log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -d %s -p %s -secio\x1b[0m\" on a different terminal\n", listenPort+2, fullAddr, BlockchainInstance.Proof)
		}
	} else {
		log.Printf("Now run \"go run main.go -c chain -l %d -d %s -p %s\" on a different terminal\n", listenPort+2, fullAddr, BlockchainInstance.Proof)
		if initAccount != "" {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -a %s -d %s -p %s\x1b[0m\" on a different terminal\n", listenPort+2, initAccount, fullAddr, BlockchainInstance.Proof)
                } else {
                        log.Printf("Now run \"go run main.go \x1b[32m -c chain -l %d -d %s -p %s\x1b[0m\" on a different terminal\n", listenPort+2, fullAddr, BlockchainInstance.Proof)
                }
	}

        return basicHost, nil
}

func HandleStream(s pnet.Stream) {

        log.Println("Got a new stream!")

        // Create a buffer stream for non blocking read and write.
        rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

        go ReadData(rw)
        go WriteData(rw)

        // stream 's' will stay open until you close it (or the other side closes it).
}

func MakeHostAndConnect(ip string, target string, listenPort int, secio bool, randseed int64, initAccount string) {
 	// Make a host that listens on the given multiaddress
        ha, err := MakeBasicHost(ip, listenPort, secio, randseed, initAccount)
        if err != nil {
                log.Fatal(err)
        }      

        if target == "" {  // bootstrap node
		log.Println("listening for connections")
                // Set a stream handler on host A. /p2p/1.0.0 is
                // a user-defined protocol name.
                ha.SetStreamHandler("/p2p/1.0.0", HandleStream)

                select {} // hang forever
                /**** This is where the listener code ends ****/
        } else {
                ha.SetStreamHandler("/p2p/1.0.0", HandleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go WriteData(rw)
		go ReadData(rw)

		select {} // hang forever

	}
}
