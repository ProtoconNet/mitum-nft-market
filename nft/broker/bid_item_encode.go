package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BidItem) unpack(
	enc encoder.Encoder,
	bn []byte,
	ba []byte,
) error {

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		it.n = n
	}

	if hinter, err := enc.Decode(ba); err != nil {
		return err
	} else if am, ok := hinter.(currency.Amount); !ok {
		return util.WrongTypeError.Errorf("not Amount; %T", hinter)
	} else {
		it.amount = am
	}

	return nil
}
