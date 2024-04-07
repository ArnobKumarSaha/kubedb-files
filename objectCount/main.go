package main

import (
	"context"
	"k8s.io/klog/v2"
	kubedb "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var (
	kbClient client.Client
	mongodb  kubedb.MongoDB
	password string

	primaryPod   string
	secondaryPod string
)

func main() {
	config := getRESTConfig()
	kbClient = getClient(config)
	err := getMongoDB()
	if err != nil {
		klog.Fatal(err)
	}

	tunnel, err := TunnelToDBService(config)
	if err != nil {
		return
	}

	err = getPrimaryAndSecondary(tunnel)
	if err != nil {
		return
	}

	tunnelPod, err := TunnelToDBPod(config, secondaryPod)
	if err != nil {
		return
	}

	start := time.Now()
	pc, sc := connectToMongo(tunnel, tunnelPod)
	defer func() {
		if err := pc.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		if err := sc.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	_, err = CompareObjectsCount(pc, sc)
	if err != nil {
		return
	}
	klog.Infof("Compares took %s", time.Since(start))
}
