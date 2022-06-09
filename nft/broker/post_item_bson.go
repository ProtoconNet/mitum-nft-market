package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (form PostForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"option":    form.option,
				"nft":       form.n,
				"closetime": form.closeTime,
				"price":     form.price,
			}))
}

type PostFormBSONUnpacker struct {
	OP string   `bson:"option"`
	NF bson.Raw `bson:"nft"`
	CT string   `bson:"closetime"`
	PR bson.Raw `bson:"price"`
}

func (form *PostForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufo PostFormBSONUnpacker
	if err := bson.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.OP, ufo.NF, ufo.CT, ufo.PR)
}

func (it BasePostItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"broker":   it.broker,
				"forms":    it.forms,
				"currency": it.cid,
			}),
	)
}

type BasePostItemBSONUnpacker struct {
	BR string   `bson:"broker"`
	FO bson.Raw `bson:"forms"`
	CR string   `bson:"currency"`
}

func (it *BasePostItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BasePostItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.BR, uit.FO, uit.CR)
}
