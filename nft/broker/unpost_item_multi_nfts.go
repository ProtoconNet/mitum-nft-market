package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var MaxNFTsUnpostItemMultiNFTs = 10

var (
	UnpostItemMultiNFTsType   = hint.Type("mitum-nft-market-unpost-multi-nfts")
	UnpostItemMultiNFTsHint   = hint.NewHint(UnpostItemMultiNFTsType, "v0.0.1")
	UnpostItemMultiNFTsHinter = UnpostItemMultiNFTs{
		BaseUnpostItem: BaseUnpostItem{
			BaseHinter: hint.NewBaseHinter(UnpostItemMultiNFTsHint),
		},
	}
)

type UnpostItemMultiNFTs struct {
	BaseUnpostItem
}

func NewUnpostItemMultiNFTs(nfts []nft.NFTID, cid currency.CurrencyID) UnpostItemMultiNFTs {
	return UnpostItemMultiNFTs{
		BaseUnpostItem: NewBaseUnpostItem(UnpostItemMultiNFTsHint, nfts, cid),
	}
}

func (it UnpostItemMultiNFTs) IsValid([]byte) error {
	if err := it.BaseUnpostItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.nfts); n > MaxNFTsUnpostItemMultiNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", n, MaxNFTsUnpostItemMultiNFTs)
	}

	return nil
}

func (it UnpostItemMultiNFTs) Rebuild() UnpostItem {
	it.BaseUnpostItem = it.BaseUnpostItem.Rebuild().(BaseUnpostItem)

	return it
}
