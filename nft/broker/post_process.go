package broker

import (
	"sync"
	"time"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"
	"github.com/relvacode/iso8601"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var PostItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PostItemProcessor)
	},
}

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

type PostItemProcessor struct {
	cp              *extensioncurrency.CurrencyPool
	h               valuehash.Hash
	sender          base.Address
	item            PostItem
	lastConfirmedAt time.Time
	posting         Posting
	pState          state.State
}

func (ipp *PostItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	if st, err := existsState(StateKeyBroker(ipp.item.Broker()), "design", getState); err != nil {
		return err
	} else if design, err := StateBrokerValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("dead broker; %q", ipp.item.Broker())
	}

	form := ipp.item.Form()
	nid := form.Details().NFT()
	if st, err := existsState(collection.StateKeyNFT(nid), "nft", getState); err != nil {
		return err
	} else if n, err := collection.StateNFTValue(st); err != nil {
		return err
	} else if stt, err := existsState(collection.StateKeyCollection(n.ID().Collection()), "design", getState); err != nil {
		return err
	} else if design, err := collection.StateCollectionValue(stt); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("dead collection; %q", nid.Collection())
	}

	if form.Option() == AuctionPostOption {
		closeTime := form.Details().(AuctionDetails).CloseTime()
		if t, err := iso8601.ParseString(closeTime.String()); err != nil {
			return err
		} else if ipp.lastConfirmedAt.After(t) || ipp.lastConfirmedAt.Equal(t) {
			return errors.Errorf("closetime is faster than last confirmed_at; %q -> %q", closeTime, ipp.lastConfirmedAt.String())
		}
	}

	posting := NewPosting(true, ipp.item.Broker(), form.Option(), form.Details())
	if err := posting.IsValid(nil); err != nil {
		return err
	}

	switch st, found, err := getState(StateKeyPosting(nid)); {
	case err != nil:
		return err
	case found:
		if pt, err := StatePostingValue(st); err != nil {
			return err
		} else if pt.Active() {
			return errors.Errorf("posting already exists; %q, %q", ipp.item.Broker(), form.Details().NFT())
		} else {
			ipp.posting = posting
			ipp.pState = st
		}
	default:
		ipp.posting = posting
		ipp.pState = st
	}

	return nil
}

func (ipp *PostItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	var states []state.State

	if st, err := SetStatePostingValue(ipp.pState, ipp.posting); err != nil {
		return nil, err
	} else {
		states = append(states, st)
	}

	return states, nil
}

func (ipp *PostItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.sender = nil
	ipp.item = PostItem{}
	ipp.posting = Posting{}
	ipp.pState = nil
	ipp.lastConfirmedAt = time.Time{}
	PostItemProcessorPool.Put(ipp)

	return nil
}

type PostProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Post
	ipps            []*PostItemProcessor
	amountStates    map[currency.CurrencyID]currency.AmountState
	required        map[currency.CurrencyID][2]currency.Big
	lastConfirmedAt time.Time
}

func NewPostProcessor(lastConfiremdAt time.Time, cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Post)
		if !ok {
			return nil, errors.Errorf("not Post; %T", op)
		}

		opp := PostProcessorPool.Get().(*PostProcessor)

		opp.cp = cp
		opp.Post = i
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil
		opp.lastConfirmedAt = lastConfiremdAt

		return opp, nil
	}
}

func (opp *PostProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(PostFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot Post nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	ipps := make([]*PostItemProcessor, len(fact.items))
	for i := range fact.items {

		c := PostItemProcessorPool.Get().(*PostItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.sender = fact.Sender()
		c.item = fact.items[i]
		c.posting = Posting{}
		c.pState = nil
		c.lastConfirmedAt = opp.lastConfirmedAt

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonError(err.Error())
		}

		ipps[i] = c
	}

	opp.ipps = ipps

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else if sts, err := CheckSenderEnoughBalance(fact.Sender(), required, getState); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee; %w", err)
	} else {
		opp.required = required
		opp.amountStates = sts
	}

	return opp, nil
}

func (opp *PostProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(PostFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process Post item; %w", err)
		} else {
			states = append(states, sts...)
		}
	}

	for k := range opp.required {
		rq := opp.required[k]
		states = append(states, opp.amountStates[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), states...)
}

func (opp *PostProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Post = Post{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil
	opp.lastConfirmedAt = time.Time{}

	PostProcessorPool.Put(opp)

	return nil
}

func (opp *PostProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(PostFact)

	items := make([]PostItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculatePostItemsFee(opp.cp, items)
}

func CalculatePostItemsFee(cp *extensioncurrency.CurrencyPool, items []PostItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for i := range items {
		it := items[i]

		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}

		if cp == nil {
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}
			continue
		}

		feeer, found := cp.Feeer(it.Currency())
		if !found {
			return nil, errors.Errorf("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[it.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}

func CheckSenderEnoughBalance(
	holder base.Address,
	required map[currency.CurrencyID][2]currency.Big,
	getState func(key string) (state.State, bool, error),
) (map[currency.CurrencyID]currency.AmountState, error) {
	sb := map[currency.CurrencyID]currency.AmountState{}

	for cid := range required {
		rq := required[cid]

		st, err := existsState(currency.StateKeyBalance(holder, cid), "currency of holder", getState)
		if err != nil {
			return nil, err
		}

		am, err := currency.StateBalanceValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError("insufficient balance of sender: %w", err)
		}

		if am.Big().Compare(rq[0]) < 0 {
			return nil, operation.NewBaseReasonError(
				"insufficient balance of sender, %s; %d !> %d", holder.String(), am.Big(), rq[0])
		} else {
			sb[cid] = currency.NewAmountState(st, cid)
		}
	}

	return sb, nil
}
