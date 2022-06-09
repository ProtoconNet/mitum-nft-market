package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BaseUnpostItem) unpack(
	enc encoder.Encoder,
	bns []byte,
	cid string,
) error {
	hns, err := enc.DecodeSlice(bns)
	if err != nil {
		return err
	}

	nfts := make([]nft.NFTID, len(hns))
	for i := range hns {
		j, ok := hns[i].(nft.NFTID)
		if !ok {
			return util.WrongTypeError.Errorf("not NFTID; %T", hns[i])
		}

		nfts[i] = j
	}

	it.nfts = nfts
	it.cid = currency.CurrencyID(cid)

	return nil
}
