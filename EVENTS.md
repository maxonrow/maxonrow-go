## Events
* `Hashing of the event is:`
* `20bytes of sha256, all small letter.`

# Fee
* 0x83e3e0f293c991a39fafe1556d8781fa1c1b6dfe - `CreatedFeeSetting(string,string)`  
`(owner: string, feeSettingName: string)`  

* 0x2591a6d8db4e2f0486d050a75270ea9d7ff915b4 - `CreatedTxFeeSetting(string,string)`    
`(owner: string, msgType: string)`  

# Kyc
* 0xf5b3603e2aa1e5fa9a7ec4a045d805146c7c9718 - `KycWhitelisted(string,string)`  
`(from: string, kycAddress: string)`  

* 0x645a6833044dc6a53b04e91633c49b303e5d0d85 - `RevokedWhitelist(string,string)`  
`(to: string, kycAddress: string)`  

*  0x3f6a6d6bea44fc3b2ed6f23a19d7c22a9410969e - `KycBinded(string,string,string)`  
`(from: string, to: string, kycAddress: string)`    

* 0x324b7ec81acbc6d956e34b88db1fd77acde375b7 - `KycUnbinded(string,string,string)`  
`(from: string, to: string, kycAddress: string)`  

# Nameservice(alias)
* 0x3e39f456518778478674f842ddd7d2c0f8523466 - `RejectedAlias(string,string)`  
`(alias: string, owner: string)`  
  
* 0x83e489ef72309acc291d6da10864e73efd0f486e - `RevokedAlias(string,string)`  
`(alias: string, owner: string)`  

* 0xdc152a0041a73ffbee055be6a7f477275604277b - `CreatedAlias(string,string,string,bignumber)`  
`(alias: string, from: string, feeTo: string, value: bignumber)`  

* 0x19563da5167518d666cc48872e6465787f4bdc36 - `ApprovedAlias(string,string)`  
`(alias: string, owner: string)`  

# Fungible Token
* 0xc44e0ed5b5505cde0e3d3233b63c2e267fc96e15 - `CreatedFungibleToken(string,string,string,bignumber)`  
`(symbol: string, owner: string, feeTo: string, value: bignumber)`  

* 0xc578a392f39d5e91824b9c61c8260f37440604bc - `ApprovedFungibleToken(string,string)`  
 `(symbol: string, owner: string)`  

* 0x254e90eaf2ba1d818c79b062d46f439acae94f98 - `RejectedFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x18ddae1cd03b74483e6e18282ffa58cb2b6ebbca - `FrozenFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0xb2f7c4ed24fe082fe74ad9428c3f0b987713b9ad- `UnfreezeFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x8743dad049ee684a0d6798c7976ffb97669c0b03 - `MintedFungibleToken(string,string,string,bignumber)`  
`(symbol: string, from: string, to: string, value: bignumber)`  

* 0x49caa496a951a0c139ddbdcf10de9186a711a2b3 - `TransferredFungibleToken(string,string,string,bignumber)`  
`(symbol: string, from: string, to: string, value: bignumber)`  

* 0x22ee863d0b0fff1cf78deaab1b29a5f185365b9a - `BurnedFungibleToken(string,string,string,bignumber)`  
`(symbol: string, owner: string, from: string, value: bignumber)`  

* 0x3d0a8cccf72a83d61ee15b35ce55be598ed90ae2 - `TransferredFungibleTokenOwnership(string,string,string)`  
`(symbol: string, from: string, to: string )`  

* 0x2cbc9a6ce2cbfab7c4d83608cda0c7799c8ee9a1 - `AcceptedFungibleTokenOwnership(string,string)`  
`(symbol: string, from: string)`  

* 0xa22424dec59cde1b4447efe16fccb342cc244411 - `FrozenFungibleTokenAccount(string,string)`  
`(symbol: string, owner: string)`  

* 0x9e9401fefbbf3e293218d8222fdfedb23cc7f476 - `UnfreezeFungibleTokenAccount(string,string)`  
`(symbol: string, owner: string)`  

# NonFungible Token
* 0x11f662e9d2dd1a7d0f582df94954113560fd7d30 - `CreatedNonFungibleToken(string,string,string,bignumber)`  
`(symbol: string, owner: string, feeTo: string, value: bignumber)`  

* 0x83138459f9be8326f1a9c3bc4184554f5e5fbbe3 - `ApprovedNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x5e1828ac38c50a2e6f3aebc769da1b780349f9c1 - `RejectedNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0xdbc39fe2e24eb3560f81989af3ac9260928d3750 - `FrozenNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x158b516fb26ccf2031abb7ef917f6ffca45146be - `UnfreezeNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0xd915a68b0824adae8500068f4791a15358b570be - `FrozenNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, owner: string)`  

* 0xf73aed41a76d330005eff923b7c47365ea5e5451 - `UnfreezeNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, owner: string)`  

* 0x6c38f2f21ca0cd532ab8e6bed6e64499f09d84d9 - `ApprovedTransferNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, owner: string, newOwner: string)`  

* 0x96721282448437667a3f958212613f7bbfcc9f68 - `RejectedTransferNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, owner: string, newOwner: string)`  

* 0xedac52598009550e6bb58d3b9518be7f0a3fd560 - `MintedNonFungibleItem(string,string,string,string)`  
`(symbol: string, itemID: string, from: string, to: string)`  

* 0x0602781615eb1bf8208da7fecab2aede96910073 - `TransferredNonFungibleItem(string,string,string,string)`  
`(symbol: string, itemID: string, from: string, to: string)`  

* 0xfa84eb62ed6bca35a9b95fb171fa34e5fa83a956 - `BurnedNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, from: string)`  

* 0xa0867945ed69b7ef4b6b7deb232c79946c68002d - `TransferredNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, from: string, to: string)`  

* 0xc906d530e1166a8b3477f78d8b1668d9c4b0c279 - `AcceptedNonFungibleTokenOwnership(string,string)`  
`(symbol: string, from: string)`  

  0x708401705567b885e3b8fa16b23271a8eea3d8f6
* 0x708401705567b885e3b8fa16b23271a8eea3d8f6 - `EndorsedNonFungibleItem(string,string,string,string)`  
`(symbol: string, itemID: string , from: string, metadata: string)`  

* 0x50691a25191c8a1223c954ba0bf81566be0b2c34 - `UpdatedNonFungibleItemMetadata(string,string,string)`  
`(symbol: string, itemID: string , from: string)`  

* 0xcdefef9a5959a9ca7863b3d92375a9c4cc24130f - `UpdatedNonFungibleTokenMetadata(string,string)`  
`(symbol: string, from: string)`  

# Bank send
* 0x2cadcfb0c336769d503d557b26fcf1e91819e7e5 - `Transferred(string,string,bignumber)`  
`(from: string, to: string, value: bignumber)`
