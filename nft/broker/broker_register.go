package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	BrokerRegisterFormType   = hint.Type("mitum-nft-market-broker-register-form")
	BrokerRegisterFormHint   = hint.NewHint(BrokerRegisterFormType, "v0.0.1")
	BrokerRegisterFormHinter = BrokerRegisterForm{BaseHinter: hint.NewBaseHinter(BrokerRegisterFormHint)}
)

type BrokerRegisterForm struct {
	hint.BaseHinter
	target    base.Address
	symbol    extensioncurrency.ContractID
	brokerage nft.PaymentParameter
	receiver  base.Address
	royalty   bool
	uri       nft.URI
}

func NewBrokerRegisterForm(target base.Address, symbol extensioncurrency.ContractID, brokerage nft.PaymentParameter, receiver base.Address, royalty bool, uri nft.URI) BrokerRegisterForm {
	return BrokerRegisterForm{
		BaseHinter: hint.NewBaseHinter(BrokerRegisterFormHint),
		target:     target,
		symbol:     symbol,
		brokerage:  brokerage,
		royalty:    royalty,
		receiver:   receiver,
		uri:        uri,
	}
}

func MustNewBrokerRegisterForm(target base.Address, symbol extensioncurrency.ContractID, brokerage nft.PaymentParameter, receiver base.Address, royalty bool, uri nft.URI) BrokerRegisterForm {
	form := NewBrokerRegisterForm(target, symbol, brokerage, receiver, royalty, uri)
	if err := form.IsValid(nil); err != nil {
		panic(err)
	}
	return form
}

func (form BrokerRegisterForm) Bytes() []byte {
	br := make([]byte, 1)
	if form.royalty {
		br[0] = 1
	} else {
		br[0] = 0
	}

	return util.ConcatBytesSlice(
		form.target.Bytes(),
		form.symbol.Bytes(),
		form.brokerage.Bytes(),
		form.receiver.Bytes(),
		br,
		form.uri.Bytes(),
	)
}

func (form BrokerRegisterForm) Target() base.Address {
	return form.target
}

func (form BrokerRegisterForm) Symbol() extensioncurrency.ContractID {
	return form.symbol
}

func (form BrokerRegisterForm) Brokerage() nft.PaymentParameter {
	return form.brokerage
}

func (form BrokerRegisterForm) Royalty() bool {
	return form.royalty
}

func (form BrokerRegisterForm) Receiver() base.Address {
	return form.receiver
}

func (form BrokerRegisterForm) Uri() nft.URI {
	return form.uri
}

func (form BrokerRegisterForm) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)
	as[0] = form.target
	as[1] = form.receiver
	return as, nil
}

func (form BrokerRegisterForm) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		form.BaseHinter,
		form.target,
		form.symbol,
		form.brokerage,
		form.receiver,
		form.uri,
	); err != nil {
		return err
	}

	return nil
}

func (form BrokerRegisterForm) Rebuild() BrokerRegisterForm {
	return form
}

var (
	BrokerRegisterFactType   = hint.Type("mitum-nft-market-broker-register-operation-fact")
	BrokerRegisterFactHint   = hint.NewHint(BrokerRegisterFactType, "v0.0.1")
	BrokerRegisterFactHinter = BrokerRegisterFact{BaseHinter: hint.NewBaseHinter(BrokerRegisterFactHint)}
	BrokerRegisterType       = hint.Type("mitum-nft-market-broker-register-operation")
	BrokerRegisterHint       = hint.NewHint(BrokerRegisterType, "v0.0.1")
	BrokerRegisterHinter     = BrokerRegister{BaseOperation: operationHinter(BrokerRegisterHint)}
)

type BrokerRegisterFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	form   BrokerRegisterForm
	cid    currency.CurrencyID
}

func NewBrokerRegisterFact(token []byte, sender base.Address, form BrokerRegisterForm, cid currency.CurrencyID) BrokerRegisterFact {
	fact := BrokerRegisterFact{
		BaseHinter: hint.NewBaseHinter(BrokerRegisterFactHint),
		token:      token,
		sender:     sender,
		form:       form,
		cid:        cid,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact BrokerRegisterFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact BrokerRegisterFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact BrokerRegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.form.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact BrokerRegisterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for BrokerRegisterFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.form,
		fact.cid); err != nil {
		return err
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact BrokerRegisterFact) Token() []byte {
	return fact.token
}

func (fact BrokerRegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact BrokerRegisterFact) Form() BrokerRegisterForm {
	return fact.form
}

func (fact BrokerRegisterFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	as = append(as, fact.sender)

	if fas, err := fact.form.Addresses(); err != nil {
		return nil, err
	} else {
		as = append(as, fas...)
	}

	return as, nil
}

func (fact BrokerRegisterFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact BrokerRegisterFact) Rebuild() BrokerRegisterFact {
	fact.form = fact.form.Rebuild()
	fact.h = fact.GenerateHash()

	return fact
}

type BrokerRegister struct {
	currency.BaseOperation
}

func NewBrokerRegister(fact BrokerRegisterFact, fs []base.FactSign, memo string) (BrokerRegister, error) {
	bo, err := currency.NewBaseOperationFromFact(BrokerRegisterHint, fact, fs, memo)
	if err != nil {
		return BrokerRegister{}, err
	}
	return BrokerRegister{BaseOperation: bo}, nil
}
