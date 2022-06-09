package cmds

import (
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type BrokerRegisterCommand struct {
	*BaseCommand
	OperationFlags
	Sender    AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency  currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Target    AddressFlag                 `arg:"" name:"target" help:"target account to register polic" required:"true"`
	BSymbol   string                      `arg:"" name:"symbol" help:"broker symbol" required:"true"`
	Brokerage uint                        `arg:"" name:"brokerage" help:"brokerage parameter 0 <= brokerage param < 100" required:"true"`
	Receiver  AddressFlag                 `arg:"" name:"receiver" help:"brokerage receiver" required:"true"`
	Royalty   bool                        `name:"royalty" help:"--royalty true if broker supports royalty" optional:""`
	Uri       string                      `name:"uri" help:"--uri broker uri" optional:""`
	sender    base.Address
	target    base.Address
	form      broker.BrokerRegisterForm
}

func NewBrokerRegisterCommand() BrokerRegisterCommand {
	return BrokerRegisterCommand{
		BaseCommand: NewBaseCommand("broker-register-operation"),
	}
}

func (cmd *BrokerRegisterCommand) Run(version util.Version) error {
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

	bs, err := operation.NewBaseSeal(
		cmd.Privatekey,
		[]operation.Operation{op},
		cmd.NetworkID.NetworkID(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create operation.Seal")
	}
	PrettyPrint(cmd.Out, cmd.Pretty, bs)

	return nil
}

func (cmd *BrokerRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Target.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid target format; %q", cmd.Target.String())
	} else {
		cmd.target = a
	}

	receiver, err := cmd.Receiver.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format; %q", cmd.Receiver.String())
	}

	symbol := extensioncurrency.ContractID(cmd.BSymbol)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}

	brokerage := nft.PaymentParameter(cmd.Brokerage)
	if err := brokerage.IsValid(nil); err != nil {
		return err
	}

	uri := nft.URI(cmd.Uri)
	if err := uri.IsValid(nil); err != nil {
		return err
	}

	form := broker.NewBrokerRegisterForm(cmd.target, symbol, brokerage, receiver, cmd.Royalty, uri)
	if err := form.IsValid(nil); err != nil {
		return err
	}
	cmd.form = form

	return nil

}

func (cmd *BrokerRegisterCommand) createOperation() (operation.Operation, error) {
	fact := broker.NewBrokerRegisterFact([]byte(cmd.Token), cmd.sender, cmd.form, cmd.Currency.CID)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := broker.NewBrokerRegister(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create broker-register operation")
	}
	return op, nil
}
