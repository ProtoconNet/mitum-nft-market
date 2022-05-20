package broker

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var SettleAuctionProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SettleAuctionProcessor)
	},
}

func (SettleAuction) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type SettleAuctionProcessor struct {
	cp *currency.CurrencyPool
	SettleAuction
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewSettleAuctionProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(SettleAuction)
		if !ok {
			return nil, errors.Errorf("not SettleAuction; %T", op)
		}

		opp := SettleAuctionProcessorPool.Get().(*SettleAuctionProcessor)

		opp.cp = cp
		opp.SettleAuction = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *SettleAuctionProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *SettleAuctionProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(SettleAuctionFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *SettleAuctionProcessor) Close() error {
	opp.cp = nil
	opp.SettleAuction = SettleAuction{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	SettleAuctionProcessorPool.Put(opp)

	return nil
}
