# Version 1.3.5
- implement length limit for itemID(28)
- implement validation for burn flag upon burning token (ft/nft)
- implement length limit for nft endorser list(10)
- fix spelling error for "APPROVE_TRANSFER_TOKEN_OWNERSHIP" and "REJECT_TRANSFER_TOKEN_OWNERSHIP"

# Version 1.3.4
- fix when burn and item is frozen, burn is not allowed
- fix when ft/nft token if frozen there is no other action can be done until it is unfrozen. Except for nft
  endorsement

# Version 1.3.3
- enable endorser list to update to empty

# Version 1.3.2
- disable delegation handlers
- unregister codec for delegations

# Version 1.3.1
- split fee calculation between token (ft and nft)
- added fee cli for setting ft/ nft multiplier, ft/ nft token fee by token action
- added error code for checking update endorser(endorser has to be whitelisted)
- added function for updating nft endorser list
- update rpc endpoint for debug/fee_info to display nft and ft multiplier
- fee with token action (nft/ft) will return default fee if not found
- added query nonfungible fee collectors
- added test case to create default action fees for token(ft nft) actions
- added updateEndorserList as token action
- MakeEndorsement added a field metadata
- added test case for MakeEndorsement

# Version 1.3.0
- Disable multisig module
- Account query fix #103
- nft multisig test done
- enable codec for fungible and nonfungible token in multisig.
- symbol length change to 40
- added validation for decimal in ft (token decimal not allow to be 0 or greater than 18)
- disable rpc end point for listing token and nonfungible token 
- added another rpc for test (AccountCdc)
- remove cdc from querry account end point (Account)

# Version 1.2.1
- Adding non fungible token module
- Adding fungible token module
- Multi-signature account
- NFT endorsment and endorsment list
- NFT transfer limit
- MltiSig: Removing pending Tx after broadcasting
- NFT: Item owner only can transffer item
- NFT: Public token (if it is public, user only can mint to themselves for nft item)
- NFT: adding metadata and properties to the item
- NFT: Check mutibility for metadata and properties
- NFT: transfer limit and mint limit checks
- NFT: REJECT-TOKEN and REJECT-TRANSFER-OWNERSHIP
- Multisig: Broadcasting Internal transaction after approval
- MultiSig: check group account is it exist, validate internal transaction message.
- MultiSig: Delete pending Tx
- MultiSig: Signer without KYC
- Checking NFT-item metadata and properties length. Bug #31
- Signature verification without public key
- Multisig account with one signer. Bug #35
- Non-fungible token item can be retransffered. Bug #59
- re-mint nft token Bug #60
- Multisig re-submission a Tx. Bug #77
- Created unique address for multisig account upon creation
- Broadcasting internal transaction when rpc is disabled
- Craeting two multisig txs at the same time, Approve the second before the first
- Check token validity before processing tx.
- Fee for message types.
- Application fee should be deducted at ante handler.
- Fee multiplier.
- Alias letters and characters issue.
- Kyc bind/unbind
- RPC methods for debugging and query fees
- RPC methods for getting accont status 
- mxwcli enable to add key from mnemonic into keyring
- Old keybase deprecated  
- mxwcli generating multisig account address from sequence and owner address
- Updating tests

# Version 1.0.0
- Launching mainnet
- Updating version

# Version 0.7.4
- Added fee cli enable deleting account fee setting

# Version 0.7.3
- Fix fee multiplier/ token multiplier for fee calculation

# Version 0.7.2
- Fix fee setting multiplier cli
- Fix get account fee setting in fee calculation
- Add fee setting cli to enable setting fee setting for token with token action

# Version 0.7.1
- Non Fungible token implementation. (disabled)
- Added Fee setting cli, able to edit, delete existing system fee setting
- Added token listing rpc
- Remove genesis rpc
- Cosmos-sdk & tendermint upgraded

# Version 0.7.0
- Nameservice module released
- Token module released
- Maintenance module released
- Added test cases for fungible token
- Fungible token support dynamic and fixed supply
- Added test cases for nameservice (alias)
- Added test cases for maintenance
- Added MXW errors code handling
- Update Fee module to support Fungible/Non-Fungible token actions
- Added new RPC methods for debugging modules
- Minor bug fixes

# Version 0.6.3 (pre MainNet)
- Updating Tendermint to fix bug#62 (Each block only contains 1 or 2 transactions)
- Fixing Deadlock situation when sending bulk txs
- Updating cosmos (v0.37.0)


