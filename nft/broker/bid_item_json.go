package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BidItemJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID       `json:"nft"`
	AM currency.Amount `json:"amount"`
}

func (it BidItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BidItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		NF:         it.n,
		AM:         it.amount,
	})
}

type BidItemJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	AM json.RawMessage `json:"amount"`
}

func (it *BidItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit BidItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.AM)
}
