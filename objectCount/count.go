package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
	"math"
)

func CompareObjectsCount(pc, sc *mongo.Client) (int64, error) {
	res := make(map[string]interface{})

	err := pc.Database("admin").RunCommand(context.Background(), bson.D{{Key: "listDatabases", Value: "1"}}).Decode(&res)
	if err != nil {
		return -1, err
	}

	masterDBs, ok := res["databases"]
	if !ok {
		return -1, fmt.Errorf("can't list master databases")
	}

	masterDbObjectCount := make(map[string]int64)
	var dbList []string
	klog.Info(fmt.Sprintf("Databases listed in the master: %v", len(masterDBs.(primitive.A))))

	for i := 0; i < len(masterDBs.(primitive.A)); i++ {
		masterDB := masterDBs.(primitive.A)[i].(map[string]interface{})
		dbname := masterDB["name"].(string)
		if !skipDB(dbname) {
			dbList = append(dbList, dbname)
			objectCounts := int64(0)
			collectionNames, err := pc.Database(dbname).ListCollectionNames(context.TODO(), bson.D{})
			if err != nil {
				return -1, err
			}

			for _, collectionName := range collectionNames {
				if !skipCollection(collectionName) {
					count, err := pc.Database(dbname).Collection(collectionName).CountDocuments(context.TODO(), bson.D{})
					if err != nil {
						return -1, err
					}
					klog.Info(fmt.Sprintf("master db=%v , collection=%v , docCount=%v", dbname, collectionName, count))

					objectCounts = objectCounts + count
				}
			}

			masterDbObjectCount[dbname] = objectCounts
		}
	}

	klog.Info(fmt.Sprintf("Object count done for master ; Now starting for %s", secondaryPod))

	for _, dbname := range dbList {
		objectCounts := int64(0)
		collectionNames, err := sc.Database(dbname).ListCollectionNames(context.TODO(), bson.D{})
		if err != nil {
			return -1, err
		}

		for _, collectionName := range collectionNames {
			if !skipCollection(collectionName) {
				count, err := sc.Database(dbname).Collection(collectionName).CountDocuments(context.TODO(), bson.D{})
				if err != nil {
					return -1, err
				}

				klog.Info(fmt.Sprintf("currentPod db=%v , collection=%v , docCount=%v", dbname, collectionName, count))
				objectCounts = objectCounts + count
			}
		}

		masterCount := masterDbObjectCount[dbname]
		replicaCount := objectCounts
		diff := (math.Abs(float64(masterCount-replicaCount)) / float64(masterCount)) * 100.0
		klog.Info(fmt.Sprintf("masterCount=%v, replicaCount=%v, diff=%v", masterCount, replicaCount, diff))
		if diff > 0 {
			return replicaCount, fmt.Errorf("object count of database %v didn't match. Master Database Object Count = %v and Replica Database Object count = %v", dbname, masterCount, replicaCount)
		}
	}

	return 0, nil
}

func skipDB(dbname string) bool {
	return dbname == "admin" ||
		dbname == "config" ||
		dbname == "local"
}

func skipCollection(collectionName string) bool {
	return collectionName == "system.profile" ||
		collectionName == "system.js" ||
		collectionName == "system.views"
}
