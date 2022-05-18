package broker

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type BrokerRegisterFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	TG base.Address        `json:"target"`
	PL BrokerPolicy        `json:"policy"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact BrokerRegisterFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BrokerRegisterFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		TG:         fact.target,
		PL:         fact.policy,
		CR:         fact.cid,
	})
}

type BrokerRegisterFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	TG base.AddressDecoder `json:"target"`
	PL json.RawMessage     `json:"policy"`
	CR string              `json:"currency"`
}

func (fact *BrokerRegisterFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact BrokerRegisterFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.TG, ufact.PL, ufact.CR)
}

func (op *BrokerRegister) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
