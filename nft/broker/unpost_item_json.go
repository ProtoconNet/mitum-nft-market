package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type UnpostItemJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID           `json:"nft"`
	CR currency.CurrencyID `json:"currency"`
}

func (it UnpostItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UnpostItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		NF:         it.nft,
		CR:         it.cid,
	})
}

type UnpostItemJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	CR string          `json:"currency"`
}

func (it *UnpostItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit UnpostItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.CR)
}
