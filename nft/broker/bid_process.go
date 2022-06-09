package broker

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var BidProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(BidProcessor)
	},
}

func (Bid) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type BidProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Bid
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewBidProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Bid)
		if !ok {
			return nil, errors.Errorf("not Bid; %T", op)
		}

		opp := BidProcessorPool.Get().(*BidProcessor)

		opp.cp = cp
		opp.Bid = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *BidProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *BidProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(BidFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *BidProcessor) Close() error {
	opp.cp = nil
	opp.Bid = Bid{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	BidProcessorPool.Put(opp)

	return nil
}
