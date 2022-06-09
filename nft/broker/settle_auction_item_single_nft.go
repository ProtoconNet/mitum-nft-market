package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	SettleAuctionItemSingleNFTType   = hint.Type("mitum-nft-market-settle-auction-single-nft")
	SettleAuctionItemSingleNFTHint   = hint.NewHint(SettleAuctionItemSingleNFTType, "v0.0.1")
	SettleAuctionItemSingleNFTHinter = SettleAuctionItemSingleNFT{
		BaseSettleAuctionItem: BaseSettleAuctionItem{
			BaseHinter: hint.NewBaseHinter(SettleAuctionItemSingleNFTHint),
		},
	}
)

type SettleAuctionItemSingleNFT struct {
	BaseSettleAuctionItem
}

func NewSettleAuctionItemSingleNFT(nftid nft.NFTID, cid currency.CurrencyID) SettleAuctionItemSingleNFT {
	return SettleAuctionItemSingleNFT{
		BaseSettleAuctionItem: NewBaseSettleAuctionItem(SettleAuctionItemSingleNFTHint, []nft.NFTID{nftid}, cid),
	}
}

func (it SettleAuctionItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseSettleAuctionItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it SettleAuctionItemSingleNFT) Rebuild() SettleAuctionItem {
	it.BaseSettleAuctionItem = it.BaseSettleAuctionItem.Rebuild().(BaseSettleAuctionItem)

	return it
}
