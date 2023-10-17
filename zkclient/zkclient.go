package zkclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/scalog/scalog/zookeeper/zookeeperpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Start() {
	zkPort := uint16(viper.GetInt("zk-port"))
	zkIpAddr := string(viper.GetString("zk-ip-address"))
	serverAddress := fmt.Sprintf("%v:%v", zkIpAddr, zkPort)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewZooKeeperClient(conn)

	for {
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
