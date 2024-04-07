package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kmodules.xyz/client-go/tools/portforward"
	"log"
	"strings"
)

func getPrimaryAndSecondary(tunnel *portforward.Tunnel) error {
	url := fmt.Sprintf("mongodb://root:%s@localhost:%v/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", password, tunnel.Local)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	adminDB := client.Database("admin")

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Run the rs.status() command
	var result bson.M
	err = adminDB.RunCommand(context.Background(), bson.D{{"replSetGetStatus", 1}}).Decode(&result)
	if err != nil {
		fmt.Println("Error running rs.status():", err)
		return err
	}

	members, ok := result["members"].(primitive.A)
	if !ok {
		fmt.Println("Error parsing members array")
		return err
	}

	primary := ""
	secondary := ""
	for _, member := range members {
		memberInfo, ok := member.(primitive.M)
		if !ok {
			fmt.Println("Error parsing member information")
			continue
		}

		// Check member state
		state := memberInfo["state"].(int32)
		if state == 1 {
			primary = memberInfo["name"].(string)
		} else if state == 2 {
			secondary = memberInfo["name"].(string)
		}
	}
	primaryPod = strings.Split(primary, ".")[0]
	secondaryPod = strings.Split(secondary, ".")[0]
	return nil
}

func connectToMongo(pTunnel, sTunnel *portforward.Tunnel) (*mongo.Client, *mongo.Client) {
	url := fmt.Sprintf("mongodb://root:%s@localhost:%v/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", password, pTunnel.Local)
	clientOptions := options.Client().ApplyURI(url)
	primaryClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	url = fmt.Sprintf("mongodb://root:%s@localhost:%v/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", password, sTunnel.Local)
	clientOptions = options.Client().ApplyURI(url)
	secondaryClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to primary %v on port %v, & secondary %v on port %v \n", primaryPod, pTunnel.Local, secondaryPod, sTunnel.Local)
	return primaryClient, secondaryClient
}

func TunnelToDBService(config *rest.Config) (*portforward.Tunnel, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	tunnel := portforward.NewTunnel(portforward.TunnelOptions{
		Client:    client.CoreV1().RESTClient(),
		Config:    config,
		Resource:  "services",
		Name:      mongodb.Name,
		Namespace: mongodb.Namespace,
		Remote:    27017,
	})

	return tunnel, tunnel.ForwardPort()
}

func TunnelToDBPod(config *rest.Config, podName string) (*portforward.Tunnel, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	tunnel := portforward.NewTunnel(portforward.TunnelOptions{
		Client:    client.CoreV1().RESTClient(),
		Config:    config,
		Resource:  "pods",
		Name:      podName,
		Namespace: mongodb.Namespace,
		Remote:    27017,
	})

	return tunnel, tunnel.ForwardPort()
}
