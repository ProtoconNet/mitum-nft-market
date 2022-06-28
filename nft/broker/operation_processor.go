package broker // nolint: dupl, revive

import (
	"fmt"
	"io"
	"sync"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/logging"
	"github.com/spikeekips/mitum/util/valuehash"
)

var operationProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(OperationProcessor)
	},
}

type DuplicationType string

const (
	DuplicationTypeSender   DuplicationType = "sender"
	DuplicationTypeCurrency DuplicationType = "currency"
)

type OperationProcessor struct {
	id string
	sync.RWMutex
	*logging.Logging
	processorHintSet     *hint.Hintmap
	cp                   *extensioncurrency.CurrencyPool
	pool                 *storage.Statepool
	fee                  map[currency.CurrencyID]currency.Big
	amountPool           map[string]currency.AmountState
	duplicated           map[string]DuplicationType
	duplicatedNewAddress map[string]struct{}
	processorClosers     *sync.Map
}

func NewOperationProcessor(cp *extensioncurrency.CurrencyPool) *OperationProcessor {
	return &OperationProcessor{
		Logging: logging.NewLogging(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "mitum-currency-operations-processor")
		}),
		processorHintSet: hint.NewHintmap(),
		cp:               cp,
	}
}

func (opr *OperationProcessor) New(pool *storage.Statepool) prprocessor.OperationProcessor {
	nopr := operationProcessorPool.Get().(*OperationProcessor)
	nopr.id = util.UUID().String()

	if nopr.Logging == nil {
		nopr.Logging = logging.NewLogging(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "mitum-currency-operations-processor")
		})
		_ = nopr.SetLogging(opr.Logging)
	}

	if nopr.processorHintSet == nil {
		nopr.processorHintSet = opr.processorHintSet
	}

	if nopr.cp == nil {
		nopr.cp = opr.cp
	}

	nopr.pool = pool
	nopr.fee = map[currency.CurrencyID]currency.Big{}
	nopr.amountPool = map[string]currency.AmountState{}
	nopr.duplicated = map[string]DuplicationType{}
	nopr.duplicatedNewAddress = map[string]struct{}{}
	nopr.processorClosers = &sync.Map{}

	nopr.Log().Debug().Str("processor_id", nopr.id).Msg("new operation processors created")

	return nopr
}

func (opr *OperationProcessor) SetProcessor(
	hinter hint.Hinter,
	newProcessor currency.GetNewProcessor,
) (prprocessor.OperationProcessor, error) {
	if err := opr.processorHintSet.Add(hinter, newProcessor); err != nil {
		return nil, err
	}
	return opr, nil
}

func (opr *OperationProcessor) setState(op valuehash.Hash, sts ...state.State) error {
	opr.Lock()
	defer opr.Unlock()

	for i := range sts {
		if t, ok := sts[i].(currency.AmountState); ok {
			if t.Fee().OverZero() {
				f := currency.ZeroBig
				if i, found := opr.fee[t.Currency()]; found {
					f = i
				}

				opr.fee[t.Currency()] = f.Add(t.Fee())
			}
		}
	}

	return opr.pool.Set(op, sts...)
}

func (opr *OperationProcessor) PreProcess(op state.Processor) (state.Processor, error) {
	var sp state.Processor
	switch i, known, err := opr.getNewProcessor(op); {
	case err != nil:
		return nil, operation.NewBaseReasonErrorFromError(err)
	case !known:
		return op, nil
	default:
		sp = i
	}

	pop, err := sp.(state.PreProcessor).PreProcess(opr.pool.Get, opr.setState)
	if err != nil {
		return nil, err
	}

	if err := opr.checkDuplication(op); err != nil {
		return nil, operation.NewBaseReasonError("duplication found: %w", err)
	}

	return pop, nil
}

