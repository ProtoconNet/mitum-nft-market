package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BrokerPolicyJSONPacker struct {
	jsonenc.HintedHead
	BR nft.PaymentParameter `json:"brokerage"`
	RC base.Address         `json:"receiver"`
	RY bool                 `json:"royalty"`
}

func (bp BrokerPolicy) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BrokerPolicyJSONPacker{
		HintedHead: jsonenc.NewHintedHead(bp.Hint()),
		BR:         bp.brokerage,
		RC:         bp.receiver,
		RY:         bp.royalty,
	})
}

type BrokerPolicyJSONUnpacker struct {
	BR uint                `json:"brokerage"`
	RC base.AddressDecoder `json:"receiver"`
	RY bool                `json:"royalty"`
}

func (cp *BrokerPolicy) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubp BrokerPolicyJSONUnpacker
	if err := enc.Unmarshal(b, &ubp); err != nil {
		return err
	}

	return cp.unpack(enc, ubp.BR, ubp.RC, ubp.RY)
}
