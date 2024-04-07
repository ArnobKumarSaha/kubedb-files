#! /bin/bash

while true
do
	kubectl patch service prometheus-kube-prometheus-kubelet -n kube-system --type='json' -p='[{"op": "remove", "path": "/spec/ports/2"}, {"op": "remove", "path": "/spec/ports/1"}]'
	sleep 10
done


