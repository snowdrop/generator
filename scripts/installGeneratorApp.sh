#!/bin/bash

# We assume that the project generator already exists and that we are authenticated
oc project generator

# Delete existing resources
echo "# Delete generator's k8s resources"
oc delete cm/generator-configmap
oc delete -f docker/generator-application.yml

echo "# Populate a new ConfigMap"
oc create configmap generator-configmap --from-file=conf/generator.yaml

echo "# Deploy the generator k8s resources"
oc apply -f docker/generator-application.yml