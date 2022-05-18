package broker

import (
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it PostItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"posting":  it.posting,
				"currency": it.cid,
			}),
	)
}

type PostItemBSONUnpacker struct {
	PO bson.Raw `bson:"posting"`
	CR string   `bson:"currency"`
}

func (it *PostItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uca PostItemBSONUnpacker
	if err := bson.Unmarshal(b, &uca); err != nil {
		return err
	}

	return it.unpack(enc, uca.PO, uca.CR)
}
