package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (bp *BrokerPolicy) unpack(
	enc encoder.Encoder,
	brokerage uint,
	bReceiver base.AddressDecoder,
	royalty bool,
) error {
	bp.brokerage = nft.PaymentParameter(brokerage)

	receiver, err := bReceiver.Encode(enc)
	if err != nil {
		return err
	}
	bp.receiver = receiver

	bp.royalty = royalty

	return nil
}
