#!/bin/sh

kind delete cluster || true
kind create cluster --config=/home/arnob/yamls/affinity/kind.config

kubectl label node kind-control-plane topology.kubernetes.io/zone=us-west-2b
kubectl label node kind-worker topology.kubernetes.io/zone=us-west-2c
kubectl label node kind-worker2 topology.kubernetes.io/zone=us-west-2d

kubectl label node kind-control-plane karpenter.sh/capacity-type=on-demand
kubectl label node kind-worker karpenter.sh/capacity-type=spot
kubectl label node kind-worker2 karpenter.sh/capacity-type=spot
