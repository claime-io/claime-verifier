#!/usr/bin/env bash

git clone https://github.com/bridges-inc/claime-registry.git
cd claime-registry
git fetch
git pull origin develop
yarn && npx hardhat compile
cd ..
