package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
	ver "github.com/maxonrow/maxonrow-go/version"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/token/fungible"
	"github.com/maxonrow/maxonrow-go/x/token/nonfungible"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	tmTypes "github.com/tendermint/tendermint/types"
	tmver "github.com/tendermint/tendermint/version"
)

// Result of querying for a tx
type ResultDecodedTx struct {
	Hash     cmn.HexBytes           `json:"hash"`
	Height   int64                  `json:"height"`
	Index    uint32                 `json:"index"`
	TxResult abci.ResponseDeliverTx `json:"tx_result"`
	Tx       string                 `json:"tx"`
	Proof    tmTypes.TxProof        `json:"proof,omitempty"`
}

type Version struct {
	MaxonrowVersion   string `json:"Maxonrow"`
	TendermintVersion string `json:"Tendermint"`
}

type FungibleTokenListInfo struct {
	Count  int
	Tokens []fungible.Token
}

type NonFungibleTokenListInfo struct {
	Count  int
	Tokens []nonfungible.Token
}

type FeeInfo struct {
	Authorizers        []sdkTypes.AccAddress
	Multiplier         string
	TokenMultiplier    string
	TokenFeeCollectors []sdkTypes.AccAddress
	AliasFeeCollectors []sdkTypes.AccAddress
	FeeSettings        []fee.FeeSetting
	AccountFeeSettings map[string]string
	MsgFeeSettings     map[string]string
	TokenFeeSetting    map[string]string
}

type KYCInfo struct {
	Providers        []sdkTypes.AccAddress
	Issuers          []sdkTypes.AccAddress
	Authorizers      []sdkTypes.AccAddress
	NumOfWhitelisted int
}

func (app *mxwApp) DecodeTx(ctx *rpctypes.Context, bz []byte) (string, error) {
	tx, err := app.txDecoder(bz)
	if err != nil {
		return "", err
	}
	js, err1 := app.cdc.MarshalJSON(tx)
	if err != nil {
		return "", err1
	}
	return string(js), nil
}

func (app *mxwApp) EncodeTx(ctx *rpctypes.Context, js string) ([]byte, error) {
	bz := parseJSON(js)
	var tx sdkAuth.StdTx
	err := app.cdc.UnmarshalJSON(bz, &tx)
	if err != nil {
		return nil, err
	}
	return app.txEncoder(tx)
}

func (app *mxwApp) DecodedTx(ctx *rpctypes.Context, hash []byte, prove bool) (*ResultDecodedTx, error) {
	res, err := rpc.Tx(ctx, hash, prove)
	if err != nil {
		return nil, err
	}

	tx, err := app.txDecoder(res.Tx)
	if err != nil {
		return nil, err
	}

	out, err := app.cdc.MarshalJSON(tx)
	if err != nil {
		return nil, err
	}

	out = sdkTypes.MustSortJSON(out)

	return &ResultDecodedTx{
		Hash:     res.Hash,
		Height:   res.Height,
		Index:    res.Index,
		TxResult: res.TxResult,
		Tx:       string(out),
		Proof:    res.Proof,
	}, nil
}

func (app *mxwApp) EncodeAndBroadcastTxSync(ctx *rpctypes.Context, js string) (*ctypes.ResultBroadcastTx, error) {
	bz := parseJSON(js)
	var tx sdkAuth.StdTx
	err := app.cdc.UnmarshalJSON(bz, &tx)
	if err != nil {
		return nil, err
	}

	txByte, err := app.txEncoder(tx)
	if err != nil {
		return nil, err
	}
	return rpc.BroadcastTxSync(ctx, txByte)
}

func (app *mxwApp) EncodeAndBroadcastTxAsync(ctx *rpctypes.Context, js string) (*ctypes.ResultBroadcastTx, error) {
	bz := parseJSON(js)
	var tx sdkAuth.StdTx
	err := app.cdc.UnmarshalJSON(bz, &tx)
	if err != nil {
		return nil, err
	}

	txByte, err := app.txEncoder(tx)
	if err != nil {
		return nil, err
	}

	return rpc.BroadcastTxAsync(ctx, txByte)
}

func (app *mxwApp) EncodeAndBroadcastTxCommit(ctx *rpctypes.Context, js string) (*ctypes.ResultBroadcastTxCommit, error) {
	bz := parseJSON(js)
	var tx sdkAuth.StdTx
	err := app.cdc.UnmarshalJSON(bz, &tx)
	if err != nil {
		return nil, err
	}

	txByte, err := app.txEncoder(tx)
	if err != nil {
		return nil, err
	}

	return rpc.BroadcastTxCommit(ctx, txByte)
}

