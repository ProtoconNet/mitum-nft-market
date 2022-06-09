package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it TradeItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"nft": it.n,
			}),
	)
}

type BaseTradeItemBSONUnpacker struct {
	NF bson.Raw `bson:"nft"`
}

func (it *TradeItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit BaseTradeItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, uit.NF)
}
