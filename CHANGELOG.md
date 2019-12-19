# Version 1.0.0
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