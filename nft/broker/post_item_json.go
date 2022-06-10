package broker

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PostFormJSONPacker struct {
	jsonenc.HintedHead
	OP PostOption      `json:"option"`
	NF nft.NFTID       `json:"nft"`
	CT PostCloseTime   `json:"closetime"`
	PR currency.Amount `json:"price"`
}

func (form PostForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		OP:         form.option,
		NF:         form.n,
		CT:         form.closeTime,
		PR:         form.price,
	})
}

type PostFormJSONUnpacker struct {
	OP string          `json:"option"`
	NF json.RawMessage `json:"nft"`
	CT string          `json:"closetime"`
	PR json.RawMessage `json:"price"`
}

func (form *PostForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufo PostFormJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.OP, ufo.NF, ufo.CT, ufo.PR)
}

type PostItemJSONPacker struct {
	jsonenc.HintedHead
	BR extensioncurrency.ContractID `json:"broker"`
	FO []PostForm                   `json:"forms"`
	CR currency.CurrencyID          `json:"currency"`
}

func (it BasePostItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		BR:         it.broker,
		FO:         it.forms,
		CR:         it.cid,
	})
}

type PostItemJSONUnpacker struct {
	BR string          `json:"broker"`
	FO json.RawMessage `json:"forms"`
	CR string          `json:"currency"`
}

func (it *BasePostItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit PostItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.BR, uit.FO, uit.CR)
}
