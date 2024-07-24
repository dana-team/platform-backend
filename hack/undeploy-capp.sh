#!/bin/bash

capprepo=https://github.com/dana-team/container-app-operator
branch=main

git clone ${capprepo} -b ${branch}
make -C container-app-operator uninstall-prereq
make -C container-app-operator undeploy
rm -rf container-app-operator