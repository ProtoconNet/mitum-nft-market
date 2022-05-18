package broker

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact *BrokerRegisterFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bTarget base.AddressDecoder,
	bPolicy []byte,
	cid string,
) error {
	sender, err := bSender.Encode(enc)
	if err != nil {
		return err
	}

	target, err := bTarget.Encode(enc)
	if err != nil {
		return err
	}

	var policy BrokerPolicy
	if hinter, err := enc.Decode(bPolicy); err != nil {
		return err
	} else if i, ok := hinter.(BrokerPolicy); !ok {
		return util.WrongTypeError.Errorf("not BrokerPolicy; %T", hinter)
	} else {
		policy = i
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.target = target
	fact.policy = policy
	fact.cid = currency.CurrencyID(cid)

	return nil
}
