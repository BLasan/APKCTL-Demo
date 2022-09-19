#!/bin/bash

./apkctl install platform

kubectl apply -f BackendService.yaml