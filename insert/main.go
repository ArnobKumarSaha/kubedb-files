package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	wg           sync.WaitGroup
	stopPrinting chan struct{}
	start        time.Time

	dataSize      = 1024 * 1024 * 1024
	batchSize     = 500000
	numGoroutines = 10
	databases     = []string{"aa", "bb", "cc"}
	collections   = []string{"one", "two", "three", "four"}
	checkpoints   []float64
)

func insert(ctx context.Context, client *mongo.Client, rt int) {
	defer wg.Done()
	fmt.Printf("Starts goroutine %d.\n", rt)
	for {
		select {
		case <-stopPrinting:
			fmt.Println("Stopping goroutine.")
			return
		default:
			batch := make([]interface{}, batchSize)
			for b := 0; b < batchSize; b++ {
				document := map[string]interface{}{
					"key":   rand.Int() % 50,
					"value": rand.Float64(),
				}
				batch[b] = document
			}
			coll := getRandomCollection(client)
			_, err := coll.InsertMany(ctx, batch)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Inserted %d documents in %s/%s \n", batchSize, coll.Database().Name(), coll.Name())
		}
	}
}

func getRandomCollection(client *mongo.Client) *mongo.Collection {
	// Generate random index for databases and collections slices
	databaseIndex := rand.Intn(len(databases))
	collectionIndex := rand.Intn(len(collections))

	db := client.Database(databases[databaseIndex])
	collection := db.Collection(collections[collectionIndex])
	return collection
}

func main() {
	start = time.Now()
	stopPrinting = make(chan struct{})
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
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)
	rand.Seed(time.Now().UnixNano())

	err = checkpointPreviousSizes(ctx, client)
	if err != nil {
		fmt.Printf("err on checkpointing = %v", err)
		return
	}
	go checkCondition(ctx, client)
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go insert(ctx, client, i)
	}
	wg.Wait()
}

func checkCondition(ctx context.Context, client *mongo.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sz, _ := calcTotalStorageSize(ctx, client)
			if int(sz) >= dataSize {
				fmt.Printf("Stopping all goroutines. Time elapsed: %v\n", time.Since(start))
				close(stopPrinting)
				return
			}
			convertToReadableUnit(sz)
		}
	}
}

func checkpointPreviousSizes(ctx context.Context, client *mongo.Client) error {
	checkpoints = make([]float64, len(databases))
	for i := 0; i < len(databases); i++ {
		databaseName := databases[i]
		dbStats, err := DBStats(ctx, client, databaseName)
		if err != nil {
			return err
		}
		storageSz, ok := dbStats["storageSize"]
		if !ok {
			return fmt.Errorf("type assertion error: can't get dataSize info")
		}

		checkpoints[i] = storageSz.(float64)
	}
	return nil
}

func calcTotalStorageSize(ctx context.Context, client *mongo.Client) (float64, error) {
	//dbList, err := ListDatabases(ctx, client)
	//if err != nil {
	//	return 0, err
	//}
	//
	//dbs, ok := dbList["databases"]
	//if !ok {
	//	return 0, fmt.Errorf("can't list database db")
	//}
	//
	//for i := 0; i < len(dbs.(primitive.A)); i++ {
	//	db, ok := dbs.(primitive.A)[i].(map[string]interface{})
	//	if !ok {
	//		return 0, fmt.Errorf("type assertion error: failed to get db info")
	//	}
	//	databaseName, ok := db["name"].(string)
	//	if !ok {
	//		return 0, fmt.Errorf("type assertion error: failed to get database name")
	//	}

	totalSize := float64(0)
	for i := 0; i < len(databases); i++ {
		databaseName := databases[i]
		dbStats, err := DBStats(ctx, client, databaseName)
		if err != nil {
			return 0, err
		}

		storageSz, ok := dbStats["storageSize"]
		if !ok {
			return 0, fmt.Errorf("type assertion error: can't get dataSize info")
		}

		totalSize += storageSz.(float64) - checkpoints[i]
	}
	return totalSize, nil
}

func ListDatabases(ctx context.Context, client *mongo.Client) (map[string]interface{}, error) {
	dbList := make(map[string]interface{})
	err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "listDatabases", Value: 1}}).Decode(&dbList)
	return dbList, err
}

func DBStats(ctx context.Context, client *mongo.Client, db string) (map[string]interface{}, error) {
	dbStats := make(map[string]interface{})
	err := client.Database(db).RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&dbStats)
	return dbStats, err
}

func convertToReadableUnit(sz float64) {
	divCount := 0
	for {
		if sz >= 1024 {
			sz = sz / 1024
			divCount += 1
		} else {
			break
		}
	}
	unit := "Byte"
	switch divCount {
	case 1:
		unit = "KiB"
		break
	case 2:
		unit = "MiB"
		break
	case 3:
		unit = "GiB"
		break
	case 4:
		unit = "TiB"
		break
	default:
		fmt.Printf("divCount %v is way too big \n", divCount)
	}
	fmt.Printf("%v %s inserted. Current time %v \n", sz, unit, time.Now())
}