func (app *mxwApp) Account(ctx *rpctypes.Context, str string) (string, error) {
	addr, err := sdkTypes.AccAddressFromBech32(str)
	if err != nil {
		return "", err
	}
	appCtx := app.NewContext(true, abci.Header{})
	acc := app.accountKeeper.GetAccount(appCtx, addr)

	out, err := app.cdc.MarshalJSON(acc)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (app *mxwApp) Validator(ctx *rpctypes.Context, str string) (string, error) {
	addr, err := sdkTypes.ValAddressFromBech32(str)
	if err != nil {
		return "", err
	}
	appCtx := app.NewContext(true, abci.Header{})
	val, exist := app.stakingKeeper.GetValidator(appCtx, addr)

	if !exist {
		return "", fmt.Errorf("Not exists")
	}
	out, err := app.cdc.MarshalJSON(val)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (app *mxwApp) KYCInfo(ctx *rpctypes.Context) (KYCInfo, error) {
	appCtx := app.NewContext(true, abci.Header{})

	var i KYCInfo
	i.Providers = app.kycKeeper.GetProviderAddresses(appCtx)
	i.Issuers = app.kycKeeper.GetIssuerAddresses(appCtx)
	i.Authorizers = app.kycKeeper.GetAuthorisedAddresses(appCtx)
	i.NumOfWhitelisted = app.kycKeeper.NumOfWhitelisted(appCtx)
	return i, nil
}

func (app *mxwApp) FeeInfo(ctx *rpctypes.Context) (FeeInfo, error) {
	appCtx := app.NewContext(true, abci.Header{})

	var i FeeInfo
	i.Authorizers = app.feeKeeper.GetAuthorisedAddresses(appCtx)
	i.Multiplier, _ = app.feeKeeper.GetFeeMultiplier(appCtx)
	i.TokenMultiplier, _ = app.feeKeeper.GetTokenFeeMultiplier(appCtx)
	i.TokenFeeCollectors = app.feeKeeper.GetFeeCollectorAddresses(appCtx, "token")
	i.AliasFeeCollectors = app.feeKeeper.GetFeeCollectorAddresses(appCtx, "alias")
	i.FeeSettings = app.feeKeeper.ListAllSysFeeSetting(appCtx)
	i.AccountFeeSettings = app.feeKeeper.ListAllAccountFeeSettings(appCtx)
	i.MsgFeeSettings = app.feeKeeper.ListAllMsgFeeSettings(appCtx)
	i.TokenFeeSetting = app.feeKeeper.ListAllMsgFeeSettings(appCtx)
	return i, nil
}

func (app *mxwApp) QueryFee(ctx *rpctypes.Context, js string) (sdkAuth.StdFee, error) {
	var fees sdkAuth.StdFee
	appCtx := app.NewContext(true, abci.Header{})
	bz := parseJSON(js)
	var tx sdkAuth.StdTx
	err := app.cdc.UnmarshalJSON(bz, &tx)
	if err != nil {
		return sdkAuth.StdFee{}, err
	}
	fee, feeErr := app.CalculateFee(appCtx, tx)
	if feeErr != nil {
		return sdkAuth.StdFee{}, feeErr
	}
	// When the fee is empty, return zero
	if fee.Empty() {
		zero := sdkTypes.Coin{Amount: sdkTypes.NewInt(0), Denom: types.CIN}
		fee = sdkTypes.Coins{zero}
	}

	fees.Amount = fee
	fees.Gas = 0

	return fees, nil
}

func (app *mxwApp) GetLatestBlockHeight(ctx *rpctypes.Context) (int64, error) {

	blockResult, err := rpc.Status(ctx)
	if err != nil {
		return 0, err
	}
	return blockResult.SyncInfo.LatestBlockHeight, nil
}

func (app *mxwApp) CheckWhitelist(ctx *rpctypes.Context, address string) (bool, error) {
	addr, addrErr := sdkTypes.AccAddressFromBech32(address)
	if addrErr != nil {
		return false, addrErr
	}
	appCtx := app.NewContext(true, abci.Header{})
	return app.kycKeeper.IsWhitelisted(appCtx, addr), nil
}

func (app *mxwApp) GetVersion(ctx *rpctypes.Context) (Version, error) {

	ver := Version{TendermintVersion: tmver.Version, MaxonrowVersion: ver.GetVersion()}
	return ver, nil
}

func parseJSON(in string) []byte {
	if json.Valid([]byte(in)) {
		return []byte(in)
	}

	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic("Not a valid json input")
	}
	return []byte(out)
}

func (app *mxwApp) FungibleTokenList(ctx *rpctypes.Context) (FungibleTokenListInfo, error) {
	appCtx := app.NewContext(true, abci.Header{})
	lst := app.tokenKeeper.ListTokens(appCtx)
	var i FungibleTokenListInfo
	i.Count = len(lst)
	i.Tokens = lst
	return i, nil
}

func (app *mxwApp) NonFungibleTokenList(ctx *rpctypes.Context) (NonFungibleTokenListInfo, error) {
	appCtx := app.NewContext(true, abci.Header{})
	lst := app.nonFungibleTokenKeeper.ListTokens(appCtx)
	var i NonFungibleTokenListInfo
	i.Count = len(lst)
	i.Tokens = lst
	return i, nil
}
