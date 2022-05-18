package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (bp BrokerPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(bp.Hint()),
		bson.M{
			"symbol":    bp.symbol,
			"brokerage": bp.brokerage,
			"receiver":  bp.receiver,
			"royalty":   bp.royalty,
		},
	))
}

type BrokerPolicyBSONUnpacker struct {
	SB string              `bson:"symbol"`
	BR uint                `bson:"brokerage"`
	RC base.AddressDecoder `bson:"receiver"`
	RY bool                `bson:"royalty"`
}

func (bp BrokerPolicy) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ucp BrokerPolicyBSONUnpacker
	if err := enc.Unmarshal(b, &ucp); err != nil {
		return err
	}

	return bp.unpack(enc, ucp.SB, ucp.BR, ucp.RC, ucp.RY)
}
