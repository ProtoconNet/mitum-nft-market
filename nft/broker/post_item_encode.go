package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (form *PostForm) unpack(
	enc encoder.Encoder,
	option string,
	bd []byte,
) error {
	form.option = PostOption(option)

	if hinter, err := enc.Decode(bd); err != nil {
		return err
	} else if details, ok := hinter.(PostDetails); !ok {
		return util.WrongTypeError.Errorf("not PostDetails; %T", hinter)
	} else {
		form.details = details
	}

	return nil
}

func (it *PostItem) unpack(
	enc encoder.Encoder,
	broker string,
	bf []byte,
	cid string,
) error {
	it.broker = extensioncurrency.ContractID(broker)

	if hinter, err := enc.Decode(bf); err != nil {
		return err
	} else if form, ok := hinter.(PostForm); !ok {
		return errors.Errorf("not PostForm; %T", hinter)
	} else {
		it.form = form
	}

	it.cid = currency.CurrencyID(cid)

	return nil
}
