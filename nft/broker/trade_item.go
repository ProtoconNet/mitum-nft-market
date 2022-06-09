package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	TradeItemType   = hint.Type("mitum-nft-market-trade-item")
	TradeItemHint   = hint.NewHint(TradeItemType, "v0.0.1")
	TradeItemHinter = TradeItem{BaseHinter: hint.NewBaseHinter(TradeItemHint)}
)

type TradeItem struct {
	hint.BaseHinter
	n nft.NFTID
}

func NewTradeItem(n nft.NFTID) TradeItem {
	return TradeItem{
		BaseHinter: hint.NewBaseHinter(TradeItemHint),
		n:          n,
	}
}

func (it TradeItem) Bytes() []byte {
	return it.n.Bytes()
}

func (it TradeItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.n); err != nil {
		return err
	}

	return nil
}

func (it TradeItem) NFT() nft.NFTID {
	return it.n
}

func (it TradeItem) Rebuild() TradeItem {
	return it
}
