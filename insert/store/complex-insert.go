package store

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
	databases   = []string{"aa", "bb", "cc"}
	collections = []string{"one", "two", "three", "four"}
)

func main() {
	fmt.Println("Starts")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:12345@127.0.0.1:27018/?directConnection=true"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()

	fmt.Println("Connected to MongoDB")

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	rand.Seed(time.Now().UnixNano())

	dataSize := 1024 * 1024 * 1024 * 4 // 1 GB
	batchSize := 100000
	var wg sync.WaitGroup

	ch := make(chan bool)
	go func(ch chan bool) {
		sz := calcTotalStorageSize()
		if sz >= dataSize {
			ch <- true
		}
	}(ch)

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := make([]interface{}, batchSize)
			for b := 0; b < batchSize; b++ {
				document := map[string]interface{}{
					"key":   rand.Int() % 50,
					"value": rand.Float64(),
				}
				batch[b] = document
			}
			//fmt.Printf("%d , ", i)
			coll := getRandomCollection(client)
			//fmt.Printf("inserting into collection %s, %v\n", coll.Name(), batch)
			_, err := coll.InsertMany(ctx, batch)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Inserted %d documents in %s/%s \n", batchSize, coll.Database().Name(), coll.Name())
		}()
	}
	wg.Wait()

	fmt.Println("Data insertion complete.")
}

func calcTotalStorageSize() int {
	return 0
}

func getRandomCollection(client *mongo.Client) *mongo.Collection {
	// Generate random index for databases and collections slices
	databaseIndex := rand.Intn(len(databases))
	collectionIndex := rand.Intn(len(collections))

	db := client.Database(databases[databaseIndex])
	collection := db.Collection(collections[collectionIndex])
	return collection
}
