Kube Aliyun
===========

[![CircleCI](https://circleci.com/gh/kubeup/kube-aliyun/tree/master.svg?style=shield)](https://circleci.com/gh/kubeup/kube-aliyun)

Aliyun essentials for Kubernetes. It provides SLB, Routes controllers and a Volume
plugin for Kubernetes to function properly on Aliyun instances.

Features
--------

* Service load balancers sync (TCP & UDP)
* Routes sync
* Volumes / PersistentVolumes
* Centralized attach/detach control. No credentials needed on nodes.
* PersistentVolumes dynamic provisioning

Docker Image
------------

- kubeup/kube-aliyun
- registry.aliyuncs.com/kubeup/kube-aliyun

Dependency
--------

Due to Kubernetes v1.6.0 flexvolume api break, the master currently supports k8s v1.6+.

To use the last v1.5.x support versionn, use the tag k8s-1.5.

Components
----------

There are two components.

**aliyun-controller** is a daemon responsible for service & route synchronization,
attach/detach control and PV provisioning. It has to run on all master nodes.

**aliyun-flexv** is a binary plugin responsible for volumes operations on nodes.
It has to be deployed on all nodes and will be called by kubelets/controller-manager when
 needed.

Deploy to Aliyun
----------------

### aliyun-controller

1. Make sure all node names are internal ip addresses.
2. Make sure node cidr will be allocated by adding `--allocate-node-cidrs=true
--configure-cloud-routes=false` to **kube-controller-manager** commandline.
3. Update the required fields in `manifests/aliyun-controller.yaml`
4. Upload it to `pod-manifest-path-of-kubelets` on all your master nodes
5. Use docker logs to check if the controller is running properly

** If your nodes can't access Aliyun metadata somehow, you need to specify 3 more
variables in env:

 - ALIYUN_REGION
 - ALIYUN_VPC
 - ALIYUN_VSWITCH

### aliyun-flexv

1. Add to **kubelet** commandline an option `--volume-plugin-dir=/opt/k8s/volume/plugins`
2. Add to **kube-controller-manager** commandline an option `--flex-volume-plugin-dir=/opt/k8s/volume/plugins`
3. Add two env variables to **kube-controller-manager**:

 - ALIYUN_ACCESS_KEY
 - ALIYUN_ACCESS_KEY_SECRET

4. Make flexv binary available on every node in a `./ailyun~flexv/` folder under
the kubelet volume plugin path. Or for your convenience, run this

```bash
  FLEXPATH=/opt/k8s/volume/plugins/aliyun~flexv; sudo mkdir $FLEXPATH -p; docker run -v $FLEXPATH:/opt kubeup/kube-aliyun:master cp /flexv /opt/
```

** Customizing volume plugin path is optional. You can just use the default which is
`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`.

Usage
-----

### Services

Just create Loadbalancer Services as usual. Currently only TCP & UDP types are
supported. Some options can be customized through annotaion on Service. Please
see [pkg/cloudprovider/providers/aliyun/loadbalancer.go](https://github.com/kubeup/kube-aliyun/blob/master/pkg/cloudprovider/providers/aliyun/loadbalancer.go) for details.

### Routes

Since we are using k8s to allocate node cidrs, we need a way to make that effective on
containers. There are several ways to do this.

* Use kubenet plugin. Details [here](https://kubernetes.io/docs/concepts/cluster-administration/network-plugins/#kubenet)
* Pass `--bip={subnet}` and `--ip-masq=false` to docker daemon

### Volumes

[example](https://github.com/kubeup/kube-aliyun/blob/master/examples/volume.yaml)

Use flexVolume in any volume spec. 

### Static PersistentVolumes

[example](https://github.com/kubeup/kube-aliyun/blob/master/examples/pv.yaml)

** Recycling policy is not supported. Use custom recycler pod if you want to.

### Dynamic PersistentVolumes and StorageClass

[example](https://github.com/kubeup/kube-aliyun/blob/master/examples/dynamic-pv.yaml)

Avaialable parameters on StorageClass:

  - diskCategory: Disk category as in Aliyun doc. Required.
  - fsType: Filesystem type. Default: ext4

More Examples
-------------

Please find a more complete setup example [here](https://github.com/kubeup/archon/tree/master/example/k8s-aliyun) which **Archon** is able to deploy 
automatically. 

License
-------

Apache Version 2.0
