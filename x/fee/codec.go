package fee

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSysFeeSetting{}, "fee/sysFeeSetting", nil)
	cdc.RegisterConcrete(MsgAssignFeeToMsg{}, "fee/assignFeeToMsg", nil)
	cdc.RegisterConcrete(MsgAssignFeeToAcc{}, "fee/assignFeeToAcc", nil)
	cdc.RegisterConcrete(MsgAssignFeeToToken{}, "fee/assignFeeToToken", nil)
	cdc.RegisterConcrete(MsgMultiplier{}, "fee/msgMultiplier", nil)
	cdc.RegisterConcrete(MsgTokenMultiplier{}, "fee/msgTokenMultiplier", nil)
	cdc.RegisterConcrete(MsgDeleteSysFeeSetting{}, "fee/deleteSysFeeSetting", nil)
	cdc.RegisterConcrete(MsgDeleteAccFeeSetting{}, "fee/deleteAccFeeSetting", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
