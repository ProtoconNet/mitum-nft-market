package broker

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var PostProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PostProcessor)
	},
}

func (Post) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type PostProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Post
	sa  state.State
	sb  currency.AmountState
	fee currency.Big
}

func NewPostProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Post)
		if !ok {
			return nil, errors.Errorf("not Post; %T", op)
		}

		opp := PostProcessorPool.Get().(*PostProcessor)

		opp.cp = cp
		opp.Post = i
		opp.sa = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *PostProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {

	return opp, nil
}

func (opp *PostProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(PostFact)

	var sts []state.State

	return setState(fact.Hash(), sts...)
}

func (opp *PostProcessor) Close() error {
	opp.cp = nil
	opp.Post = Post{}
	opp.sa = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	PostProcessorPool.Put(opp)

	return nil
}
