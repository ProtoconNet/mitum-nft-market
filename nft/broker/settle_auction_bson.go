package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact SettleAuctionFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":     fact.h,
				"token":    fact.token,
				"sender":   fact.sender,
				"nft":      fact.nft,
				"currency": fact.cid,
			}))
}

type SettleAuctionFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	NF bson.Raw            `bson:"nft"`
	CR string              `bson:"currency"`
}

func (fact *SettleAuctionFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact SettleAuctionFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.NF, ufact.CR)
}

func (op *SettleAuction) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
