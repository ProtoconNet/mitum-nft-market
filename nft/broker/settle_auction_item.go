package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseSettleAuctionItem struct {
	hint.BaseHinter
	nfts []nft.NFTID
	cid  currency.CurrencyID
}

func NewBaseSettleAuctionItem(ht hint.Hint, nfts []nft.NFTID, cid currency.CurrencyID) BaseSettleAuctionItem {
	return BaseSettleAuctionItem{
		BaseHinter: hint.NewBaseHinter(ht),
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseSettleAuctionItem) Bytes() []byte {
	ns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		ns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (it BaseSettleAuctionItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseSettleAuctionItem")
	}

	foundNFT := map[string]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}
		nft := it.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return errors.Errorf("duplicated nft found; %s", nft)
		}
		foundNFT[nft] = true
	}

	return nil
}

func (it BaseSettleAuctionItem) NFTs() []nft.NFTID {
	return it.nfts
}

func (it BaseSettleAuctionItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseSettleAuctionItem) Rebuild() SettleAuctionItem {
	nfts := make([]nft.NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
