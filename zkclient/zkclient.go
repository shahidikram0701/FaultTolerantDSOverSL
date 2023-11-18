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

func getClient(ip string, idx int) (pb.ZooKeeperClient, *grpc.ClientConn) {
	zkPort := int32(viper.GetInt("zk-port"))
	zids_ := viper.Get("zids").([]interface{})
	var zids []int

	for _, zid_ := range zids_ {
		zids = append(zids, zid_.(int))
	}

	serverAddress := fmt.Sprintf("%v:%v", ip, zkPort+int32(zids[idx]))

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	// defer conn.Close()

	client := pb.NewZooKeeperClient(conn)

	return client, conn

}

func Start() {
	zkNodesIp := viper.GetStringSlice("zk-servers")

	rand.Seed(time.Now().UnixNano())

	clientNum := rand.Intn(len(zkNodesIp))

	for {
		client, conn := getClient(zkNodesIp[clientNum], clientNum)
		defer conn.Close()
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
			for {
				createResponse, err := client.CreateZNode(context.Background(), znode)
				if err != nil {
					log.Printf("CreateZNode failed: %v", err)
					clientNum = (clientNum + 1) % len(zkNodesIp)
					client, conn = getClient(zkNodesIp[clientNum], clientNum)
					defer conn.Close()
				} else {
					fmt.Printf("Created ZNode: %v\n", createResponse.Path)
					break
				}
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

			// path = strings.Trim(path, "/") // Remove leading/trailing slashes

			pathObj := &pb.Path{Path: path}
			for {
				getResponse, err := client.GetZNode(context.Background(), pathObj)
				if err != nil {
					log.Printf("GetZNode failed: %v", err)
					clientNum = (clientNum + 1) % len(zkNodesIp)
					client, conn = getClient(zkNodesIp[clientNum], clientNum)
					defer conn.Close()
				} else {
					fmt.Printf("ZNode Path: %v\nZNode Data: %v\n", getResponse.Path, string(getResponse.Data))
					break
				}
			}

		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)

		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}
	}
}
