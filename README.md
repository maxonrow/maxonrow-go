
# Maxonrow

Maxonrow is the sole main chain service provider in the world to offer real-name identity verification for KYC process and AML safety measures. Its a blockchain powered by Tendermint & Cosmos-SDK. We utilize GDPOS (Governed Delegated Proof of Stake) as a consensus mechanism, which is a mechanism of equity entrustment unique to Maxonrow.

# Building the source

Prerequisites
* Go 1.12.5+ [installed](https://github.com/golang/go)
* Docker [installed](https://docs.docker.com/engine/installation/)(Optional)
* Docker Compose [installed](https://docs.docker.com/compose/install/)(Optional)


## Steps for building the project executables

1. Clone the project repository

    `git clone https://github.com/mxw-blockchain/maxonrow-go/`


2. Change to the project directory.

    `cd  github.com/mxw-blockchain/maxonrow-go/`


3.  Get all the dependecies and build project the binary

    `make all`


## Running `MXWONROW` Full node on the local network

1. Initial the blockchain network with a valid chain-id. Below command which genrate the config,genesis and account in respective folder


    `./build/mxwd init --home ~/.mxw`


2. Generate validator transcation. Below command which create the gentx folder with gentx transcation of validator account-1

    `./build/mxwd gentx --name acc-1 --home ~/.mxw`


3. Run the blockchain using the below command

    `./build/mxwd start --home ~/.mxw`


### Acoount and Transaction

  - Get all the account details

    `./build/mxwcli keys list --home ~/.mxw`

  - Query the account

    `mxwcli query account $(mxwcli keys show acc-1 --address) --chain-id maxonrow-chain`, which verify the balance of acc-1

  - Send Transaction

    `mxwcli tx send $(mxwcli keys show acc-1 --address) $(mxwcli keys show acc-2 --address) 1000cin --fees 10000000000000000cin --gas 0 --memo "TRANSFER" --chain-id maxonrow-chain`


## `MXWONROW` come with the many test cases,you can find all our test case under test folder in project soruce.

* To run our test case:

    `make test`