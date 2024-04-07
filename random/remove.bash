#!/bin/bash


names=$(kubectl get recommendation -n demo -o=jsonpath='{.items[*].metadata.name}')


for name in ${names}
do
    outdated=$(kubectl get recommendation -n demo $name -ojsonpath={.status.outdated})

    if [ $outdated == "true" ]
    then
        kubectl delete recommendation -n demo $name
    fi
done