func (opr *OperationProcessor) Process(op state.Processor) error {
	switch op.(type) {
	case *extensioncurrency.TransfersProcessor,
		*extensioncurrency.CreateAccountsProcessor,
		*extensioncurrency.KeyUpdaterProcessor,
		*extensioncurrency.CurrencyRegisterProcessor,
		*extensioncurrency.CurrencyPolicyUpdaterProcessor,
		*extensioncurrency.SuffrageInflationProcessor,
		*extensioncurrency.CreateContractAccountsProcessor,
		*extensioncurrency.WithdrawsProcessor,
		*collection.DelegateProcessor,
		*collection.ApproveProcessor,
		*collection.CollectionRegisterProcessor,
		*collection.MintProcessor,
		*collection.TransferProcessor,
		*collection.BurnProcessor,
		*collection.SignProcessor,
		*BrokerRegisterProcessor,
		*PostProcessor:
		return opr.process(op)
	case currency.Transfers,
		currency.CreateAccounts,
		currency.KeyUpdater,
		extensioncurrency.CurrencyRegister,
		extensioncurrency.CurrencyPolicyUpdater,
		extensioncurrency.SuffrageInflation,
		extensioncurrency.CreateContractAccounts,
		extensioncurrency.Withdraws,
		collection.Delegate,
		collection.Approve,
		collection.CollectionRegister,
		collection.Mint,
		collection.Transfer,
		collection.Burn,
		collection.Sign,
		BrokerRegister,
		Post:
		pr, err := opr.PreProcess(op)
		if err != nil {
			return err
		}
		return opr.process(pr)
	default:
		return op.Process(opr.pool.Get, opr.pool.Set)
	}
}

func (opr *OperationProcessor) process(op state.Processor) error {
	var sp state.Processor

	switch t := op.(type) {
	case *extensioncurrency.TransfersProcessor:
		sp = t
	case *extensioncurrency.CreateAccountsProcessor:
		sp = t
	case *extensioncurrency.KeyUpdaterProcessor:
		sp = t
	case *extensioncurrency.CreateContractAccountsProcessor:
		sp = t
	case *extensioncurrency.WithdrawsProcessor:
		sp = t
	case *collection.ApproveProcessor:
		sp = t
	case *collection.CollectionRegisterProcessor:
		sp = t
	case *collection.MintProcessor:
		sp = t
	case *collection.TransferProcessor:
		sp = t
	case *collection.BurnProcessor:
		sp = t
	case *collection.SignProcessor:
		sp = t
	case *BrokerRegisterProcessor:
		sp = t
	case *PostProcessor:
		sp = t
	default:
		return op.Process(opr.pool.Get, opr.pool.Set)
	}

	return sp.Process(opr.pool.Get, opr.setState)
}

