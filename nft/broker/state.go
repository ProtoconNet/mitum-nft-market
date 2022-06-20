package broker

import (
	"fmt"
	"strings"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var (
	StateKeyBrokerPrefix = "broker:"
)

var (
	StateKeyPostingSuffix = ":posting"
)

func StateKeyBroker(id extensioncurrency.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyBrokerPrefix, id)
}

func IsStateBrokerKey(key string) bool {
	return strings.HasPrefix(key, StateKeyBrokerPrefix)
}

func StateBrokerValue(st state.State) (nft.Design, error) {
	value := st.Value()
	if value == nil {
		return nft.Design{}, util.NotFoundError.Errorf("design not found in State")
	}

	if design, ok := value.Interface().(nft.Design); !ok {
		return nft.Design{}, errors.Errorf("invalid design value found; %T", value.Interface())
	} else {
		return design, nil
	}
}

func SetStateBrokerValue(st state.State, design nft.Design) (state.State, error) {
	if vd, err := state.NewHintedValue(design); err != nil {
		return nil, err
	} else {
		return st.SetValue(vd)
	}
}

func StateKeyPosting(nid nft.NFTID) string {
	return fmt.Sprintf("%s%s", nid, StateKeyPostingSuffix)
}

func IsStatePostingKey(key string) bool {
	return strings.HasSuffix(key, StateKeyPostingSuffix)
}

func StatePostingValue(st state.State) (Posting, error) {
	value := st.Value()
	if value == nil {
		return Posting{}, util.NotFoundError.Errorf("posting not found in State")
	}

	if posting, ok := value.Interface().(Posting); !ok {
		return Posting{}, errors.Errorf("invalid posting value found; %T", value.Interface())
	} else {
		return posting, nil
	}
}

func SetStatePostingValue(st state.State, posting Posting) (state.State, error) {
	if vp, err := state.NewHintedValue(posting); err != nil {
		return nil, err
	} else {
		return st.SetValue(vp)
	}
}

func checkExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return nil
	default:
		return operation.NewBaseReasonError("state, %q already exists", key)
	}
}

func existsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, operation.NewBaseReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, operation.NewBaseReasonError("%s already exists", name)
	default:
		return st, nil
	}
}
