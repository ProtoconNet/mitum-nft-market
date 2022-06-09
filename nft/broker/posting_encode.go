package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (posting *Posting) unpack(
	enc encoder.Encoder,
	active bool,
	broker string,
	option string,
	bNFT []byte,
	closetime string,
	bPrice []byte,
) error {
	posting.active = active
	posting.broker = extensioncurrency.ContractID(broker)
	posting.option = PostOption(option)

	if hinter, err := enc.Decode(bNFT); err != nil {
		return err
	} else if nft, ok := hinter.(nft.NFTID); !ok {
		return errors.Errorf("not NFTID; %T", hinter)
	} else {
		posting.nft = nft
	}

	posting.closeTime = PostCloseTime(closetime)

	if hinter, err := enc.Decode(bPrice); err != nil {
		return err
	} else if price, ok := hinter.(currency.Amount); !ok {
		return errors.Errorf("not Amount; %T", hinter)
	} else {
		posting.price = price
	}

	return nil
}

func (bid *Bidding) unpack(
	enc encoder.Encoder,
	bBidder base.AddressDecoder,
	bAmount []byte,
) error {
	bidder, err := bBidder.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(bAmount); err != nil {
		return err
	} else if amount, ok := hinter.(currency.Amount); !ok {
		return errors.Errorf("not Amount; %T", hinter)
	} else {
		bid.amount = amount
	}

	bid.bidder = bidder

	return nil
}