func (opr *OperationProcessor) checkDuplication(op state.Processor) error { // nolint:revive
	opr.Lock()
	defer opr.Unlock()

	var did string
	var didtype DuplicationType
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.Transfers:
		did = t.Fact().(currency.TransfersFact).Sender().String()
		didtype = DuplicationTypeSender
	case currency.CreateAccounts:
		fact := t.Fact().(currency.CreateAccountsFact)
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.KeyUpdater:
		did = t.Fact().(currency.KeyUpdaterFact).Target().String()
		didtype = DuplicationTypeSender
	case extensioncurrency.CurrencyRegister:
		did = t.Fact().(extensioncurrency.CurrencyRegisterFact).Currency().Currency().String()
		didtype = DuplicationTypeCurrency
	case extensioncurrency.CurrencyPolicyUpdater:
		did = t.Fact().(extensioncurrency.CurrencyPolicyUpdaterFact).Currency().String()
		didtype = DuplicationTypeCurrency
	case extensioncurrency.CreateContractAccounts:
		did = t.Fact().(extensioncurrency.CreateContractAccountsFact).Sender().String()
		didtype = DuplicationTypeSender
	case extensioncurrency.Withdraws:
		did = t.Fact().(extensioncurrency.WithdrawsFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Delegate:
		did = t.Fact().(collection.DelegateFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Approve:
		did = t.Fact().(collection.ApproveFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.CollectionRegister:
		did = t.Fact().(collection.CollectionRegisterFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Mint:
		did = t.Fact().(collection.MintFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Transfer:
		did = t.Fact().(collection.TransferFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Burn:
		did = t.Fact().(collection.BurnFact).Sender().String()
		didtype = DuplicationTypeSender
	case collection.Sign:
		did = t.Fact().(collection.SignFact).Sender().String()
		didtype = DuplicationTypeSender
	case BrokerRegister:
		did = t.Fact().(BrokerRegisterFact).Sender().String()
		didtype = DuplicationTypeSender
	case Post:
		did = t.Fact().(PostFact).Sender().String()
		didtype = DuplicationTypeSender
	default:
		return nil
	}

	if len(did) > 0 {
		if _, found := opr.duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicated currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.duplicated[did] = didtype
	}

	if len(newAddresses) > 0 {
		if err := opr.checkNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func (opr *OperationProcessor) checkNewAddressDuplication(as []base.Address) error {
	for i := range as {
		if _, found := opr.duplicatedNewAddress[as[i].String()]; found {
			return errors.Errorf("new address already processed")
		}
	}

	for i := range as {
		opr.duplicatedNewAddress[as[i].String()] = struct{}{}
	}

	return nil
}

func (opr *OperationProcessor) Close() error {
	opr.Lock()
	defer opr.Unlock()

	defer opr.close()

	if opr.cp != nil && len(opr.fee) > 0 {
		op := extensioncurrency.NewFeeOperation(extensioncurrency.NewFeeOperationFact(opr.pool.Height(), opr.fee))

		pr := extensioncurrency.NewFeeOperationProcessor(opr.cp, op)
		if err := pr.Process(opr.pool.Get, opr.pool.Set); err != nil {
			return err
		}
		opr.pool.AddOperations(op)
	}

	return nil
}

func (opr *OperationProcessor) Cancel() error {
	opr.Lock()
	defer opr.Unlock()

	defer opr.close()

	return nil
}

func (opr *OperationProcessor) getNewProcessor(op state.Processor) (state.Processor, bool, error) {
	switch i, err := opr.getNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.Transfers,
		currency.CreateAccounts,
		currency.KeyUpdater,
		extensioncurrency.CurrencyRegister,
		extensioncurrency.CurrencyPolicyUpdater,
		extensioncurrency.SuffrageInflation,
		extensioncurrency.CreateContractAccounts,
		extensioncurrency.Withdraws,
		collection.Delegate,
		collection.Approve,
		collection.CollectionRegister,
		collection.Mint,
		collection.Transfer,
		collection.Burn,
		collection.Sign,
		BrokerRegister,
		Post:

		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return op, false, nil
	}
}

func (opr *OperationProcessor) getNewProcessorFromHintset(op state.Processor) (state.Processor, error) {
	var f currency.GetNewProcessor
	if hinter, ok := op.(hint.Hinter); !ok {
		return nil, nil
	} else if i, err := opr.processorHintSet.Compatible(hinter); err != nil {
		if errors.Is(err, util.NotFoundError) {
			return nil, nil
		}

		return nil, err
	} else if j, ok := i.(currency.GetNewProcessor); !ok {
		return nil, errors.Errorf("invalid GetNewProcessor func, %q", i)
	} else {
		f = j
	}

	opp, err := f(op)
	if err != nil {
		return nil, err
	}

	h := op.(valuehash.Hasher).Hash().String()
	_, iscloser := opp.(io.Closer)
	if iscloser {
		opr.processorClosers.Store(h, opp)
		iscloser = true
	}

	opr.Log().Debug().
		Str("operation", h).
		Str("processor", fmt.Sprintf("%T", opp)).
		Bool("is_closer", iscloser).
		Msg("operation processor created")

	return opp, nil
}

func (opr *OperationProcessor) close() {
	opr.processorClosers.Range(func(_, v interface{}) bool {
		err := v.(io.Closer).Close()
		if err != nil {
			opr.Log().Error().Err(err).Str("op", fmt.Sprintf("%T", v)).Msg("failed to close operation processor")
		} else {
			opr.Log().Debug().Str("processor", fmt.Sprintf("%T", v)).Msg("operation processor closed")
		}

		return true
	})

	opr.pool = nil
	opr.fee = nil
	opr.amountPool = nil
	opr.duplicated = nil
	opr.duplicatedNewAddress = nil
	opr.processorClosers = nil

	operationProcessorPool.Put(opr)
	opr.Log().Debug().Str("processor_id", opr.id).Msg("operation processors closed")
}
