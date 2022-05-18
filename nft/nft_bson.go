package nft

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (nid NFTID) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(nid.Hint()),
		bson.M{
			"collection": nid.collection,
			"idx":        nid.idx,
		}),
	)
}

type NFTIDBSONUnpacker struct {
	CL string       `bson:"collection"`
	IX currency.Big `bson:"id"`
}

func (nid *NFTID) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var unid NFTIDBSONUnpacker
	if err := enc.Unmarshal(b, &unid); err != nil {
		return err
	}

	return nid.unpack(enc, unid.CL, unid.IX)
}

func (cr Copyrighter) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(cr.Hint()),
		bson.M{
			"set":     cr.set,
			"address": cr.address,
		}),
	)
}

type CopyrighterBSONUnpacker struct {
	ST bool                `bson:"set"`
	AD base.AddressDecoder `bson:"address"`
}

func (cr *Copyrighter) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ucr CopyrighterBSONUnpacker
	if err := enc.Unmarshal(b, &ucr); err != nil {
		return err
	}

	return cr.unpack(enc, ucr.ST, ucr.AD)
}

func (nft NFT) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(nft.Hint()),
		bson.M{
			"id":          nft.id,
			"owner":       nft.owner,
			"hash":        nft.hash,
			"uri":         nft.uri,
			"approved":    nft.approved,
			"copyrighter": nft.copyrighter,
		}),
	)
}

type NFTBSONUnpacker struct {
	ID bson.Raw            `bson:"id"`
	ON base.AddressDecoder `bson:"owner"`
	HS string              `bson:"hash"`
	UR string              `bson:"uri"`
	AP base.AddressDecoder `bson:"approved"`
	CP bson.Raw            `bson:"copyrighter"`
}

func (nft *NFT) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var unft NFTBSONUnpacker
	if err := enc.Unmarshal(b, &unft); err != nil {
		return err
	}

	return nft.unpack(enc, unft.ID, unft.ON, unft.HS, unft.UR, unft.AP, unft.CP)
}
