package broker

import (
	"sync"
	"time"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"
	"github.com/relvacode/iso8601"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util/valuehash"
)

var UnpostItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UnpostItemProcessor)
	},
}

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

type UnpostItemProcessor struct {
	cp              *extensioncurrency.CurrencyPool
	h               valuehash.Hash
	sender          base.Address
	item            UnpostItem
	lastConfirmedAt time.Time
	posting         Posting
	pState          state.State
}

func (ipp *UnpostItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := ipp.item.IsValid(nil); err != nil {
		return err
	}

	var n nft.NFT
	nid := ipp.item.NFT()
	if st, err := existsState(collection.StateKeyNFT(nid), "nft", getState); err != nil {
		return err
	} else if nv, err := collection.StateNFTValue(st); err != nil {
		return err
	} else if !nv.Active() {
		return errors.Errorf("burned nft; %q", nid)
	} else if st, err = existsState(collection.StateKeyCollection(nv.ID().Collection()), "design", getState); err != nil {
		return err
	} else if design, err := collection.StateCollectionValue(st); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("deactivated collection; %q", nid.Collection())
	} else {
		n = nv
	}

	if !n.Owner().Equal(ipp.sender) {
		if err := checkExistsState(currency.StateKeyAccount(n.Owner()), getState); err != nil {
			return err
		} else if st, err := existsState(collection.StateKeyAgents(n.Owner(), n.ID().Collection()), "agents", getState); err != nil {
			return errors.Errorf("unathorized sender; %q", ipp.sender)
		} else if box, err := collection.StateAgentsValue(st); err != nil {
			return err
		} else if !box.Exists(ipp.sender) {
			return errors.Errorf("unathorized sender; %q", ipp.sender)
		}
	}

	var posting Posting
	if st, err := existsState(StateKeyPosting(nid), "posting", getState); err != nil {
		return err
	} else if p, err := StatePostingValue(st); err != nil {
		return err
	} else if !p.Active() {
		return errors.Errorf("already unposted nft; %q", nid)
	} else if stt, err := existsState(StateKeyBroker(p.Broker()), "design", getState); err != nil {
		return err
	} else if design, err := StateBrokerValue(stt); err != nil {
		return err
	} else if !design.Active() {
		return errors.Errorf("deactivated broker; %q", p.Broker())
	} else {
		posting = p
		ipp.pState = st
	}

	if posting.Option() == AuctionPostOption {
		closeTime := posting.Details().(AuctionDetails).CloseTime()
		if t, err := iso8601.ParseString(closeTime.String()); err != nil {
			return err
		} else if ipp.lastConfirmedAt.After(t) || ipp.lastConfirmedAt.Equal(t) {
			return errors.Errorf("this auction is already over; %q, closetime: %q, last block: %q", nid, closeTime)
		}
	}

	posting = NewPosting(false, posting.Broker(), posting.Option(), posting.Details())
	if err := posting.IsValid(nil); err != nil {
		return err
	}
	ipp.posting = posting

	return nil
}

func (ipp *UnpostItemProcessor) Process(
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

func (ipp *UnpostItemProcessor) Close() error {
	ipp.cp = nil
	ipp.h = nil
	ipp.sender = nil
	ipp.item = UnpostItem{}
	ipp.posting = Posting{}
	ipp.pState = nil
	ipp.lastConfirmedAt = time.Time{}
	UnpostItemProcessorPool.Put(ipp)

	return nil
}

type UnpostProcessor struct {
	cp *extensioncurrency.CurrencyPool
	Unpost
	ipps         []*UnpostItemProcessor
	amountStates map[currency.CurrencyID]currency.AmountState
	required     map[currency.CurrencyID][2]currency.Big
	mst          storage.Database
}

func NewUnpostProcessor(mst storage.Database, cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(Unpost)
		if !ok {
			return nil, errors.Errorf("not Unpost; %T", op)
		}

		opp := UnpostProcessorPool.Get().(*UnpostProcessor)

		opp.cp = cp
		opp.Unpost = i
		opp.ipps = nil
		opp.amountStates = nil
		opp.required = nil
		opp.mst = mst

		return opp, nil
	}
}

func (opp *UnpostProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(UnpostFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot Unpost nfts; %q", fact.Sender())
	}

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	var lastConfirmedAt time.Time
	switch m, found, err := opp.mst.LastManifest(); {
	case err != nil:
		return nil, err
	case !found:
		return nil, err
	default:
		lastConfirmedAt = m.ConfirmedAt()
	}

	ipps := make([]*UnpostItemProcessor, len(fact.items))
	for i := range fact.items {

		c := UnpostItemProcessorPool.Get().(*UnpostItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.sender = fact.Sender()
		c.item = fact.items[i]
		c.posting = Posting{}
		c.pState = nil
		c.lastConfirmedAt = lastConfirmedAt

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

func (opp *UnpostProcessor) Process(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(UnpostFact)

	var states []state.State

	for i := range opp.ipps {
		if sts, err := opp.ipps[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process Unpost item; %w", err)
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

func (opp *UnpostProcessor) Close() error {
	for i := range opp.ipps {
		_ = opp.ipps[i].Close()
	}

	opp.cp = nil
	opp.Unpost = Unpost{}
	opp.ipps = nil
	opp.amountStates = nil
	opp.required = nil
	opp.mst = nil

	UnpostProcessorPool.Put(opp)

	return nil
}

func (opp *UnpostProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(UnpostFact)

	items := make([]UnpostItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateUnpostItemsFee(opp.cp, items)
}

func CalculateUnpostItemsFee(cp *extensioncurrency.CurrencyPool, items []UnpostItem) (map[currency.CurrencyID][2]currency.Big, error) {
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
