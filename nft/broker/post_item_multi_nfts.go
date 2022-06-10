package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsPostItemMultiNFTs = 10

var (
	PostItemMultiNFTsType   = hint.Type("mitum-nft-post-multi-nfts")
	PostItemMultiNFTsHint   = hint.NewHint(PostItemMultiNFTsType, "v0.0.1")
	PostItemMultiNFTsHinter = PostItemMultiNFTs{
		BasePostItem: BasePostItem{
			BaseHinter: hint.NewBaseHinter(PostItemMultiNFTsHint),
		},
	}
)

type PostItemMultiNFTs struct {
	BasePostItem
}

func NewPostItemMultiNFTs(broker extensioncurrency.ContractID, forms []PostForm, cid currency.CurrencyID) PostItemMultiNFTs {
	return PostItemMultiNFTs{
		BasePostItem: NewBasePostItem(PostItemMultiNFTsHint, broker, forms, cid),
	}
}

func (it PostItemMultiNFTs) IsValid([]byte) error {
	if err := it.BasePostItem.IsValid(nil); err != nil {
		return err
	}

	if l := len(it.forms); l > MaxNFTsPostItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", l, MaxNFTsPostItemMultiNFTs)
	}

	return nil
}

func (it PostItemMultiNFTs) Rebuild() PostItem {
	it.BasePostItem = it.BasePostItem.Rebuild().(BasePostItem)

	return it
}
