package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (form PostForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"option":  form.option,
				"details": form.details,
			}))
}

type PostFormBSONUnpacker struct {
	OP string   `bson:"option"`
	DE bson.Raw `bson:"details"`
}

func (form *PostForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufo PostFormBSONUnpacker
	if err := bson.Unmarshal(b, &ufo); err != nil {
		return err
	}

	return form.unpack(enc, ufo.OP, ufo.DE)
}

func (it PostItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"broker":   it.broker,
				"form":     it.form,
				"currency": it.cid,
			}),
	)
}

type PostItemBSONUnpacker struct {
	BR string   `bson:"broker"`
	FO bson.Raw `bson:"form"`
	CR string   `bson:"currency"`
}

func (it *PostItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit PostItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.BR, uit.FO, uit.CR)
}
