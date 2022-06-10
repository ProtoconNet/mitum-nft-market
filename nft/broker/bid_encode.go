package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *BidFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bNFT []byte,
	bAmount []byte,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(bNFT); err != nil {
		return err
	} else if nft, ok := hinter.(nft.NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		fact.nft = nft
	}

	if hinter, err := enc.Decode(bAmount); err != nil {
		return err
	} else if amount, ok := hinter.(currency.Amount); !ok {
		return util.WrongTypeError.Errorf("not Amount; %T", hinter)
	} else {
		fact.amount = amount
	}

	fact.h = h
	fact.token = token
	fact.sender = sender

	return nil
}
