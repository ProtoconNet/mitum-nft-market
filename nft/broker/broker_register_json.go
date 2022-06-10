package broker

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type BrokerRegisterFormJSONPacker struct {
	jsonenc.HintedHead
	TG base.Address                 `json:"target"`
	SB extensioncurrency.ContractID `json:"symbol"`
	BR nft.PaymentParameter         `json:"brokerage"`
	RC base.Address                 `json:"receiver"`
	RY bool                         `json:"royalty"`
	UR nft.URI                      `json:"uri"`
}

func (form BrokerRegisterForm) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BrokerRegisterFormJSONPacker{
		HintedHead: jsonenc.NewHintedHead(form.Hint()),
		TG:         form.target,
		SB:         form.symbol,
		BR:         form.brokerage,
		RC:         form.receiver,
		RY:         form.royalty,
		UR:         form.uri,
	})
}

type BrokerRegisterFormJSONUnpacker struct {
	TG base.AddressDecoder `json:"target"`
	SB string              `json:"symbol"`
	BR uint                `json:"brokerage"`
	RY bool                `json:"royalty"`
	RC base.AddressDecoder `json:"receiver"`
	UR string              `json:"uri"`
}

func (form *BrokerRegisterForm) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uf BrokerRegisterFormJSONUnpacker
	if err := enc.Unmarshal(b, &uf); err != nil {
		return err
	}

	return form.unpack(enc, uf.TG, uf.SB, uf.BR, uf.RC, uf.RY, uf.UR)
}

type BrokerRegisterFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	FO BrokerRegisterForm  `json:"form"`
	CR currency.CurrencyID `json:"currency"`
}

func (fact BrokerRegisterFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BrokerRegisterFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		FO:         fact.form,
		CR:         fact.cid,
	})
}

type BrokerRegisterFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	FO json.RawMessage     `json:"form"`
	CR string              `json:"currency"`
}

func (fact *BrokerRegisterFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufact BrokerRegisterFactJSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.FO, ufact.CR)
}

func (op *BrokerRegister) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
