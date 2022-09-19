#!/bin/bash

./apkctl uninstall platform

kubectl delete -f BackendService.yaml

kubectl delete httproute/petstore

kubectl delete httproute/demo-api