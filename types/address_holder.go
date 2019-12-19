package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type AddressHolder []sdkTypes.AccAddress

func (a *AddressHolder) AppendAccAddrs(addrs []sdkTypes.AccAddress) {
	for _, addr := range addrs {
		a.Append(addr)
	}
}

func (a *AddressHolder) Append(addr sdkTypes.AccAddress) bool {
	if _, exist := a.Contains(addr); exist {
		return false
	}
	*a = append(*a, addr)
	return true
}

func (a *AddressHolder) Contains(addr sdkTypes.AccAddress) (int, bool) {
	for i, e := range *a {
		if e.Equals(addr) {
			return i, true
		}
	}
	return -1, false
}

func (a *AddressHolder) Remove(addr sdkTypes.AccAddress) bool {
	index, exists := a.Contains(addr)
	if exists {
		s := *a
		*a = append(s[:index], s[index+1:]...)
		return true
	}
	return false
}

func (a *AddressHolder) Size() int {
	return len(*a)
}
