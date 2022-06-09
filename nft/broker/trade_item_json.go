package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type TradeItemJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID `json:"nft"`
}

func (it TradeItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(TradeItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		NF:         it.n,
	})
}

type TradeItemJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
}

func (it *TradeItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit TradeItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF)
}
