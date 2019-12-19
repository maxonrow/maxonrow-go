## Events
* `Hashing of the event is:`
* `20bytes of sha256, all small letter.`

# Fee
* 0x83e3e0f293c991a39fafe1556d8781fa1c1b6dfe - `CreatedFeeSetting(string,string)`
* 0x2591a6d8db4e2f0486d050a75270ea9d7ff915b4 - `CreatedTxFeeSetting(string,string)`

# Kyc
* 0xf5b3603e2aa1e5fa9a7ec4a045d805146c7c9718 - `KycWhitelisted(string,string)`
* 0x645a6833044dc6a53b04e91633c49b303e5d0d85 - `RevokedWhitelist(string,string)`

# Nameservice(alias)
* 0x3e39f456518778478674f842ddd7d2c0f8523466 - `RejectedAlias(string,string)`
* 0xe9d027f679c35f0c185230d29d245cb110dd896f - `RemovedAlias(string,string)`
* 0xdc152a0041a73ffbee055be6a7f477275604277b - `CreatedAlias(string,string,string,bignumber)`
* 0x19563da5167518d666cc48872e6465787f4bdc36 - `ApprovedAlias(string,string)`

# Fungible Token
* 0xc44e0ed5b5505cde0e3d3233b63c2e267fc96e15 - `CreatedFungibleToken(string,string,string,bignumber)`
* `(symbol: string, owner: string, feeTo: string, value: bignumber)`
* 0xc578a392f39d5e91824b9c61c8260f37440604bc - `ApprovedFungibleToken(string,string)`
* `(symbol: string, owner: string)`
* 0x254e90eaf2ba1d818c79b062d46f439acae94f98 - `RejectedFungibleToken(string,string)`
* `(symbol: string, owner: string)`
* 0x18ddae1cd03b74483e6e18282ffa58cb2b6ebbca - `FrozenFungibleToken(string,string)`
* `(symbol: string, owner: string)`
* 0xb2f7c4ed24fe082fe74ad9428c3f0b987713b9ad- `UnfreezeFungibleToken(string,string)`
* `(symbol: string, owner: string)`
* 0x8743dad049ee684a0d6798c7976ffb97669c0b03 - `MintedFungibleToken(string,string,string,bignumber)`
* `(symbol: string, from: string, to: string, value: bignumber)`
* 0x49caa496a951a0c139ddbdcf10de9186a711a2b3 - `TransferredFungibleToken(string,string,string,bignumber)`
* `(symbol: string, from: string, to: string, value: bignumber)`
* 0x22ee863d0b0fff1cf78deaab1b29a5f185365b9a - `BurnedFungibleToken(string,string,string,bignumber)`
* `(symbol: string, owner: string, from: string, value: bignumber)`
* 0x3d0a8cccf72a83d61ee15b35ce55be598ed90ae2 - `TransferredFungibleTokenOwnership(string,string,string)`
* `(symbol: string, from: string, to: string )`
* 0x2cbc9a6ce2cbfab7c4d83608cda0c7799c8ee9a1 - `AcceptedFungibleTokenOwnership(string,string)`
* `(symbol: string, from: string)`
* 0xa22424dec59cde1b4447efe16fccb342cc244411 - `FrozenFungibleTokenAccount(string,string)`
* `(symbol: string, owner: string)`
* 0x9e9401fefbbf3e293218d8222fdfedb23cc7f476 - `UnfreezeFungibleTokenAccount(string,string)`
* `(symbol: string, owner: string)`

# Bank send
* 0x2cadcfb0c336769d503d557b26fcf1e91819e7e5 - `Transferred(string,string,bignumber)`
