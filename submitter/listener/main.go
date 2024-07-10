package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	pb "submitter/protos/model_server" // 更新为实际生成的protobuf包路径

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL       = "https://eth-sepolia.g.alchemy.com/v2/tMSj2Flz7SAnCGwTpvD_hPK2L9zFVbHn"
	contractAddress = "YOUR_CONTRACT_ADDRESS"
	contractABI     = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
)

func listenModelEvent() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	contractAddressHex := common.HexToAddress(contractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddressHex},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error: %v", err)
		case vLog := <-logs:
			fmt.Println("Log:", vLog)

			var transferEvent struct {
				From  common.Address
				To    common.Address
				Value *big.Int
			}

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log: %v", err)
			}

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Value: %s\n", transferEvent.Value.String())
		}
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewModelServerClient(conn)

	request := &pb.ModelRequest{
		RequestId: uuid.New().String(),
		ModelId:   1,
		Input:     "sample input",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.ModelCall(ctx, request)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Printf("Response: %v", response)
}
