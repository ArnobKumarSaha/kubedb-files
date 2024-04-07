#!/bin/bash
# usage: ./k8s-cluster-setup.bash db backup builtin-prom operator-prom


LICENSE=/home/arnob/yamls/license/gke.txt
KUBEDB=v2024.2.14
KUBESTASH=v2024.2.22

INSTALL_DB=false
INSTALL_BACKUP=false
INSTALL_BUILTIN=false
INSTALL_OPERATOR=false


for arg in "$@"; do
    case $arg in
        "db")
            INSTALL_DB=true
            ;;
        "backup")
            INSTALL_BACKUP=true
            ;;
        "builtin-prom")
            INSTALL_BUILTIN=true
            ;;
        "operator-prom")
            INSTALL_OPERATOR=true
            ;;
        *)
            echo "Unknown parameter: $arg"
            ;;
    esac
done



if [ "$INSTALL_DB" = true ]; then
echo "Installing DB"
helm upgrade -i kubedb oci://ghcr.io/appscode-charts/kubedb \
    --version ${KUBEDB} \
    --namespace kubedb --create-namespace \
    --set-file global.license=${LICENSE} \
    --set kubedb-provisioner.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-webhook-server.server.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-ops-manager.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-autoscaler.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set global.featureGates.RabbitMQ=true \
    --set global.featureGates.ZooKeeper=true \
    --set global.featureGates.Solr=true \
    --set global.featureGates.Druid=true \
    --set global.featureGates.Pgpool=true \
    --set global.featureGates.FerretDB=true \
    --set global.featureGates.Singlestore=true \
    --set kubedb-provisioner.imagePullPolicy=Always \
    --set kubedb-webhook-server.imagePullPolicy=Always \
    --set kubedb-provisioner.operator.tag=v0.43.0-petset-0 \
    --set kubedb-ops-manager.operator.tag=v0.30.0-petset.0
fi



if [ "$INSTALL_BACKUP" = true ]; then
echo "Installing Backup"
helm install kubestash oci://ghcr.io/appscode-charts/kubestash \
  --version ${KUBESTASH} \
  --namespace kubestash --create-namespace \
  --set-file global.license=${LICENSE}
kubectl apply -f ~/go/src/github.com/random/external-snapshotter/client/config/crd/ || true
kubectl apply -f ~/go/src/github.com/random/external-snapshotter/deploy/kubernetes/snapshot-controller/ || true
kubectl apply -f ~/go/src/github.com/random/external-snapshotter/deploy/kubernetes/csi-snapshotter/ || true

sleep 5
kubectl create ns demo || true
kubectl apply -f ~/yamls/archiver/new/volumesnapshotclass.yaml
kubectl create secret generic -n demo linode-secret \
  --from-literal=AWS_ACCESS_KEY_ID=${LINODE_ACCESS_KEY_ID} \
  --from-literal=AWS_SECRET_ACCESS_KEY=${LINODE_SECRET_ACCESS_KEY} \
  --from-literal=AWS_ENDPOINT=https://us-east-1.linodeobjects.com
kubectl apply -f ~/yamls/archiver/new/retention-policy.yaml
kubectl apply -f ~/yamls/archiver/new/encrypt-secret.yaml
kubectl apply -f ~/yamls/archiver/new/backupstorage-lke.yaml
kubectl apply -f ~/yamls/archiver/new/archiver.yaml
kubectl apply -f ~/yamls/archiver/new/db1.yaml
fi


if [ "$INSTALL_BUILTIN" = true ] || [ "$INSTALL_OPERATOR" = true ]; then
helm install metrics-server metrics-server/metrics-server --set=args={--kubelet-insecure-tls}
helm upgrade -i kubedb-metrics appscode/kubedb-metrics -n kubedb --create-namespace --version=${KUBEDB}
fi


if [ "$INSTALL_BUILTIN" = true ]; then
echo "Installing built-in prometheus"
kubectl create ns monitoring || true
kubectl apply -f ~/yamls/builtin-prometheus/deployment.yaml
kubectl apply -f ~/yamls/builtin-prometheus/configmap.yaml
kubectl apply -f ~/yamls/builtin-prometheus/rbac.yaml
kubectl apply -f ~/yamls/builtin-prometheus/grafana.yaml
kubectl apply -f ~/yamls/builtin-prometheus/kube-state-metrics.yaml
kubectl apply -f ~/yamls/builtin-prometheus/svc.yaml

helm install panopticon appscode/panopticon -n kubeops \
    --create-namespace \
    --set monitoring.enabled=true \
    --set monitoring.agent=prometheus.io/builtin \
    --set-file license=${LICENSE}
fi


if [ "$INSTALL_OPERATOR" = true ]; then
echo "Installing prometheus operator"
helm install prometheus prometheus-community/kube-prometheus-stack --create-namespace -n monitoring
helm install panopticon appscode/panopticon -n kubeops \
    --create-namespace \
    --set monitoring.enabled=true \
    --set monitoring.agent=prometheus.io/operator \
    --set monitoring.serviceMonitor.labels.release=prometheus \
    --set-file license=${LICENSE}
fi






# helm install gg appscode/kubedb-grafana-dashboards -n kubeops --create-namespace --version=v2024.2.14 -f /home/arnob/yamls/grafana.yaml


helm upgrade -i kubedb oci://ghcr.io/appscode-charts/kubedb \
    --version v2024.2.14 \
    --namespace kubedb --create-namespace \
    --set-file global.license=/home/arnob/yamls/license/aws.txt \
    --set kubedb-provisioner.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-webhook-server.server.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-ops-manager.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-autoscaler.operator.securityContext.seccompProfile.type=RuntimeDefault \
    --set kubedb-provisioner.imagePullPolicy=Always \
    --set kubedb-webhook-server.imagePullPolicy=Always \
    --set kubedb-provisioner.operator.tag=v0.43.0-petset-0 \
    --set kubedb-ops-manager.operator.tag=v0.30.0-petset.0 \
    --set kubedb-webhook-server.server.tag=v0.19.0-petset.0 \
    --set kubedb-crd-manager.image.tag=v0.0.7-petset.0 \
    --set petset.enabled=true