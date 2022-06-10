package broker

import (
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/operation"
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
	cp *extensioncurrency.CurrencyPool
	BrokerRegister
	DesignState state.State
	design      nft.Design
	amountState currency.AmountState
	fee         currency.Big
}

func NewBrokerRegisterProcessor(cp *extensioncurrency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {

		i, ok := op.(BrokerRegister)
		if !ok {
			return nil, errors.Errorf("not BrokerRegister; %T", op)
		}

		opp := BrokerRegisterProcessorPool.Get().(*BrokerRegisterProcessor)

		opp.cp = cp
		opp.BrokerRegister = i
		opp.DesignState = nil
		opp.design = nft.Design{}
		opp.amountState = currency.AmountState{}
		opp.fee = currency.ZeroBig

		return opp, nil
	}
}

func (opp *BrokerRegisterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(BrokerRegisterFact)

	if err := fact.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot register a broker; %q", fact.Sender())
	}

	if st, err := existsState(extensioncurrency.StateKeyContractAccount(fact.Form().Target()), "contract account", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if ca, err := extensioncurrency.StateContractAccountValue(st); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else if !ca.Owner().Equal(fact.Sender()) {
		return nil, operation.NewBaseReasonError("not owner of contract account; %q", fact.Form().Target())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Form().Receiver()), getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	if err := checkNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Form().Receiver()), getState); err != nil {
		return nil, operation.NewBaseReasonError("contract account cannot receive brokerage fee; %q", fact.Form().Receiver())
	}

	if st, err := notExistsState(StateKeyBroker(fact.Form().Symbol()), "design", getState); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	} else {
		opp.DesignState = st
	}

	policy := NewBrokerPolicy(fact.Form().Brokerage(), fact.Form().Receiver(), fact.Form().Royalty())
	if err := policy.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}

	design := nft.NewDesign(fact.Form().Target(), fact.Sender(), fact.Form().Symbol(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, operation.NewBaseReasonError(err.Error())
	}
	opp.design = design

	if err := checkFactSignsByState(fact.Sender(), opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing; %w", err)
	}

	if st, err := existsState(
		currency.StateKeyBalance(fact.Sender(), fact.Currency()), "balance of sender", getState); err != nil {
		return nil, err
	} else {
		opp.amountState = currency.NewAmountState(st, fact.Currency())
	}

	feeer, found := opp.cp.Feeer(fact.Currency())
	if !found {
		return nil, operation.NewBaseReasonError("currency not found; %q", fact.Currency())
	}

	fee, err := feeer.Fee(currency.ZeroBig)
	if err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	}
	switch b, err := currency.StateBalanceValue(opp.amountState); {
	case err != nil:
		return nil, operation.NewBaseReasonErrorFromError(err)
	case b.Big().Compare(fee) < 0:
		return nil, operation.NewBaseReasonError("insufficient balance with fee")
	default:
		opp.fee = fee
	}

	return opp, nil
}

func (opp *BrokerRegisterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(BrokerRegisterFact)

	var states []state.State

	if st, err := SetStateBrokerValue(opp.DesignState, opp.design); err != nil {
		return operation.NewBaseReasonError(err.Error())
	} else {
		states = append(states, st)
	}

	opp.amountState = opp.amountState.Sub(opp.fee).AddFee(opp.fee)
	states = append(states, opp.amountState)

	return setState(fact.Hash(), states...)
}

func (opp *BrokerRegisterProcessor) Close() error {
	opp.cp = nil
	opp.BrokerRegister = BrokerRegister{}
	opp.DesignState = nil
	opp.design = nft.Design{}
	opp.amountState = currency.AmountState{}
	opp.fee = currency.ZeroBig

	BrokerRegisterProcessorPool.Put(opp)

	return nil
}
