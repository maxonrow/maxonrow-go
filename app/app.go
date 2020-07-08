package app

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	sdkDist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	sdkParams "github.com/cosmos/cosmos-sdk/x/params"
	sdkStaking "github.com/cosmos/cosmos-sdk/x/staking"
	sdkSupply "github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/maxonrow/maxonrow-go/genesis"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
	"github.com/maxonrow/maxonrow-go/x/nameservice"
	fungible "github.com/maxonrow/maxonrow-go/x/token/fungible"
	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	rpccore "github.com/tendermint/tendermint/rpc/core"
	rpc "github.com/tendermint/tendermint/rpc/lib/server"
	tm "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	appName = "maxonrow-go"
)

var (
	ModuleBasics module.BasicManager
)

type mxwApp struct {
	*baseapp.BaseApp
	cdc       *codec.Codec
	txDecoder sdkTypes.TxDecoder
	txEncoder sdkTypes.TxEncoder
	logger    log.Logger

	// Handlers
	authAnteHandler sdkTypes.AnteHandler

	// Storage keys
	keyMain         *sdkTypes.KVStoreKey
	keyAccount      *sdkTypes.KVStoreKey
	keyNSnames      *sdkTypes.KVStoreKey
	keyNSowners     *sdkTypes.KVStoreKey
	keySupply       *sdkTypes.KVStoreKey
	keyStaking      *sdkTypes.KVStoreKey
	keyParams       *sdkTypes.KVStoreKey
	tkeyParams      *sdkTypes.TransientStoreKey
	keyDistr        *sdkTypes.KVStoreKey
	keyToken        *sdkTypes.KVStoreKey
	keyFee          *sdkTypes.KVStoreKey
	KeyKyc          *sdkTypes.KVStoreKey
	KeyKycData      *sdkTypes.KVStoreKey
	KeyMaintenance  *sdkTypes.KVStoreKey
	KeyValidatorSet *sdkTypes.KVStoreKey

	// Keepers
	accountKeeper          sdkAuth.AccountKeeper
	supplyKeeper           sdkSupply.Keeper
	bankKeeper             sdkBank.Keeper
	stakingKeeper          sdkStaking.Keeper
	distrKeeper            sdkDist.Keeper
	paramsKeeper           sdkParams.Keeper
	nsKeeper               nameservice.Keeper
	kycKeeper              kyc.Keeper
	fungibleTokenKeeper    fungible.Keeper
	nonFungibleTokenKeeper nonFungible.Keeper
	feeKeeper              fee.Keeper
	maintenanceKeeper      maintenance.Keeper

	router sdkTypes.Router

	mm *module.Manager

	chainID     string
	blockHeight int64
}

func init() {

	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.SetCoinType(376)
	config.SetFullFundraiserPath("44'/376'/0'/0/0")
	config.SetKeyringServiceName("mxw")

	config.Seal()

	// cosmos-sdk was using big.NewInt(6), but mxw is supporting 18 for staking.
	sdkTypes.PowerReduction = sdkTypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	ModuleBasics = module.NewBasicManager(
		sdkAuth.AppModuleBasic{},
		sdkBank.AppModuleBasic{},
		sdkStaking.AppModuleBasic{},
		sdkDist.AppModuleBasic{},
		sdkParams.AppModuleBasic{},
		sdkSupply.AppModuleBasic{},
	)
}

