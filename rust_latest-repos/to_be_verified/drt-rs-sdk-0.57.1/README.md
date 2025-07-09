# The Dharitri Rust Tool Set


[![Build Status](https://img.shields.io/github/actions/workflow/status/TerraDharitri/drt-rs-sdk/actions.yml?branch=main)](https://github.com/TerraDharitri/drt-rs-sdk/actions/workflows/actions.yml?query=branch%3Amain) [![Dependency Status](https://deps.rs/repo/github/TerraDharitri/drt-rs-sdk/status.svg)](https://deps.rs/repo/github/TerraDharitri/drt-rs-sdk) [![Contributors](https://img.shields.io/github/contributors/TerraDharitri/drt-rs-sdk)](https://github.com/TerraDharitri/drt-rs-sdk/graphs/contributors)

This repository contains a wide variety of tools, aimed primarily at smart contract developers.

The repo contains:
- The most complete smart contract framework on Dharitri:
    - The base framework;
    - A complete build system, which relies on the smart contract code directly;
    - A powerful debugger, based on a partial implementation of the Dharitri VM, in Rust.
    - A framework for writing both black-box and white-box tests. They rely on the standard Dharitri blockchain scenario format.
    - The official data serializer and deserializer for smart contract data. Can be used both on- and off-chain.
- A large collection of smart contract examples and feature tests, together with some of the core smart contracts used on the blockchain (e.g. the wrapped rewa swap, multisig, etc.).
- A framework for interacting with the blockchain, based on the smart contract logic, especially suitable for developers.
- A code snippet generator.

# Documentation

Most documentation can be found at https://docs.dharitri.org/developers/overview/

# Getting started

The crowdfunding tutorial is a great place to start: https://docs.dharitri.org/developers/tutorials/crowdfunding-p1/

# IDE

The framework is designed to be easiest to use with the IDE VSCode extension: https://marketplace.visualstudio.com/items?itemName=Numbat.vscode-numbat-ide

# Building contracts

A comprehensive build guide can be found here: https://docs.dharitri.org/developers/developer-reference/sc-build-reference/

# Debugging contracts

The debugger guide: https://docs.dharitri.org/developers/developer-reference/sc-debugging/
