package broker

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var TradeProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TradeProcessor)
	},
}

func (Trade) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type TradeProcessor struct {
	cp *currency.CurrencyPool
	Trade
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewTradeProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Trade)
		if !ok {
			return nil, errors.Errorf("not Trade; %T", op)
		}

		opp := TradeProcessorPool.Get().(*TradeProcessor)

		opp.cp = cp
		opp.Trade = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *TradeProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *TradeProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(TradeFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *TradeProcessor) Close() error {
	opp.cp = nil
	opp.Trade = Trade{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	TradeProcessorPool.Put(opp)

	return nil
}
