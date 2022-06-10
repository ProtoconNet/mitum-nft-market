package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	BrokerPolicyType   = hint.Type("mitum-nft-market-broker-policy")
	BrokerPolicyHint   = hint.NewHint(BrokerPolicyType, "v0.0.1")
	BrokerPolicyHinter = BrokerPolicy{BaseHinter: hint.NewBaseHinter(BrokerPolicyHint)}
)

type BrokerPolicy struct {
	hint.BaseHinter
	brokerage nft.PaymentParameter
	receiver  base.Address
	royalty   bool
}

func NewBrokerPolicy(brokerage nft.PaymentParameter, receiver base.Address, royalty bool) BrokerPolicy {
	return BrokerPolicy{
		BaseHinter: hint.NewBaseHinter(BrokerPolicyHint),
		brokerage:  brokerage,
		receiver:   receiver,
		royalty:    royalty,
	}
}

func MustNewBrokerPolicy(brokerage nft.PaymentParameter, receiver base.Address, royalty bool) BrokerPolicy {
	policy := NewBrokerPolicy(brokerage, receiver, royalty)

	if err := policy.IsValid(nil); err != nil {
		panic(err)
	}

	return policy
}

func (policy BrokerPolicy) Bytes() []byte {
	br := make([]byte, 1)
	if policy.royalty {
		br[0] = 1
	} else {
		br[0] = 0
	}

	return util.ConcatBytesSlice(
		policy.brokerage.Bytes(),
		policy.receiver.Bytes(),
		br,
	)
}

func (policy BrokerPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		policy.BaseHinter,
		policy.brokerage,
		policy.receiver); err != nil {
		return err
	}

	return nil
}

func (policy BrokerPolicy) Brokerage() nft.PaymentParameter {
	return policy.brokerage
}

func (policy BrokerPolicy) Receiver() base.Address {
	return policy.receiver
}

func (policy BrokerPolicy) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = policy.receiver
	return as, nil
}

func (policy BrokerPolicy) Royalty() bool {
	return policy.royalty
}

func (policy BrokerPolicy) Rebuild() BrokerPolicy {
	return policy
}
