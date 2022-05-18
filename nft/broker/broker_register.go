package broker

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	BrokerRegisterFactType   = hint.Type("mitum-nft-broker-register-operation-fact")
	BrokerRegisterFactHint   = hint.NewHint(BrokerRegisterFactType, "v0.0.1")
	BrokerRegisterFactHinter = BrokerRegisterFact{BaseHinter: hint.NewBaseHinter(BrokerRegisterFactHint)}
	BrokerRegisterType       = hint.Type("mitum-nft-broker-register-operation")
	BrokerRegisterHint       = hint.NewHint(BrokerRegisterType, "v0.0.1")
	BrokerRegisterHinter     = BrokerRegister{BaseOperation: operationHinter(BrokerRegisterHint)}
)

type BrokerRegisterFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	target base.Address
	policy BrokerPolicy
	cid    currency.CurrencyID
}

func NewBrokerRegisterFact(token []byte, sender base.Address, target base.Address, policy BrokerPolicy, cid currency.CurrencyID) BrokerRegisterFact {
	fact := BrokerRegisterFact{
		BaseHinter: hint.NewBaseHinter(BrokerRegisterFactHint),
		token:      token,
		sender:     sender,
		target:     target,
		policy:     policy,
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
		fact.target.Bytes(),
		fact.policy.Bytes(),
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
		fact.target,
		fact.policy,
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

func (fact BrokerRegisterFact) Target() base.Address {
	return fact.target
}

func (fact BrokerRegisterFact) Policy() BrokerPolicy {
	return fact.policy
}

func (fact BrokerRegisterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.target

	return as, nil
}

func (fact BrokerRegisterFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact BrokerRegisterFact) Rebuild() BrokerRegisterFact {
	policy := fact.policy.Rebuild()
	fact.policy = policy

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
