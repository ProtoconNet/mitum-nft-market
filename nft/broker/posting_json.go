package broker

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type PostingJSONPacker struct {
	jsonenc.HintedHead
	BR extensioncurrency.ContractID `json:"broker"`
	OP PostOption                   `json:"option"`
	NF nft.NFTID                    `json:"nft"`
	CT PostCloseTime                `json:"closetime"`
	PR currency.Amount              `json:"price"`
}

func (posting Posting) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostingJSONPacker{
		HintedHead: jsonenc.NewHintedHead(posting.Hint()),
		BR:         posting.broker,
		OP:         posting.option,
		NF:         posting.nft,
		CT:         posting.closeTime,
		PR:         posting.price,
	})
}

type PostingJSONUnpacker struct {
	BR string          `json:"broker"`
	OP string          `json:"option"`
	NF json.RawMessage `json:"nft"`
	CT string          `json:"closetime"`
	PR json.RawMessage `json:"price"`
}

func (cp *Posting) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var upt PostingJSONUnpacker
	if err := enc.Unmarshal(b, &upt); err != nil {
		return err
	}

	return cp.unpack(enc, upt.BR, upt.OP, upt.NF, upt.CT, upt.PR)
}

type BiddingJSONPacker struct {
	jsonenc.HintedHead
	BD base.Address    `json:"bidder"`
	AM currency.Amount `json:"amount"`
}

func (bidding Bidding) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BiddingJSONPacker{
		HintedHead: jsonenc.NewHintedHead(bidding.Hint()),
		BD:         bidding.bidder,
		AM:         bidding.amount,
	})
}

type BiddingJSONUnpacker struct {
	BD base.AddressDecoder `json:"bidder"`
	AM json.RawMessage     `json:"amount"`
}

func (bid *Bidding) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubd BiddingJSONUnpacker
	if err := enc.Unmarshal(b, &ubd); err != nil {
		return err
	}

	return bid.unpack(enc, ubd.BD, ubd.AM)
}
