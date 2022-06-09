package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	PostItemSingleNFTType   = hint.Type("mitum-nft-market-post-single-nft")
	PostItemSingleNFTHint   = hint.NewHint(PostItemSingleNFTType, "v0.0.1")
	PostItemSingleNFTHinter = PostItemSingleNFT{
		BasePostItem: BasePostItem{
			BaseHinter: hint.NewBaseHinter(PostItemSingleNFTHint),
		},
	}
)

type PostItemSingleNFT struct {
	BasePostItem
}

func NewPostItemSingleNFT(broker extensioncurrency.ContractID, form PostForm, cid currency.CurrencyID) PostItemSingleNFT {
	return PostItemSingleNFT{
		BasePostItem: NewBasePostItem(PostItemSingleNFTHint, broker, []PostForm{form}, cid),
	}
}

func (it PostItemSingleNFT) IsValid([]byte) error {
	if err := it.BasePostItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.forms); n != 1 {
		return isvalid.InvalidError.Errorf("only one nft allowed; %d", n)
	}

	return nil
}

func (it PostItemSingleNFT) Rebuild() PostItem {
	it.BasePostItem = it.BasePostItem.Rebuild().(BasePostItem)

	return it
}