func NewMXWApp(logger log.Logger, db dbm.DB) *mxwApp {
	cdc := MakeDefaultCodec()

	base := baseapp.NewBaseApp(appName, logger, db, sdkAuth.DefaultTxDecoder(cdc))

	app := &mxwApp{
		BaseApp:         base,
		cdc:             cdc,
		logger:          logger,
		keyMain:         sdkTypes.NewKVStoreKey("main"),
		keyAccount:      sdkTypes.NewKVStoreKey("acc"),
		keyNSnames:      sdkTypes.NewKVStoreKey("ns_names"),
		keyNSowners:     sdkTypes.NewKVStoreKey("ns_owners"),
		keySupply:       sdkTypes.NewKVStoreKey("supply"),
		keyStaking:      sdkTypes.NewKVStoreKey("staking"),
		keyParams:       sdkTypes.NewKVStoreKey("params"),
		tkeyParams:      sdkTypes.NewTransientStoreKey("transient_params"),
		keyDistr:        sdkTypes.NewKVStoreKey("distr"),
		keyToken:        sdkTypes.NewKVStoreKey("token"),
		keyFee:          sdkTypes.NewKVStoreKey("fee"),
		KeyKyc:          sdkTypes.NewKVStoreKey("kyc"),
		KeyKycData:      sdkTypes.NewKVStoreKey("kycData"),
		KeyMaintenance:  sdkTypes.NewKVStoreKey("maintenance"),
		KeyValidatorSet: sdkTypes.NewKVStoreKey("validator_set"),
	}

	app.txDecoder = sdkAuth.DefaultTxDecoder(cdc)
	app.txEncoder = sdkAuth.DefaultTxEncoder(cdc)

	// account permissions
	maccPerms := map[string][]string{
		sdkAuth.FeeCollectorName:     nil,
		sdkDist.ModuleName:           nil,
		mint.ModuleName:              []string{sdkSupply.Minter},
		sdkStaking.BondedPoolName:    []string{sdkSupply.Burner, sdkSupply.Staking},
		sdkStaking.NotBondedPoolName: []string{sdkSupply.Burner, sdkSupply.Staking},
		gov.ModuleName:               []string{sdkSupply.Burner},
	}

	app.paramsKeeper = sdkParams.NewKeeper(
		app.cdc,
		app.keyParams,
		app.tkeyParams,
		sdkParams.DefaultCodespace,
	)

	app.accountKeeper = sdkAuth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(sdkAuth.DefaultParamspace),
		sdkAuth.ProtoBaseAccount,
	)

	app.bankKeeper = sdkBank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(sdkBank.DefaultParamspace),
		sdkBank.DefaultCodespace,
		app.ModuleAccountAddrs(maccPerms),
	)

	app.supplyKeeper = sdkSupply.NewKeeper(app.cdc,
		app.keySupply,
		app.accountKeeper,
		app.bankKeeper,
		maccPerms)

	// NS Keeper
	app.nsKeeper = nameservice.NewKeeper(
		app.keyNSnames,
		app.keyNSowners,
		&app.accountKeeper,
		app.bankKeeper,
		app.cdc,
	)

	app.stakingKeeper = sdkStaking.NewKeeper(
		cdc,
		app.keyStaking,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(sdkStaking.DefaultParamspace),
		sdkStaking.DefaultCodespace,
	)

	app.distrKeeper = sdkDist.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(sdkDist.DefaultParamspace),
		app.stakingKeeper,
		app.supplyKeeper,
		sdkDist.DefaultCodespace,
		sdkAuth.FeeCollectorName,
		app.ModuleAccountAddrs(maccPerms),
	)

	app.fungibleTokenKeeper = fungible.NewKeeper(cdc, &app.accountKeeper, &app.feeKeeper, app.keyToken)
	app.nonFungibleTokenKeeper = nonFungible.NewKeeper(cdc, &app.accountKeeper, &app.feeKeeper, app.keyToken)
	app.feeKeeper = fee.NewKeeper(cdc, app.keyFee)
	app.kycKeeper = kyc.NewKeeper(cdc, &app.accountKeeper, app.KeyKyc, app.KeyKycData)
	app.maintenanceKeeper = maintenance.NewKeeper(cdc, app.KeyMaintenance, app.KeyValidatorSet, app.executeProposal)

	// Registering hooks from distribution and slashing module to be called
	// on different events in the consensus
	app.stakingKeeper.SetHooks(
		sdkStaking.NewMultiStakingHooks(app.distrKeeper.Hooks()))

	// AnteHandler is executed before every transaction, it verifies transactions,
	// verifies their signatures and manages fees via feeCollectionKeeper
	app.authAnteHandler = app.NewAnteHandler()
	app.SetAnteHandler(app.anteHandler)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetInitChainer(app.initChainer)
	app.SetEndBlocker(app.endBlocker)

	app.mm = module.NewManager(
		sdkAuth.NewAppModule(app.accountKeeper),
		sdkBank.NewAppModule(app.bankKeeper, app.accountKeeper),
		sdkSupply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		sdkDist.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		sdkStaking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
	)

	app.Router().
		//AddRoute("auth", auth.NewHandler(app.accountKeeper, app.kycKeeper, app.txEncoder)).
		AddRoute("bank", bank.NewHandler(app.bankKeeper, app.accountKeeper)).
		AddRoute("staking", sdkStaking.NewHandler(app.stakingKeeper)).
		AddRoute("distribution", sdkDist.NewHandler(app.distrKeeper)).
		AddRoute("nameservice", nameservice.NewHandler(app.nsKeeper)).
		AddRoute("kyc", kyc.NewHandler(&app.kycKeeper)).
		AddRoute("token", fungible.NewHandler(&app.fungibleTokenKeeper)).
		AddRoute("nonFungible", nonFungible.NewHandler(&app.nonFungibleTokenKeeper)).
		AddRoute("fee", fee.NewHandler(&app.feeKeeper)).
		AddRoute("maintenance", maintenance.NewHandler(&app.maintenanceKeeper, &app.accountKeeper))

	app.QueryRouter().
		AddRoute(sdkAuth.QuerierRoute, sdkAuth.NewQuerier(app.accountKeeper)).
		AddRoute("bank", bank.NewQuerier(app.feeKeeper)).
		AddRoute("staking", sdkStaking.NewQuerier(app.stakingKeeper)).
		AddRoute("distribution", sdkDist.NewQuerier(app.distrKeeper)).
		AddRoute("nameservice", nameservice.NewQuerier(app.cdc, app.nsKeeper, app.feeKeeper)).
		AddRoute("kyc", kyc.NewQuerier(&app.kycKeeper, &app.feeKeeper)).
		AddRoute("token", fungible.NewQuerier(app.cdc, &app.fungibleTokenKeeper, &app.feeKeeper)).
		AddRoute("nonFungible", nonFungible.NewQuerier(app.cdc, &app.nonFungibleTokenKeeper, &app.feeKeeper)).
		AddRoute("fee", fee.NewQuerier(app.cdc, &app.feeKeeper)).
		AddRoute("maintenance", maintenance.NewQuerier(&app.maintenanceKeeper))
		//AddRoute("auth", auth.NewQuerier(app.cdc, app.accountKeeper))

	app.router = app.Router()
	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keySupply,
		app.keyNSnames,
		app.keyNSowners,
		app.keyStaking,
		app.keyParams,
		app.keyDistr,
		app.KeyKyc,
		app.KeyKycData,
		app.tkeyParams,
		app.keyToken,
		app.keyFee,
		app.KeyMaintenance,
		app.KeyValidatorSet,
	)

	if err := app.LoadLatestVersion(app.keyMain); err != nil {
		common.Exit(err.Error())
	}

	rpccore.Routes["query_fee"] = rpc.NewRPCFunc(app.QueryFee, "tx")
	rpccore.Routes["latest_block_height"] = rpc.NewRPCFunc(app.GetLatestBlockHeight, "")
	rpccore.Routes["is_whitelisted"] = rpc.NewRPCFunc(app.CheckWhitelist, "address")
	rpccore.Routes["validator"] = rpc.NewRPCFunc(app.Validator, "address")
	rpccore.Routes["account"] = rpc.NewRPCFunc(app.Account, "address")
	rpccore.Routes["account_cdc"] = rpc.NewRPCFunc(app.AccountCdc, "address")
	rpccore.Routes["decode_tx"] = rpc.NewRPCFunc(app.DecodeTx, "tx")
	rpccore.Routes["encode_tx"] = rpc.NewRPCFunc(app.EncodeTx, "json")
	rpccore.Routes["decoded_tx"] = rpc.NewRPCFunc(app.DecodedTx, "hash,prove")
	rpccore.Routes["encode_and_broadcast_tx_sync"] = rpc.NewRPCFunc(app.EncodeAndBroadcastTxSync, "json")
	rpccore.Routes["encode_and_broadcast_tx_async"] = rpc.NewRPCFunc(app.EncodeAndBroadcastTxAsync, "json")
	rpccore.Routes["encode_and_broadcast_tx_commit"] = rpc.NewRPCFunc(app.EncodeAndBroadcastTxCommit, "json")
	rpccore.Routes["version"] = rpc.NewRPCFunc(app.GetVersion, "")

	// We need to customized it
	rpccore.Routes["debug/fee_info"] = rpc.NewRPCFunc(app.FeeInfo, "")
	rpccore.Routes["debug/kyc_info"] = rpc.NewRPCFunc(app.KYCInfo, "")
	//rpccore.Routes["debug/fungible_token_list"] = rpc.NewRPCFunc(app.FungibleTokenList, "")
	//rpccore.Routes["debug/non_fungible_token_list"] = rpc.NewRPCFunc(app.NonFungibleTokenList, "")

	delete(rpccore.Routes, "genesis")

	return app
}

