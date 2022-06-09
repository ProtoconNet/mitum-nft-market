package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type UnpostItemJSONPacker struct {
	jsonenc.HintedHead
	NS []nft.NFTID         `json:"nfts"`
	CR currency.CurrencyID `json:"currency"`
}

func (it BaseUnpostItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UnpostItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		NS:         it.nfts,
		CR:         it.cid,
	})
}

type UnpostItemJSONUnpacker struct {
	NS json.RawMessage `json:"nfts"`
	CR string          `json:"currency"`
}

func (it *BaseUnpostItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uit UnpostItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NS, uit.CR)
}
