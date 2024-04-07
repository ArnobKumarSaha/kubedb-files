My tries : 
 - Tried with both multi-node kind cluster & single node kind cluster.
 - Installed `kind`, `kubectl` by hand.



## Docker Info

```
runner@acdc-3-3:~/runner/_work/mongodb/mongodb$ which docker 
/usr/bin/docker


runner@acdc-3-3:~/runner/_work/mongodb/mongodb$ docker images
REPOSITORY                    TAG                                           IMAGE ID       CREATED         SIZE
ghcr.io/kubedb/mg-operator    v0.34.0-beta.1-1-gbb0716221-dbg_linux_amd64   1f41fe28a039   2 minutes ago   283MB
ghcr.io/kubedb/mg-operator    v0.34.0-beta.1-1-gbb0716221_linux_amd64       45d75e256e61   2 minutes ago   169MB
ghcr.io/appscode/golang-dev   1.21                                          dcaba40783a0   13 days ago     1.51GB
kindest/node                  <none>                                        89e7dc9f9131   7 months ago    932MB


$ docker save ghcr.io/kubedb/mg-operator:v0.34.0-beta.1-1-gbb0716221_linux_amd64 > my.tar
$ tar -xvf my.tar 
blobs/
blobs/sha256/
blobs/sha256/10e9b9e72178b3e5fe8ff566a61eefb03ad585f5eae5cc0f2b08dbe3c874cecb
blobs/sha256/1766806cf32c3d7ff0faa9f5630b74fcdcc14e33a378b019eca480fac32c58d1
blobs/sha256/1a73b54f556b477f0a8b939d13c504a3b4f4db71f7a09c63afbc10acb3de5849
blobs/sha256/3a0fd3584ebcfaffef1040bd3004a26a8d36f2f2bd1e802fbf29238a690fbf85
blobs/sha256/45d75e256e6151910dd36136f9cdda2fbb1b460b2daa27f110da13656791a61e
blobs/sha256/464d6bd6658b01f626665b55bc0028e5a6fab4e49d738e300214fc09f9fe81e7
blobs/sha256/6f1e3682d66abd0d92095699687c4c9a03b87f68a65b5dd9b79a4c1e87309936
blobs/sha256/6f4d4102c049625bda5ffcd0e671ae65aea1d600735d3dc3585309f7d99cc083
blobs/sha256/70c791e25c15d8b3f812bc972959f53be23d34c4579706b2e342f53b4e6b3aa4
blobs/sha256/8700526eb61e505508e6b3d4e0b31f78a4eb5f98502ead6deeb16c8a53587740
blobs/sha256/8f2b2d741d53c8e25da6a3b40a1f45aba2612cf5c1ff06e642286358e756d9a8
blobs/sha256/9921b5b9328039b1348a586f4cecae542339a41f9eaebb26b3bbf3924d5dd675
blobs/sha256/9959f8de333696342eedd49db37331a9c3efa84899f1236226c18db2e856f2aa
blobs/sha256/b76ced86194643349ac6a239046e75a387c4106cf83d7cef5d8c634ead18e365
blobs/sha256/b9dff753472714a6a0c13099a746cf3f397b05c6d0ebfae4ce1ab04e38c26802
blobs/sha256/bcf97653e3511f8357ed3ee6a9dd5e3ff9507364d6ebde6de7fe37bd06bc4a96
blobs/sha256/d52f02c6501c9c4410568f0bf6ff30d30d8290f57794c308fe36ea78393afac2
blobs/sha256/dfd833c306f4709fc39a5cd6028f5fe130a6199e156c9042f4defbecd447ef0e
blobs/sha256/e486c66c8915b29dc4241fd6154d237182d727db74512249eac3e648974b424f
blobs/sha256/e5c3350918fab9f177b9497466174b2a29bd622a56d0991741a2ee7d321a52e8
blobs/sha256/e624a5370eca2b8266e74d179326e2a8767d361db14d13edd9fb57e408731784
blobs/sha256/ff5700ec54186528cbae40f54c24b1a34fb7c01527beaa1232868c16e2353f52
index.json
manifest.json
oci-layout
repositories

$ cat manifest.json 

``````




## Inside kind Node

`docker exec -it kind-control-plane bash`

```
root@kind-control-plane:/# ctr -n=k8s.io image list | grep mg-operator
ghcr.io/kubedb/mg-operator:v0.34.0-beta.1-1-gbb0716221_linux_amd64                        application/vnd.oci.image.manifest.v1+json                sha256:f391c20ff102a8af4bef5bde49b73842c3b88dd28ef2c4bef940e3fba8faf879 161.9 MiB linux/amd64                                                                   io.cri-containerd.image=managed 

root@kind-control-plane:/# crictl images | grep mg-operator
ghcr.io/kubedb/mg-operator                 v0.34.0-beta.1-1-gbb0716221_linux_amd64   45d75e256e615       170MB
```

Notice that, `mg-operator` image is available inside the kind node. But its a wrong information. We cant use this image in any pod. Doing that says, ImagePullBackOff.




Followed these steps : 
https://github.com/containerd/containerd/issues/4887#issuecomment-1371151081

It seems like the issue is resolved, but actually not. 
When We create a pod with that image (ImagePullPolicy: IfNotPresent), it gives the same `wrong diff id calculated on extraction` error.




docker exec --privileged -i kind-control-plane ctr --namespace=k8s.io images import --all-platforms --digests --snapshotter=overlayfs -
docker exec -ti kind-control-plane ctr -n=k8s.io image import /image-filename.tar --no-unpack --all-platforms --digests