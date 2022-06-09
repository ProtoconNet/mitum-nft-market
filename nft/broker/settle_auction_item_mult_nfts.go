package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsSettleAuctionItemMultiNFTs = 10

var (
	SettleAuctionItemMultiNFTsType   = hint.Type("mitum-nft-market-settle-auction-multi-nfts")
	SettleAuctionItemMultiNFTsHint   = hint.NewHint(SettleAuctionItemMultiNFTsType, "v0.0.1")
	SettleAuctionItemMultiNFTsHinter = SettleAuctionItemMultiNFTs{
		BaseSettleAuctionItem: BaseSettleAuctionItem{
			BaseHinter: hint.NewBaseHinter(SettleAuctionItemMultiNFTsHint),
		},
	}
)

type SettleAuctionItemMultiNFTs struct {
	BaseSettleAuctionItem
}

func NewSettleAuctionItemMultiNFTs(nfts []nft.NFTID, cid currency.CurrencyID) SettleAuctionItemMultiNFTs {
	return SettleAuctionItemMultiNFTs{
		BaseSettleAuctionItem: NewBaseSettleAuctionItem(SettleAuctionItemMultiNFTsHint, nfts, cid),
	}
}

func (it SettleAuctionItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseSettleAuctionItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n > MaxNFTsSettleAuctionItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", n, MaxNFTsSettleAuctionItemMultiNFTs)
	}

	return nil
}

func (it SettleAuctionItemMultiNFTs) Rebuild() SettleAuctionItem {
	it.BaseSettleAuctionItem = it.BaseSettleAuctionItem.Rebuild().(BaseSettleAuctionItem)

	return it
}
