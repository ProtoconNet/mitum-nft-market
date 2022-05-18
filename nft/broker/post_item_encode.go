package broker

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *PostItem) unpack(enc encoder.Encoder, bPosting []byte, cid string) error {

	if hinter, err := enc.Decode(bPosting); err != nil {
		return err
	} else if posting, ok := hinter.(Posting); !ok {
		return util.WrongTypeError.Errorf("not Posting; %T", hinter)
	} else {
		it.posting = posting
	}

	it.cid = currency.CurrencyID(cid)

	return nil
}