func (app *mxwApp) initChainer(ctx sdkTypes.Context, req abci.RequestInitChain) abci.ResponseInitChain {

	// The initial state as stored in genesis file
	genesisStateJSON := req.AppStateBytes
	genesisState := new(genesis.GenesisState)

	if err := app.cdc.UnmarshalJSON(genesisStateJSON, genesisState); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal genesis state: %s", err))
	}

	if err := genesisState.Validate(); err != nil {
		panic(err)
	}

	sdkAuth.InitGenesis(ctx, app.accountKeeper, genesisState.AuthState)

	// Setting up initial accounts
	for _, initialAccount := range genesisState.Accounts {
		initialAccount.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, initialAccount)
		app.logger.Info(fmt.Sprintf("Account %v starting with %v coins", initialAccount.Address.String(), initialAccount.Coins))
	}

	sdkBank.InitGenesis(ctx, app.bankKeeper, genesisState.BankState)

	initialValidators := sdkStaking.InitGenesis(ctx, app.stakingKeeper, app.accountKeeper, app.supplyKeeper, genesisState.StakingState)
	for i, validator := range initialValidators {
		app.logger.Info(fmt.Sprintf("Validator %d. Key: %s. Power: %d", i, validator.PubKey.String(), validator.Power))
	}

	sdkDist.InitGenesis(ctx, app.distrKeeper, app.supplyKeeper, genesisState.DistrState)
	kyc.InitGenesis(ctx, &app.kycKeeper, genesisState.KycState)
	fungible.InitGenesis(ctx, &app.fungibleTokenKeeper, genesisState.TokenState)
	nameservice.InitGenesis(ctx, app.nsKeeper, genesisState.NameServiceState)
	fee.InitGenesis(ctx, &app.feeKeeper, genesisState.FeeState)
	maintenance.InitGenesis(ctx, &app.maintenanceKeeper, genesisState.MaintenanceState)

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx sdkAuth.StdTx
			err := app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(abci.RequestDeliverTx{Tx: bz})
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		initialValidators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	initResponse := abci.ResponseInitChain{
		Validators: initialValidators,
	}

	return initResponse
}

