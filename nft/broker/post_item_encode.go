package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (form *PostForm) unpack(
	enc encoder.Encoder,
	option string,
	bn []byte,
	closetime string,
	bp []byte,
) error {
	form.option = PostOption(option)

	if hinter, err := enc.Decode(bn); err != nil {
		return err
	} else if n, ok := hinter.(nft.NFTID); !ok {
		return util.WrongTypeError.Errorf("not NFTID; %T", hinter)
	} else {
		form.n = n
	}

	form.closeTime = PostCloseTime(closetime)

	if hinter, err := enc.Decode(bp); err != nil {
		return err
	} else if am, ok := hinter.(currency.Amount); !ok {
		return util.WrongTypeError.Errorf("not Amount; %T", hinter)
	} else {
		form.price = am
	}

	return nil
}

func (it *BasePostItem) unpack(
	enc encoder.Encoder,
	broker string,
	bForms []byte,
	cid string,
) error {
	it.broker = extensioncurrency.ContractID(broker)

	hForms, err := enc.DecodeSlice(bForms)
	if err != nil {
		return err
	}

	forms := make([]PostForm, len(hForms))
	for i := range hForms {
		j, ok := hForms[i].(PostForm)
		if !ok {
			return util.WrongTypeError.Errorf("not PostForm; %T", hForms[i])
		}
		forms[i] = j
	}
	it.forms = forms

	it.cid = currency.CurrencyID(cid)

	return nil
}
