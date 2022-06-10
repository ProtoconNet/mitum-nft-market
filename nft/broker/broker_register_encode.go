package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form *BrokerRegisterForm) unpack(
	enc encoder.Encoder,
	bt base.AddressDecoder,
	symbol string,
	brokerage uint,
	br base.AddressDecoder,
	royalty bool,
	uri string,
) error {
	target, err := bt.Encode(enc)
	if err != nil {
		return err
	}
	form.target = target

	receiver, err := br.Encode(enc)
	if err != nil {
		return err
	}
	form.receiver = receiver

	form.symbol = extensioncurrency.ContractID(symbol)
	form.brokerage = nft.PaymentParameter(brokerage)
	form.royalty = royalty
	form.uri = nft.URI(uri)

	return nil
}

func (fact *BrokerRegisterFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bs base.AddressDecoder,
	bf []byte,
	cid string,
) error {
	sender, err := bs.Encode(enc)
	if err != nil {
		return err
	}

	var form BrokerRegisterForm
	if hinter, err := enc.Decode(bf); err != nil {
		return err
	} else if i, ok := hinter.(BrokerRegisterForm); !ok {
		return util.WrongTypeError.Errorf("not BrokerRegisterForm; %T", hinter)
	} else {
		form = i
	}

	fact.h = h
	fact.token = token
	fact.sender = sender
	fact.form = form
	fact.cid = currency.CurrencyID(cid)

	return nil
}
