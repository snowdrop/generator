#!/bin/bash

# We assume that the project generator already exists and that we are authenticated
oc project generator

# Delete existing resources
oc delete cm/generator-configmap

oc delete -f docker/generator.yml
oc create configmap generator-configmap --from-file=conf/generator.yaml
oc apply -f docker/generator.yml