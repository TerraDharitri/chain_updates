#!/bin/sh

# Exit with nonzero exit code if anything fails
set -e

# Install prerequisites
echo "Installing yarn..."
npm install --global yarn
echo "Installing yalc..."
npm install -global yalc


# Prepare drt-sdk-dapp for publishing
git clone https://github.com/TerraDharitri/drt-sdk-dapp.git

echo "cd drt-sdk-dapp..."
cd drt-sdk-dapp
git checkout development

echo "Installing dependencies for drt-sdk-dapp..."
yarn install

echo "Building drt-sdk-dapp..."
yarn build

echo "Publishing drt-sdk-dapp..."
cd dist
yalc publish
cd ../..


# Consume drt-sdk-dapp in drt-template-dapp
git clone https://github.com/TerraDharitri/drt-template-dapp.git

echo "cd drt-template-dapp..."
cd drt-template-dapp

echo "Installing dependencies drt-template-dapp..."
yarn install

echo "Linking drt-sdk-dapp..."
yalc add @terradharitri/sdk-dapp

echo "Building drt-template-dapp..."
yarn build:devnet


# Consume drt-sdk-dapp in drt-nextjs-template-dapp
git clone https://github.com/TerraDharitri/drt-nextjs-template-dapp.git

echo "cd drt-nextjs-template-dapp..."
cd drt-nextjs-template-dapp

echo "Installing dependencies drt-nextjs-template-dapp..."
yarn install

echo "Linking drt-sdk-dapp..."
yalc add @terradharitri/sdk-dapp

echo "Building drt-nextjs-template-dapp..."
yarn build-devnet

echo "Script executed successfully!"
