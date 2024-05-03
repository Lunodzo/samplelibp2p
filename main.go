package main

import (
	//"bufio"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// A function for a simple chat stream application
func handleChatStream(stream network.Stream) {
    // Create a buffer stream for non blocking read and write.
    buf := make([]byte, 1024)
    for{
        // Read data from the stream
        n, err := stream.Read(buf)
        if err != nil {
            if err != io.EOF{
                fmt.Println("Error reading from stream", err)
            }
            return
        }
        message := string(buf[:n])
        fmt.Println("Received message:", message)
    }
}

// A function to send a message
func sendMessage(node host.Host, peer peer.ID, message string) {
    stream, err := node.NewStream(context.Background(), peer, ping.ID)
    if err != nil {
        fmt.Println("Error opening stream", err)
        return
    }
    _, err = stream.Write([]byte(message))
    if err != nil {
        fmt.Println("Error sending message", err)
    }
    stream.Close()
}

func main() {
	// start a libp2p node that listens on a random local TCP port,
    // but without running the built-in ping protocol
    node, err := libp2p.New(
        libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
        libp2p.Ping(false),
    )
    if err != nil {
        panic(err)
    }

    // configure our own ping protocol
    //pingService := &ping.PingService{Host: node}
    node.SetStreamHandler(ping.ID, handleChatStream)

    // print the node's PeerInfo in multiaddr format
    peerInfo := peerstore.AddrInfo{
        ID:    node.ID(),
        Addrs: node.Addrs(),
    }
    addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
    if err != nil {
        panic(err)
    }
    fmt.Println("libp2p node address:", addrs[0])

    // if a remote peer has been passed on the command line, connect to it
    // and send it 5 ping messages, otherwise wait for a signal to stop
    if len(os.Args) > 1 {
        addr, err := multiaddr.NewMultiaddr(os.Args[1])
        if err != nil {
            panic(err)
        }
        peer, err := peerstore.AddrInfoFromP2pAddr(addr)
        if err != nil {
            panic(err)
        }
        if err := node.Connect(context.Background(), *peer); err != nil {
            panic(err)
        }

        // Print the peer we are connected to and the peer connecting from
        fmt.Println("Connected to", peer.ID, "from", node.ID())


       /*  fmt.Println("sending 5 ping messages to", addr)
        ch := pingService.Ping(context.Background(), peer.ID)
        for i := 0; i < 5; i++ {
            res := <-ch
            fmt.Println("pinged", addr, "in", res.RTT)
        } */

        // Other operations can be done here
        // For example, to send a message to the connected peer
        // Loop to allow sending multiple messages
        for{
            fmt.Print("Enter message to send: ")
            reader := bufio.NewReader(os.Stdin)
            message, _ := reader.ReadString('\n')
            sendMessage(node, peer.ID, message)
        }
        //sendMessage(node, peer.ID, "Hello Peer!")
    } else {
        // wait for a SIGINT or SIGTERM signal
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        <-ch
        fmt.Println("Received signal, shutting down...")
    }

    // shut the node down
    if err := node.Close(); err != nil {
        panic(err)
    }
}
