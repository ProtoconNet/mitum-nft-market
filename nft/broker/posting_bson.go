package broker

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (posting Posting) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(posting.Hint()),
		bson.M{
			"active":    posting.active,
			"broker":    posting.broker,
			"option":    posting.option,
			"nft":       posting.n,
			"closetime": posting.closeTime,
			"price":     posting.price,
		}),
	)
}

type PostingBSONUnpacker struct {
	AC bool     `bson:"active"`
	BR string   `bson:"broker"`
	OP string   `bson:"option"`
	NF bson.Raw `bson:"nft"`
	CT string   `bson:"closetime"`
	PR bson.Raw `bson:"price"`
}

func (posting *Posting) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var upt PostingBSONUnpacker
	if err := enc.Unmarshal(b, &upt); err != nil {
		return err
	}

	return posting.unpack(enc, upt.AC, upt.BR, upt.OP, upt.NF, upt.CT, upt.PR)
}

func (bid Bidding) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(bid.Hint()),
		bson.M{
			"bidder": bid.bidder,
			"amount": bid.amount,
		}),
	)
}

type BiddingBSONUnpacker struct {
	BD base.AddressDecoder `bson:"bidder"`
	AM bson.Raw            `bson:"amount"`
}

func (bid *Bidding) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubid BiddingBSONUnpacker
	if err := enc.Unmarshal(b, &ubid); err != nil {
		return err
	}

	return bid.unpack(enc, ubid.BD, ubid.AM)
}
