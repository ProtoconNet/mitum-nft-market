package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseUnpostItem struct {
	hint.BaseHinter
	nfts []nft.NFTID
	cid  currency.CurrencyID
}

func NewBaseUnpostItem(ht hint.Hint, nfts []nft.NFTID, cid currency.CurrencyID) BaseUnpostItem {
	return BaseUnpostItem{
		BaseHinter: hint.NewBaseHinter(ht),
		nfts:       nfts,
		cid:        cid,
	}
}

func (it BaseUnpostItem) Bytes() []byte {
	ns := make([][]byte, len(it.nfts))

	for i := range it.nfts {
		ns[i] = it.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (it BaseUnpostItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, it.BaseHinter, it.cid); err != nil {
		return err
	}

	if len(it.nfts) < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for BaseUnpostItem")
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range it.nfts {
		if err := it.nfts[i].IsValid(nil); err != nil {
			return err
		}
		n := it.nfts[i]
		if _, found := foundNFT[n]; found {
			return errors.Errorf("duplicated nft found; %s", n)
		}
		foundNFT[n] = true
	}

	return nil
}

func (it BaseUnpostItem) NFTs() []nft.NFTID {
	return it.nfts
}

func (it BaseUnpostItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseUnpostItem) Rebuild() UnpostItem {
	nfts := make([]nft.NFTID, len(it.nfts))
	for i := range it.nfts {
		nfts[i] = it.nfts[i]
	}
	it.nfts = nfts

	return it
}
