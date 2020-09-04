# Maxonrow docker

This dockerfile is useful for running maxonrow locally and testing purpose. It has some kyc accounts. Some fee settings are set and also maintenance group are defined.

First of all there are two version of docker file which is one from master branch(Dockerfile) and develop branch(Dockerfile.develop). The develop branch will be the latest unrelease version, as for master branch it is the version of the latest released.

This docker file has pre defined [genesis](../tests/config/genesis.json) file and [config](../tests/config/config.toml) file.

Accounts including private keys are define [here](../tests/config/keys.json)

To build the docker run `docker build . --tag maxonrow`

To start the container, run `docker run -p 26656:26656 -p 26657:26657 --name maxonrow -d maxonrow`

To stop container, run `docker stop maxonrow`

Check rpc endpoint: http://localhost:26657
