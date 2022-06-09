package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	UnpostItemSingleNFTType   = hint.Type("mitum-nft-market-unpost-single-nft")
	UnpostItemSingleNFTHint   = hint.NewHint(UnpostItemSingleNFTType, "v0.0.1")
	UnpostItemSingleNFTHinter = UnpostItemSingleNFT{
		BaseUnpostItem: BaseUnpostItem{
			BaseHinter: hint.NewBaseHinter(UnpostItemSingleNFTHint),
		},
	}
)

type UnpostItemSingleNFT struct {
	BaseUnpostItem
}

func NewUnpostItemSingleNFT(nftid nft.NFTID, cid currency.CurrencyID) UnpostItemSingleNFT {
	return UnpostItemSingleNFT{
		BaseUnpostItem: NewBaseUnpostItem(UnpostItemSingleNFTHint, []nft.NFTID{nftid}, cid),
	}
}

func (it UnpostItemSingleNFT) IsValid([]byte) error {
	if err := it.BaseUnpostItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it UnpostItemSingleNFT) Rebuild() UnpostItem {
	it.BaseUnpostItem = it.BaseUnpostItem.Rebuild().(BaseUnpostItem)

	return it
}
