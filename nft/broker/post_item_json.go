package broker

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PostItemJSONPacker struct {
	jsonenc.HintedHead
	PO Posting             `json:"posting"`
	CR currency.CurrencyID `json:"currency"`
}

func (it PostItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		PO:         it.posting,
		CR:         it.cid,
	})
}

type PostItemJSONUnpacker struct {
	PO json.RawMessage `json:"posting"`
	CR string          `json:"currency"`
}

func (it *PostItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var upn PostItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &upn); err != nil {
		return err
	}

	return it.unpack(enc, upn.PO, upn.CR)
}
