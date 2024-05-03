Simpe P2P Application

It allows two peers to communicate. In this case one peer is listening and the other sending messages. This application uses libp2p library in GO, specifically the PING module in libp2p. The libp2p ping protocol is a simple liveness check that peers can use to test the connectivity and performance between two peers. The libp2p ping protocol is different from the ping command line utility (ICMP ping), as it requires an already established libp2p connection. A peer opens a new stream on an existing libp2p connection and sends a ping request with a random 32 byte payload. The receiver echoes these 32 bytes back on the same stream. By measuring the time between the request and response, the initiator can calculate the round-trip time of the underlying libp2p connection. The stream can be reused for future pings from the initiator.

The ping protocol ID is /ipfs/ping/1.0.0. It can also be read in the imported GO modules



This application can be extended to support actual conversations. 
