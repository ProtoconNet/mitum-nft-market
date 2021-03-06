package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (details *SellDetails) unpack(
	enc encoder.Encoder,
	bn []byte,
	bp []byte,
) error {

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return errors.Errorf("not NFTID; %T", hinter)
	} else {
		details.nft = n
	}

	if hinter, err := enc.Decode(bp); err != nil {
		return err
	} else if am, ok := hinter.(currency.Amount); !ok {
		return errors.Errorf("not Amount; %T", hinter)
	} else {
		details.price = am
	}

	return nil
}

func (details *AuctionDetails) unpack(
	enc encoder.Encoder,
	bn []byte,
	closeTime string,
	bp []byte,
) error {

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return errors.Errorf("not NFTID; %T", hinter)
	} else {
		details.nft = n
	}

	details.closeTime = PostCloseTime(closeTime)

	if hinter, err := enc.Decode(bp); err != nil {
		return err
	} else if am, ok := hinter.(currency.Amount); !ok {
		return errors.Errorf("not Amount; %T", hinter)
	} else {
		details.price = am
	}

	return nil
}

func (posting *Posting) unpack(
	enc encoder.Encoder,
	active bool,
	broker string,
	option string,
	bd []byte,
) error {
	posting.active = active
	posting.broker = extensioncurrency.ContractID(broker)
	posting.option = PostOption(option)

	if hinter, err := enc.Decode(bd); err != nil {
		return err
	} else if details, ok := hinter.(PostDetails); !ok {
		return errors.Errorf("not PostDetails; %T", hinter)
	} else {
		posting.details = details
	}

	return nil
}

func (bid *Bidding) unpack(
	enc encoder.Encoder,
	bb base.AddressDecoder,
	ba []byte,
) error {
	bidder, err := bb.Encode(enc)
	if err != nil {
		return err
	}

	if hinter, err := enc.Decode(ba); err != nil {
		return err
	} else if am, ok := hinter.(currency.Amount); !ok {
		return errors.Errorf("not Amount; %T", hinter)
	} else {
		bid.amount = am
	}

	bid.bidder = bidder

	return nil
}
