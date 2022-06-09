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
	broker := NewBrokerPolicy(brokerage, receiver, royalty)

	if err := broker.IsValid(nil); err != nil {
		panic(err)
	}

	return broker
}

func (broker BrokerPolicy) Bytes() []byte {
	if broker.royalty {
		return util.ConcatBytesSlice(
			broker.brokerage.Bytes(),
			broker.receiver.Bytes(),
			[]byte{1},
		)
	}

	return util.ConcatBytesSlice(
		broker.brokerage.Bytes(),
		broker.receiver.Bytes(),
		[]byte{0},
	)
}

func (broker BrokerPolicy) IsValid([]byte) error {

	if err := isvalid.Check(nil, false,
		broker.BaseHinter,
		broker.brokerage,
		broker.receiver); err != nil {
		return err
	}

	return nil
}

func (broker BrokerPolicy) Brokerage() nft.PaymentParameter {
	return broker.brokerage
}

func (broker BrokerPolicy) Receiver() base.Address {
	return broker.receiver
}

func (broker BrokerPolicy) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = broker.receiver
	return as, nil
}

func (broker BrokerPolicy) Royalty() bool {
	return broker.royalty
}

func (broker BrokerPolicy) Rebuild() BrokerPolicy {
	return broker
}
