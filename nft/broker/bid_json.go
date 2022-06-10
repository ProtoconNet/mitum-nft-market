package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type BidFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash  `json:"hash"`
	TK []byte          `json:"token"`
	SD base.Address    `json:"sender"`
	NF nft.NFTID       `json:"nft"`
	AM currency.Amount `json:"amount"`
}

func (fact BidFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BidFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		NF:         fact.nft,
		AM:         fact.amount,
	})
}

type BidFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	NF json.RawMessage     `json:"nft"`
	AM json.RawMessage     `json:"amount"`
}

func (fact *BidFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact BidFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.NF, ufact.AM)
}

func (op *Bid) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
