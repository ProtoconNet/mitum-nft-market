package broker

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PostFormJSONPacker struct {
	jsonenc.HintedHead
	OP PostOption  `json:"option"`
	DE PostDetails `json:"details"`
}

func (form PostForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		OP:         form.option,
		DE:         form.details,
	})
}

type PostFormJSONUnpacker struct {
	OP string          `json:"option"`
	DE json.RawMessage `json:"details"`
}

func (form *PostForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufo PostFormJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.OP, ufo.DE)
}

type PostItemJSONPacker struct {
	jsonenc.HintedHead
	BR extensioncurrency.ContractID `json:"broker"`
	FO PostForm                     `json:"form"`
	CR currency.CurrencyID          `json:"currency"`
}

func (it PostItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		BR:         it.broker,
		FO:         it.form,
		CR:         it.cid,
	})
}

type PostItemJSONUnpacker struct {
	BR string          `json:"broker"`
	FO json.RawMessage `json:"form"`
	CR string          `json:"currency"`
}

func (it *PostItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit PostItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.BR, uit.FO, uit.CR)
}
