package cmds

import (
	"bytes"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/seal"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

type CreateContractAccountCommand struct {
	*BaseCommand
	OperationFlags
	Sender    AddressFlag                       `arg:"" name:"sender" help:"sender address" required:"true"`
	Threshold uint                              `help:"threshold for keys (default: ${create_account_threshold})" default:"${create_account_threshold}"` // nolint
	Keys      []currencycmds.KeyFlag            `name:"key" help:"key for new account (ex: \"<public key>,<weight>\")" sep:"@"`
	Seal      mitumcmds.FileLoad                `help:"seal" optional:""`
	Amounts   []currencycmds.CurrencyAmountFlag `arg:"" name:"currency-amount" help:"amount (ex: \"<currency>,<amount>\")"`
	sender    base.Address
	keys      currency.BaseAccountKeys
}

func NewCreateContractAccountCommand() CreateContractAccountCommand {
	return CreateContractAccountCommand{
		BaseCommand: NewBaseCommand("create-account-operation"),
	}
}

func (cmd *CreateContractAccountCommand) Run(version util.Version) error { // nolint:dupl
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Wrap(err, "failed to initialize command")
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	sl, err := LoadSealAndAddOperation(
		cmd.Seal.Bytes(),
		cmd.Privatekey,
		cmd.NetworkID.NetworkID(),
		op,
	)
	if err != nil {
		return err
	}
	currencycmds.PrettyPrint(cmd.Out, cmd.Pretty, sl)

	return nil
}

func (cmd *CreateContractAccountCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = a

	if len(cmd.Keys) < 1 {
		return errors.Errorf("--key must be given at least one")
	}

	if len(cmd.Amounts) < 1 {
		return errors.Errorf("empty currency-amount, must be given at least one")
	}

	{
		ks := make([]currency.AccountKey, len(cmd.Keys))
		for i := range cmd.Keys {
			ks[i] = cmd.Keys[i].Key
		}

		if kys, err := currency.NewBaseAccountKeys(ks, cmd.Threshold); err != nil {
			return err
		} else if err := kys.IsValid(nil); err != nil {
			return err
		} else {
			cmd.keys = kys
		}
	}

	return nil
}

func (cmd *CreateContractAccountCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []extension.CreateContractAccountsItem
	for j := range i {
		if t, ok := i[j].(extension.CreateContractAccounts); ok {
			items = t.Fact().(extension.CreateContractAccountsFact).Items()
		}
	}

	ams := make([]currency.Amount, len(cmd.Amounts))
	for i := range cmd.Amounts {
		a := cmd.Amounts[i]
		am := currency.NewAmount(a.Big, a.CID)
		if err = am.IsValid(nil); err != nil {
			return nil, err
		}

		ams[i] = am
	}

	item := extension.NewCreateContractAccountsItemMultiAmounts(cmd.keys, ams)
	if err = item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := extension.NewCreateContractAccountsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := extension.NewCreateContractAccounts(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create create-contract-account operation")
	}
	return op, nil
}

func loadOperations(b []byte, networkID base.NetworkID) ([]operation.Operation, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, nil
	}

	var sl seal.Seal
	if s, err := LoadSeal(b, networkID); err != nil {
		return nil, err
	} else if so, ok := s.(operation.Seal); !ok {
		return nil, errors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return nil, errors.Errorf("seal is not operation.SealUpdater, %T", s)
	} else {
		sl = so
	}

	return sl.(operation.Seal).Operations(), nil
}

func LoadSeal(b []byte, networkID base.NetworkID) (seal.Seal, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, errors.Errorf("empty input")
	}

	var sl seal.Seal
	if err := encoder.Decode(b, jenc, &sl); err != nil {
		return nil, err
	}

	if err := sl.IsValid(networkID); err != nil {
		return nil, errors.Wrap(err, "invalid seal")
	}

	return sl, nil
}

func LoadSealAndAddOperation(
	b []byte,
	privatekey key.Privatekey,
	networkID base.NetworkID,
	op operation.Operation,
) (operation.Seal, error) {
	if b == nil {
		bs, err := operation.NewBaseSeal(
			privatekey,
			[]operation.Operation{op},
			networkID,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create operation.Seal")
		}
		return bs, nil
	}

	var sl operation.Seal
	if s, err := LoadSeal(b, networkID); err != nil {
		return nil, err
	} else if so, ok := s.(operation.Seal); !ok {
		return nil, errors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return nil, errors.Errorf("seal is not operation.SealUpdater, %T", s)
	} else {
		sl = so
	}

	// NOTE add operation to existing seal
	sl = sl.(operation.SealUpdater).SetOperations([]operation.Operation{op}).(operation.Seal)

	s, err := currencycmds.SignSeal(sl, privatekey, networkID)
	if err != nil {
		return nil, err
	}
	sl = s.(operation.Seal)

	return sl, nil
}
