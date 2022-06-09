package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *TradeItem) unpack(
	enc encoder.Encoder,
	bn []byte,
) error {

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		it.n = n
	}

	return nil
}
