package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BaseUnpostItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"nfts":     it.nfts,
				"currency": it.cid,
			}),
	)
}

type BaseUnpostItemBSONUnpacker struct {
	NS bson.Raw `bson:"nfts"`
	CR string   `bson:"currency"`
}

func (it *BaseUnpostItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BaseUnpostItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NS, uit.CR)
}
