package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	databases   = []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"}
	collections = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
)

func main() {

	// Connect to MongoDB
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:W0;k;8H*ILSDXla8@127.0.0.1:27021/?directConnection=true"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:0n*yJ4vV7ZK(_KL5@127.0.0.1:27018/?directConnection=true"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	rand.Seed(time.Now().UnixNano())

	// Insert 1 GB of random data
	dataSize := 1024 * 1024 * 1024 * 2 // 1 GB
	batchSize := 1000                  // Adjust based on your system's capacity
	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < dataSize/batchSize; i++ {
				batch := make([]interface{}, batchSize)

				for j := 0; j < batchSize; j++ {
					document := map[string]interface{}{
						"key":   rand.Int(),
						"value": rand.Float64(),
					}
					batch[j] = document
				}
				coll := getRandomCollection(client)
				_, err := coll.InsertMany(ctx, batch)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Inserted %d documents in %s/%s \n", batchSize, coll.Database().Name(), coll.Name())
			}
		}()
	}
	wg.Wait()

	fmt.Println("Data insertion complete.")
}

func getRandomCollection(client *mongo.Client) *mongo.Collection {
	// Generate random index for databases and collections slices
	databaseIndex := rand.Intn(len(databases))
	collectionIndex := rand.Intn(len(collections))

	db := client.Database(databases[databaseIndex])
	collection := db.Collection(collections[collectionIndex])
	return collection
}
