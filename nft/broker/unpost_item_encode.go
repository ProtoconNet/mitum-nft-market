package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *UnpostItem) unpack(
	enc encoder.Encoder,
	bn []byte,
	cid string,
) error {
	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return errors.Errorf("not NFTID; %T", hinter)
	} else {
		it.nft = n
	}

	it.cid = currency.CurrencyID(cid)

	return nil
}
