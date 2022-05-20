package broker

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var BrokerRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(BrokerRegisterProcessor)
	},
}

func (BrokerRegister) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type BrokerRegisterProcessor struct {
	cp *currency.CurrencyPool
	BrokerRegister
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewBrokerRegisterProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(BrokerRegister)
		if !ok {
			return nil, errors.Errorf("not BrokerRegister; %T", op)
		}

		opp := BrokerRegisterProcessorPool.Get().(*BrokerRegisterProcessor)

		opp.cp = cp
		opp.BrokerRegister = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *BrokerRegisterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *BrokerRegisterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(BrokerRegisterFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *BrokerRegisterProcessor) Close() error {
	opp.cp = nil
	opp.BrokerRegister = BrokerRegister{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	BrokerRegisterProcessorPool.Put(opp)

	return nil
}
