package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scm = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scm))
	utilruntime.Must(kubedbscheme.AddToScheme(scm))
}

func getRESTConfig() *rest.Config {
	kubeConfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func getClient(config *rest.Config) client.Client {
	kc, err := client.New(config, client.Options{
		Scheme: scm,
		Mapper: nil,
	})
	if err != nil {
		panic(err)
	}
	return kc
}

func getMongoDB() error {
	name := os.Getenv("MONGODB_NAME")
	ns := os.Getenv("MONGODB_NAMESPACE")

	err := kbClient.Get(context.TODO(), types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}, &mongodb)
	if err != nil {
		return err
	}

	var authSecret corev1.Secret
	err = kbClient.Get(context.TODO(), types.NamespacedName{
		Namespace: ns,
		Name:      mongodb.Spec.AuthSecret.Name,
	}, &authSecret)
	if err != nil {
		return err
	}

	password = string(authSecret.Data["password"])
	return nil
}
