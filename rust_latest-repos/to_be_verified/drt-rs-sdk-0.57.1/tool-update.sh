#!/bin/bash

# Locally update the VM version

VM_TAG="v0.0.3"

echo "Before:"
drtpy config dump
drtpy config set dependencies.vmtools.tag $VM_TAG
drtpy config set dependencies.vmtools.urlTemplate.linux https://github.com/TerraDharitri/drt-go-chain-vm/archive/{TAG}.tar.gz
echo "After:"
drtpy config dump

drtpy deps install vmtools --overwrite
