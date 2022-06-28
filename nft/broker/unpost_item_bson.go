package broker

import (
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it UnpostItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"nft":      it.nft,
				"currency": it.cid,
			}),
	)
}

type UnpostItemBSONUnpacker struct {
	NF bson.Raw `bson:"nft"`
	CR string   `bson:"currency"`
}

func (it *UnpostItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit UnpostItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF, uit.CR)
}