func (app *mxwApp) beginBlocker(ctx sdkTypes.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *mxwApp) endBlocker(ctx sdkTypes.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *mxwApp) anteHandler(ctx sdkTypes.Context, tx sdkTypes.Tx, simulate bool) (sdkTypes.Context, error) {
	stdTx, ok := tx.(sdkAuth.StdTx)
	if !ok {
		return ctx, sdkTypes.ErrInternal("Tx must be StdTx.")
	}
	if !app.kycKeeper.CheckTx(ctx, stdTx) {
		return ctx, sdkTypes.NewError("mxw", 1000, "All signers must pass kyc.")
	}

	return app.authAnteHandler(ctx, tx, simulate)
}

func (app *mxwApp) ExportStateAndValidators() (json.RawMessage, []tm.GenesisValidator, error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*sdkAuth.BaseAccount{}

	appendAccountsFn := func(acc exported.Account) bool {
		account := &sdkAuth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)

		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	authState := sdkAuth.ExportGenesis(ctx, app.accountKeeper)
	bankState := sdkBank.ExportGenesis(ctx, app.bankKeeper)
	StakingState := sdkStaking.ExportGenesis(ctx, app.stakingKeeper)
	distrState := sdkDist.ExportGenesis(ctx, app.distrKeeper)
	kycState := kyc.ExportGenesis(&app.kycKeeper)
	tokenState := fungible.ExportGenesis(&app.fungibleTokenKeeper)
	feeState := fee.ExportGenesis(&app.feeKeeper)
	nameServiceState := nameservice.ExportGenesis(&app.nsKeeper)
	maintenanceState := maintenance.ExportGenesis(ctx, &app.maintenanceKeeper)

	appState := genesis.GenesisState{
		AuthState:        authState,
		Accounts:         accounts,
		BankState:        bankState,
		StakingState:     StakingState,
		DistrState:       distrState,
		KycState:         kycState,
		TokenState:       tokenState,
		FeeState:         feeState,
		NameServiceState: nameServiceState,
		MaintenanceState: maintenanceState,
	}

	appStateJSON, err := codec.MarshalJSONIndent(app.cdc, appState)
	if err != nil {
		return nil, nil, err
	}

	validators := sdkStaking.WriteValidators(ctx, app.stakingKeeper)

	return appStateJSON, validators, nil
}

// Uses go-amino which is a fork of protobuf3
// Here the codec implementation is injected into different modules
func MakeDefaultCodec() *codec.Codec {
	var cdc = codec.New()

	// cosmos-sdk using interface to register all the modules codec.
	ModuleBasics.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	nameservice.RegisterCodec(cdc)
	sdkTypes.RegisterCodec(cdc)
	kyc.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	fungible.RegisterCodec(cdc)
	nonFungible.RegisterCodec(cdc)
	fee.RegisterCodec(cdc)
	maintenance.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	return cdc
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *mxwApp) ModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
