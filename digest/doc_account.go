package digest

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

type AccountDoc struct {
	mongodbstorage.BaseDoc
	address string
	height  base.Height
	pubs    []string
}

func NewAccountDoc(rs AccountValue, enc encoder.Encoder) (AccountDoc, error) {
	b, err := mongodbstorage.NewBaseDoc(nil, rs, enc)
	if err != nil {
		return AccountDoc{}, err
	}

	var pubs []string
	if keys := rs.Account().Keys(); keys != nil {
		ks := keys.Keys()
		pubs = make([]string, len(ks))
		for i := range ks {
			k := ks[i].Key()
			pubs[i] = k.String()
		}
	}

	address := rs.ac.Address()
	return AccountDoc{
		BaseDoc: b,
		address: address.String(),
		height:  rs.height,
		pubs:    pubs,
	}, nil
}

func (doc AccountDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["address"] = doc.address
	m["height"] = doc.height
	m["pubs"] = doc.pubs

	return bsonenc.Marshal(m)
}

type BalanceDoc struct {
	mongodbstorage.BaseDoc
	st state.State
	am currency.Amount
}

// NewBalanceDoc gets the State of Amount
func NewBalanceDoc(st state.State, enc encoder.Encoder) (BalanceDoc, error) {
	am, err := currency.StateBalanceValue(st)
	if err != nil {
		return BalanceDoc{}, errors.Wrap(err, "balanceDoc needs Amount state")
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return BalanceDoc{}, err
	}

	return BalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc BalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	address := doc.st.Key()[:len(doc.st.Key())-len(currency.StateKeyBalanceSuffix)-len(doc.am.Currency())-1]
	m["address"] = address
	m["currency"] = doc.am.Currency().String()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type ContractAccountStatusDoc struct {
	mongodbstorage.BaseDoc
	st state.State
}

// NewContractAccountStatusDoc gets the State of contract account status
func NewContractAccountStatusDoc(st state.State, enc encoder.Encoder) (ContractAccountStatusDoc, error) {
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return ContractAccountStatusDoc{}, err
	}
	return ContractAccountStatusDoc{
		BaseDoc: b,
		st:      st,
	}, nil
}

func (doc ContractAccountStatusDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	address := doc.st.Key()[:len(doc.st.Key())-len(extensioncurrency.StateKeyContractAccountSuffix)]
	m["address"] = address
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
