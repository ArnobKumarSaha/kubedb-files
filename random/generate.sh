#!/bin/bash
# Input:  ./extract_kubecfg_cert.sh my-cluster-name username
# Output: ./my-cluster-name-ca.crt ./username.crt ./username.key

# Exit on error
abort(){
  echo $1 && exit 1
}

# Prerequistes
for i in yq jq kubectl
do
 if ! which $i; then
  abort "$i is not instaled." 
 fi
done

cluster_name=$1
user=$2
echo "name =>  $cluster_name $user"




if ! kubectl config get-clusters | grep -q "^$cluster_name$"; 
then
  abort "Usage: $0 <cluster-name> <username>"
fi
if [ -z "$user" ];
then
  abort "Usage: $0 <cluster-name> <username>"
fi

kube_path=$HOME/.kube
kube_config=$kube_path/config
if [ ! -f $kube_config ];
then
  abort "No $kube_config file."
fi


echo "kube-config === $kube_config"

TMPJSON=$kube_path/kubecfg.json
# Convert yaml to json
# yq $kube_config -j > $TMPJSON

TMPJSON=config.json

# Get CA cert
cat $TMPJSON | jq --arg x $cluster_name -r \
	'.clusters[] | select(.name==$x) | .cluster | ."certificate-authority-data" ' | base64 -d > ${cluster_name}-ca.crt
if [ ! -s ${cluster_name}-ca.crt ];
then
  abort "Cannot find ${cluster_name}'s cert."
fi
# Get user client cert
cat $TMPJSON | jq --arg x $user -r \
	'.users[] | select(.name==$x) | .user | ."client-certificate-data" ' | base64 -d > $user.crt
if [ ! -s $user.crt ]; 
then
  abort "Cannot find $user's cert."
fi

echo "okokoko"
# Get user client key
cat $TMPJSON | jq --arg x $user -r \
	'.users[] | select(.name==$x) | .user | ."client-key-data" ' | base64 -d > $user.key
if [ ! -s $user.key ]; 
then
  abort "Cannot find $user's key."
fi

echo "${cluster_name}-ca.crt, $user.crt, and $user's key are generated in the current directory." 

# Clean up
rm -rf $TMPJSON
