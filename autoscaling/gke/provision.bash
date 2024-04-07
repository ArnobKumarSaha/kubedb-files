#!/bin/bash

CLUSTER_NAME=arnob
ZONE=us-central1-c

gcloud container clusters create "${CLUSTER_NAME}" \
  --zone="${ZONE}" \
  --num-nodes=3 \
  --machine-type=e2-standard-4 \
  --enable-autoupgrade \
  --enable-autorepair  # --node-taints=nodepool=default:NoSchedule 

gcloud container node-pools create medium \
  --cluster="${CLUSTER_NAME}" \
  --zone="${ZONE}" \
  --machine-type=e2-medium \
  --num-nodes=0 \
  --enable-autoscaling \
  --min-nodes=0 \
  --max-nodes=3 \
  --node-taints=nodepool=medium:NoSchedule \
  --node-labels=database=yes

gcloud container node-pools create std4 \
  --cluster="${CLUSTER_NAME}" \
  --zone="${ZONE}" \
  --machine-type=e2-standard-4 \
  --num-nodes=0 \
  --enable-autoscaling \
  --min-nodes=0 \
  --max-nodes=3 \
  --node-taints=nodepool=std4:NoSchedule \
  --node-labels=database=yes

gcloud container node-pools create std8 \
  --cluster="${CLUSTER_NAME}" \
  --zone="${ZONE}" \
  --machine-type=e2-standard-8 \
  --num-nodes=0 \
  --enable-autoscaling \
  --min-nodes=0 \
  --max-nodes=3 \
  --node-taints=nodepool=std8:NoSchedule \
  --node-labels=database=yes

# gcloud container node-pools update default-pool \
#   --cluster="${CLUSTER_NAME}" \
#   --zone="${ZONE}" \
#   --node-taints=nodepool=default:NoSchedule \
#   --quiet