# Version 0.6.2 (pre MainNet)
- Fixing UATNet hangs after upgrading from 0.6.1
- Fixing app mismatch issue hash when a node re-run
- Fixing fee setting name in genesis file
- Add version RPC method
- gentx now support --moniker flag
- Calculate fee gets corrects message type and route
- Fixing some minor issues


# Version 0.6.1 (pre MainNet)
- Rejecting transactions with empty fee
- Updating system test for fee module
- Fixing some minor issues
- Updating cosmos (v0.36.0)


# Version 0.6.0 (pre MainNet)
- **create_non_empty_block**; Waiting for transactions to create a new block
- Removing slashing module. In MXW only verified validator can make blocks (validator set)
- Updating rewarding mechanism. Fees are distributed between all validator who sign a block. Not only proposer.
- Make sure `app_hash` doesn't change when a block is empty.
- Adding multiplier to fee settings to control fee economy.
- Improving fee calculation to overcome rounding issue for floating values.
- Registering MXW coin types to support [BIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
- Adding zero fee settings
- Adding new messages for assigning fees to modules and premium accounts (fee-selector)
- Adding new message for whitelisting validators (validator set)
- Assigning Fee & KYC messages to zero-fee-settings (KYC transactions are fee free)
- Adding new cli function for sending fee transactions (experimental)
- Encrypting private-key for key-ceremony application
- Updating Cosmos-sdk and Tendermint libraries
- Prompt for moniker in gentx (used for key-ceremony)
- Adding new RPC methods for querying whitelisted account
- Adding new RPC method for calculation fee
- Assigning fee-setting to special account (premium account). Useful for airdrop tokens.
- Returning code errors for special faults like "signature verification fault". It's useful for mobile sdk
- Adding tests for cli-send-tx
- Adding tests for fee-calculation
- Adding test for fee-keeper
- Adding system tests for bank/fees/kyc
- fixing some bugs and issues

# Version 0.5.2
- Added more test cases for kyc, revamp kyc test case
- Split kyc module kv store into two
- Added checking for revocation whitelist, to prevent revocation without whitelisting
- Added nameservice genesis set genesis alias owner

# Version 0.5.1
- Added test case for fee calculation for rounding issue
- Fee now calculated with Int instead of float
- Update uatnet deployment files


# Version 0.5.0
- Modify the coin type to 376 which already registered at https://github.com/satoshilabs/slips/blob/master/slip-0044.md
- Added kyc test cases
- Remove address from signature in kyc, nameservice, and token module


# Version 0.3.0
- Adding maintenance module (experimental, not tested)
- Cosmos version updated to the latest version (master)
- Staking supports 18 decimals
- Events structure updated
- Removing formTag and adding new function for creating mxwEvents `MakeMxwEvents`
- Using crypto.Pubkey in kyc/token/ns msgs instead of string and using amino marshaljson for signbytes.
- Adding Events.md for introducing mxw events
- Adding new command for initializing multiple nodes with gentxs: `mxwd init-auto`
- Adding new command for creating key-pairs for key ceremony: `mxwcli create-keypair`
- Updating tx for fee module: `mxwcli tx add-fee`
- Removing mint from genesis
- Minor fixing on the token and alias events
- Fixing some bugs

# Version 0.2.1
- Remove `omitempty` from ResultLog json field
- RevokedAlias event fixed, use alias owner instead of the tx signer
- Fixed big int issue on comparing min and max fee

# Version 0.2.0
- Removing mint module
- Re-designing the fee structure
- keepers won't reset when the blockchain restart
- Keepers instance (pointer) are shared between module
- Removing block reward and block fee from genesis file
- Revoking kyc features
- Updating queries for transaction fees
- Updating tests for rpc module
- Change the devnet image used to be the same as uatnet
- Modify the image used in uatnet to the correct version
- Improve error handling for calculate fees
- Improve implementation of validate signatures
- Updating bank-tx type from `cosmos` to `mxw`
- Fix zero fee txs. Failed transaction should not be recorded in blockchain

# Version 0.1.1
- Hot-fixing unmarshaling transaction in JSON format

# Version 0.1.0
- Reorganizing codebase
- Updating REST APIs for light client
- Adding JSON-RPC for query and broadcasting transaction in JSON format
- Adding JSON-RPC for encoding/decoding transactions
- Saving transactions in Binary format
- Updating units to handle 18 decimal places
- Fixing some bugs and issues
