package broker

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var UnpostProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UnpostProcessor)
	},
}

func (Unpost) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type UnpostProcessor struct {
	cp *currency.CurrencyPool
	Unpost
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewUnpostProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Unpost)
		if !ok {
			return nil, errors.Errorf("not Unpost; %T", op)
		}

		opp := UnpostProcessorPool.Get().(*UnpostProcessor)

		opp.cp = cp
		opp.Unpost = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *UnpostProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *UnpostProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(UnpostFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *UnpostProcessor) Close() error {
	opp.cp = nil
	opp.Unpost = Unpost{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	UnpostProcessorPool.Put(opp)

	return nil
}
