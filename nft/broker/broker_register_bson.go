package broker

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (form BrokerRegisterForm) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(form.Hint()),
			bson.M{
				"target":    form.target,
				"symbol":    form.symbol,
				"brokerage": form.brokerage,
				"receiver":  form.receiver,
				"royalty":   form.royalty,
				"uri":       form.uri,
			}))
}

type BrokerRegisterFormBSONUnpacker struct {
	TG base.AddressDecoder `bson:"target"`
	SB string              `bson:"symbol"`
	BR uint                `bson:"brokerage"`
	RC base.AddressDecoder `bson:"receiver"`
	RY bool                `bson:"royalty"`
	UR string              `bson:"uri"`
}

func (form *BrokerRegisterForm) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uf BrokerRegisterFormBSONUnpacker
	if err := bson.Unmarshal(b, &uf); err != nil {
		return err
	}

	return form.unpack(enc, uf.TG, uf.SB, uf.BR, uf.RC, uf.RY, uf.UR)
}

func (fact BrokerRegisterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":     fact.h,
				"token":    fact.token,
				"sender":   fact.sender,
				"form":     fact.form,
				"currency": fact.cid,
			}))
}

type BrokerRegisterFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	FO bson.Raw            `bson:"form"`
	CR string              `bson:"currency"`
}

func (fact *BrokerRegisterFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact BrokerRegisterFactBSONUnpacker
	if err := bson.Unmarshal(b, &ufact); err != nil {
		return err
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, ufact.FO, ufact.CR)
}

func (op *BrokerRegister) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
