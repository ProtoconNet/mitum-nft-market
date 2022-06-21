package broker

import (
	"encoding/json"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type SellDetailsJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID       `json:"nft"`
	PR currency.Amount `json:"price"`
}

func (details SellDetails) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SellDetailsJSONPacker{
		HintedHead: jsonenc.NewHintedHead(details.Hint()),
		NF:         details.nft,
		PR:         details.price,
	})
}

type SellDetailsJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	PR json.RawMessage `json:"price"`
}

func (details *SellDetails) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ude SellDetailsJSONUnpacker
	if err := enc.Unmarshal(b, &ude); err != nil {
		return err
	}

	return details.unpack(enc, ude.NF, ude.PR)
}

type AuctionDetailsJSONPacker struct {
	jsonenc.HintedHead
	NF nft.NFTID       `json:"nft"`
	CT PostCloseTime   `json:"closetime"`
	PR currency.Amount `json:"price"`
}

func (details AuctionDetails) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(AuctionDetailsJSONPacker{
		HintedHead: jsonenc.NewHintedHead(details.Hint()),
		NF:         details.nft,
		CT:         details.closeTime,
		PR:         details.price,
	})
}

type AuctionDetailsJSONUnpacker struct {
	NF json.RawMessage `json:"nft"`
	CT string          `json:"closetime"`
	PR json.RawMessage `json:"price"`
}

func (details *AuctionDetails) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ude AuctionDetailsJSONUnpacker
	if err := enc.Unmarshal(b, &ude); err != nil {
		return err
	}

	return details.unpack(enc, ude.NF, ude.CT, ude.PR)
}

type PostingJSONPacker struct {
	jsonenc.HintedHead
	AC bool                         `json:"active"`
	BR extensioncurrency.ContractID `json:"broker"`
	OP PostOption                   `json:"option"`
	DE PostDetails                  `json:"details"`
}

func (posting Posting) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(PostingJSONPacker{
		HintedHead: jsonenc.NewHintedHead(posting.Hint()),
		AC:         posting.active,
		BR:         posting.broker,
		OP:         posting.option,
		DE:         posting.details,
	})
}

type PostingJSONUnpacker struct {
	AC bool            `json:"active"`
	BR string          `json:"broker"`
	OP string          `json:"option"`
	DE json.RawMessage `json:"details"`
}

func (cp *Posting) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var upt PostingJSONUnpacker
	if err := enc.Unmarshal(b, &upt); err != nil {
		return err
	}

	return cp.unpack(enc, upt.AC, upt.BR, upt.OP, upt.DE)
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
