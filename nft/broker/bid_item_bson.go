package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BidItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"nft":    it.n,
				"amount": it.amount,
			}),
	)
}

type BidItemBSONUnpacker struct {
	NF bson.Raw `bson:"nft"`
	AM bson.Raw `bson:"amount"`
}

func (it *BidItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BidItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.AM)
}
