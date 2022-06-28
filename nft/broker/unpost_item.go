package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	UnpostItemType   = hint.Type("mitum-nft-market-unpost-item")
	UnpostItemHint   = hint.NewHint(UnpostItemType, "v0.0.1")
	UnpostItemHinter = UnpostItem{
		BaseHinter: hint.NewBaseHinter(UnpostItemHint),
	}
)

type UnpostItem struct {
	hint.BaseHinter
	nft nft.NFTID
	cid currency.CurrencyID
}

func NewUnpostItem(n nft.NFTID, cid currency.CurrencyID) UnpostItem {
	return UnpostItem{
		BaseHinter: hint.NewBaseHinter(UnpostItemHint),
		nft:        n,
		cid:        cid,
	}
}

func (it UnpostItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.nft.Bytes(),
		it.cid.Bytes(),
	)
}

func (it UnpostItem) IsValid([]byte) error {
	return isvalid.Check(nil, false,
		it.BaseHinter,
		it.nft,
		it.cid)
}

func (it UnpostItem) NFT() nft.NFTID {
	return it.nft
}

func (it UnpostItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it UnpostItem) Rebuild() UnpostItem {
	return it
}
