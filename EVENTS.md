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

*  0x3F6A6D6BEA44FC3B2ED6F23A19D7C22A9410969E - `KycBinded(string,string,string)`  
`(from: string, to: string, kycAddress: string)`    

* 0x324B7EC81ACBC6D956E34B88DB1FD77ACDE375B7 - `KycUnbinded(string,string,string)`  
`(from: string, to: string, kycAddress: string)`  

# Nameservice(alias)
* 0x3e39f456518778478674f842ddd7d2c0f8523466 - `RejectedAlias(string,string)` 
`(alias: string, owner: string)`  
  
* 0x83E489EF72309ACC291D6DA10864E73EFD0F486E - `RevokedAlias(string,string)`  
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
* 0x11F662E9D2DD1A7D0F582DF94954113560FD7D30 - `CreatedNonFungibleToken(string,string,string,bignumber)`  
`(symbol: string, owner: string, feeTo: string, value: bignumber)`  

* 0x83138459F9BE8326F1A9C3BC4184554F5E5FBBE3 - `ApprovedNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x5E1828AC38C50A2E6F3AEBC769DA1B780349F9C1 - `RejectedNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0xDBC39FE2E24EB3560F81989AF3AC9260928D3750 - `FrozenNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0x158B516FB26CCF2031ABB7EF917F6FFCA45146BE - `UnfreezeNonFungibleToken(string,string)`  
`(symbol: string, owner: string)`  

* 0xD915A68B0824ADAE8500068F4791A15358B570BE - `FrozenNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, owner: string)`  

* 0xF73AED41A76D330005EFF923B7C47365EA5E5451 - `UnfreezeNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, owner: string)`  

* 0x6C38F2F21CA0CD532AB8E6BED6E64499F09D84D9 - `ApprovedTransferNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, owner: string, newOwner: string)`  

* 0x96721282448437667A3F958212613F7BBFCC9F68 - `RejectedTransferNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, owner: string, newOwner: string)`  

* 0xEDAC52598009550E6BB58D3B9518BE7F0A3FD560 - `MintedNonFungibleItem(string,string,string,string)`  
`(symbol: string, itemID: string, from: string, to: string)`  

* 0x0602781615EB1BF8208DA7FECAB2AEDE96910073 - `TransferredNonFungibleItem(string,string,string,string)`  
`(symbol: string, itemID: string, from: string, to: string)`  

* 0xFA84EB62ED6BCA35A9B95FB171FA34E5FA83A956 - `BurnedNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string, from: string)`  

* 0xA0867945ED69B7EF4B6B7DEB232C79946C68002D - `TransferredNonFungibleTokenOwnership(string,string,string)`  
`(symbol: string, from: string, to: string)`  

* 0xC906D530E1166A8B3477F78D8B1668D9C4B0C279 - `AcceptedNonFungibleTokenOwnership(string,string)`  
`(symbol: string, from: string)`  

* 0xE4018B8FDADE0A9A89DD3A65A5A6ED4389DE4F10 - `EndorsedNonFungibleItem(string,string,string)`  
`(symbol: string, itemID: string , from: string)`  

* 0x50691A25191C8A1223C954BA0BF81566BE0B2C34 - `UpdatedNonFungibleItemMetadata(string,string,string)`  
`(symbol: string, itemID: string , from: string)`  

* 0xCDEFEF9A5959A9CA7863B3D92375A9C4CC24130F - `UpdatedNonFungibleTokenMetadata(string,string)`  
`(symbol: string, from: string)`  

# Bank send
* 0x2cadcfb0c336769d503d557b26fcf1e91819e7e5 - `Transferred(string,string,bignumber)`  
`(from: string, to: string, value: bignumber)`
