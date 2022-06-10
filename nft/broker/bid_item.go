package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	BidItemType   = hint.Type("mitum-nft-market-bid-item")
	BidItemHint   = hint.NewHint(BidItemType, "v0.0.1")
	BidItemHinter = BidItem{BaseHinter: hint.NewBaseHinter(BidItemHint)}
)

type BidItem struct {
	hint.BaseHinter
	n      nft.NFTID
	amount currency.Amount
}

func NewBidItem(n nft.NFTID, amount currency.Amount) BidItem {
	return BidItem{
		BaseHinter: hint.NewBaseHinter(BidItemHint),
		n:          n,
		amount:     amount,
	}
}

func MustNewBidItem(n nft.NFTID, amount currency.Amount) BidItem {
	item := NewBidItem(n, amount)
	if err := item.IsValid(nil); err != nil {
		panic(err)
	}
	return item
}

func (it BidItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.n.Bytes(),
		it.amount.Bytes(),
	)
}

func (it BidItem) IsValid([]byte) error {
	if !it.amount.Big().OverZero() {
		return isvalid.InvalidError.Errorf("bid must be greater than zero")
	}

	if err := isvalid.Check(nil, false, it.BaseHinter, it.n, it.amount); err != nil {
		return err
	}

	return nil
}

func (it BidItem) NFT() nft.NFTID {
	return it.n
}

func (it BidItem) Amount() currency.Amount {
	return it.amount
}

func (it BidItem) Rebuild() BidItem {
	return it
}
