package zkclient

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pb "github.com/scalog/scalog/zookeeper/zookeeperpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Start() {
	zkNodesIp := viper.GetStringSlice("zk-servers")
	clients := []pb.ZooKeeperClient{}

	zkPort := int32(viper.GetInt("zk-port"))

	zids_ := viper.Get("zids").([]interface{})
	var zids []int

	for _, zid_ := range zids_ {
		zids = append(zids, zid_.(int))
	}

	rand.Seed(time.Now().UnixNano())

	for idx, ip := range zkNodesIp {
		fmt.Printf("initiating client to: %v:%v\n", ip, zkPort+int32(zids[idx]))
		serverAddress := fmt.Sprintf("%v:%v", ip, zkPort+int32(zids[idx]))

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewZooKeeperClient(conn)

		clients = append(clients, client)
	}

	for {
		client := clients[rand.Intn(len(clients))]
		fmt.Println("Choose an action:")
		fmt.Println("1. Create ZNode")
		fmt.Println("2. Read ZNode")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Enter ZNode Path: ")
			var path string
			fmt.Scanln(&path)
			fmt.Print("Enter ZNode Data: ")
			var data string
			fmt.Scanln(&data)

			znode := &pb.ZNode{
				Path: path,
				Data: []byte(data),
			}
			createResponse, err := client.CreateZNode(context.Background(), znode)
			if err != nil {
				log.Printf("CreateZNode failed: %v", err)
			} else {
				fmt.Printf("Created ZNode: %v\n", createResponse.Path)
			}

		case 2:
			fmt.Print("Enter ZNode Path: ")
			var path string
			fmt.Scanln(&path)

			path = strings.TrimSpace(path)
			if path == "" {
				log.Println("Invalid ZNode Path")
				continue
			}

			path = strings.Trim(path, "/") // Remove leading/trailing slashes

			pathObj := &pb.Path{Path: path}
			getResponse, err := client.GetZNode(context.Background(), pathObj)
			if err != nil {
				log.Printf("GetZNode failed: %v", err)
			} else {
				fmt.Printf("ZNode Path: %v\nZNode Data: %v\n", getResponse.Path, string(getResponse.Data))
			}

		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)

		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}
	}
}
