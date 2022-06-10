package broker

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type SettleAuctionFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	NF nft.NFTID           `json:"nft"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact SettleAuctionFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SettleAuctionFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		NF:         fact.nft,
		CR:         fact.cid,
	})
}

type SettleAuctionFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	NF json.RawMessage     `json:"nft"`
	CR string              `json:"currency"`
}

func (fact *SettleAuctionFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact SettleAuctionFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.NF, ufact.CR)
}

func (op *SettleAuction) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
