# Bitcoin Handshake Task

This app is a solution for the test [requirements](https://github.com/eigerco/recruitment-exercises/blob/main/golang/bitcoin-handshake.md).

## Requirements
* OS: Linux
* Golang: v 1.22
* Bitcoin Node

## What does App do
This app connects to the bitcoin node by provided node host/port. 
After tcp connection is made, the app sends two specific messages to establish the handshake: version and verack message.
First, the app sends version message. If the message is accepted by the node, the app will receive a version message from the node.
After a version message, the app sends verack message and waits a verack message in a response from the node.
If the app gets version and verack messages, it means that the app made a handshake.

The app writes its results into console logs.
When the app gets all necessary messages from the node, it prints them to the console.
If everything went good, you'll see message in the log:
```shell
{"level":"info","msg":"All necessary messages for connection are received."}
```

## How to run

### Run Bitcoin Node
You can use either public Bitcoin node in a testnet or run you own node locally.
To run bitcoin node locally, please follow this official instructions:
* official [code repo](https://github.com/bitcoin/bitcoin) 
* instruction for [Unix system](https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md)
* if you wish to run the node on another OS system, please find the appropriate instruction [here](https://github.com/bitcoin/bitcoin/tree/master/doc) 

### Run Go App
To run Go App, you need to have:
* Go v1.22
* make

#### Run the app:

Download dependencies
```shell
    make dep
```

Start the app
```go
    go run main.go --node.host=<NODE_HOST> --node.port=<NODE_PORT>
```

```shell
NODE_HOST - address of bitcoin node
NODE_PORT - port of bitcoin node. If you connect to testnet, the port is 18333
```

Default values for incoming parameters:
```shell
  -node.host string
        Host of blockchain node (default "127.0.0.1")
  -node.port int
        Port of blockchain node (default 18333)
```

Example of the logs results:
```shell
go run main.go --node.host=127.0.0.1 --node.port=18333                                
{"level":"info","msg":"App starting...","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"Initializing all services...","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"Connecting to bitcoin node, host 127.0.0.1, port 18333...","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"starting reading incoming messages from node","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"starting handshake","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"Starting listening incoming messages from node...","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"got version message","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"{Header:{Magic:118034699 Command:version Length:102 Checksum:[151 53 207 2]} Payload:{Version:70016 Services:1037 Timestamp:1713282931 AddrRecv:{Timestamp:0 Services:1037 IP::: Port:0} AddrFrom:{Timestamp:0 Services:0 IP:31.50.193.200 Port:55752} Nonce:8996837035657133647 UserAgent: StartHeight:1632841488 Relay:true} Error:\u003cnil\u003e}\n","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"version message received successfully, trying to send verack message","time":"2024-04-16T16:55:31+01:00"}
{"level":"warning","msg":"err while parsing msg payload: unknown command, can't parse payload: wtxidrelay. Skipping","time":"2024-04-16T16:55:31+01:00"}
{"level":"warning","msg":"err while parsing msg payload: unknown command, can't parse payload: sendaddrv2. Skipping","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"got verack message","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"{Header:{Magic:118034699 Command:verack Length:0 Checksum:[93 246 224 226]} Payload:{} Error:\u003cnil\u003e}\n","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"verack message received successfully","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"All necessary messages for connection are received.","time":"2024-04-16T16:55:31+01:00"}
{"level":"info","msg":"Stopping the App...","time":"2024-04-16T16:55:31+01:00"}
{"level":"warning","msg":"stopping receiving thread by context done","time":"2024-04-16T16:55:31+01:00"}
```