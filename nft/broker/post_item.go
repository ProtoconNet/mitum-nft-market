package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	PostItemType   = hint.Type("mitum-nft-post-item")
	PostItemHint   = hint.NewHint(PostItemType, "v0.0.1")
	PostItemHinter = PostItem{BaseHinter: hint.NewBaseHinter(PostItemHint)}
)

type PostItem struct {
	hint.BaseHinter
	posting Posting
	cid     currency.CurrencyID
}

func NewPostItem(posting Posting, cid currency.CurrencyID) PostItem {
	return PostItem{
		BaseHinter: hint.NewBaseHinter(PostItemHint),
		posting:    posting,
		cid:        cid,
	}
}

func (it PostItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.posting.Bytes(),
		it.cid.Bytes(),
	)
}

func (it PostItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		it.BaseHinter,
		it.posting,
		it.cid); err != nil {
		return err
	}

	return nil
}

func (it PostItem) NFT() nft.NFTID {
	return it.posting.NFT()
}

func (it PostItem) Posting() Posting {
	return it.posting
}

func (it PostItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it PostItem) Rebuild() PostItem {
	posting := it.posting.Rebuild()
	it.posting = posting

	return it
}
