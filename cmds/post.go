package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type PostCommand struct {
	*BaseCommand
	OperationFlags
	Sender    AddressFlag                     `arg:"" name:"sender" help:"sender address; nft owner or agent" required:"true"`
	Currency  currencycmds.CurrencyIDFlag     `arg:"" name:"currency" help:"currency id" required:"true"`
	BSymbol   string                          `arg:"" name:"broker" help:"broker symbol" required:"true"`
	NFT       NFTIDFlag                       `arg:"" name:"nft" help:"target nft to bid (ex: \"<symbol>,<idx>\")" required:"true"`
	Amount    currencycmds.CurrencyAmountFlag `arg:"" name:"price" help:"amount to bid (ex: \"<currency>,<amount>\")" required:"true"`
	CloseTime PostCloseTimeFlag               `name:"closetime" help:"post close time (ex: \"yyyy-MM-ddTHH:mm:ssZ\")" optional:""`
	Option    string                          `name:"option" help:"post option (sell|auction)" optional:""`
	sender    base.Address
	form      broker.PostForm
	broker    extensioncurrency.ContractID
}

func NewPostCommand() PostCommand {
	return PostCommand{
		BaseCommand: NewBaseCommand("post-operation"),
	}
}

func (cmd *PostCommand) Run(version util.Version) error {
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

func (cmd *PostCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	var option broker.PostOption
	if len(cmd.Option) < 1 {
		option = broker.SellPostOption
	} else if cmd.Option == broker.SellPostOption.String() {
		option = broker.SellPostOption
	} else if cmd.Option == broker.AuctionPostOption.String() {
		option = broker.AuctionPostOption
	} else {
		return errors.Errorf("wrong option; %q", cmd.Option)
	}
	if err := option.IsValid(nil); err != nil {
		return err
	}

	if option == broker.AuctionPostOption && len(cmd.CloseTime.s) < 1 {
		return errors.Errorf("empty post close time")
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}

	price := currency.NewAmount(cmd.Amount.Big, cmd.Amount.CID)
	if !price.Big().OverZero() {
		return errors.Errorf("price should be over zero")
	}

	var details broker.PostDetails
	if option == broker.SellPostOption {
		details = broker.NewSellDetails(n, price)
	} else {
		closeTime := broker.PostCloseTime(cmd.CloseTime.s)
		if err := closeTime.IsValid(nil); err != nil {
			return err
		}
		details = broker.NewAuctionDetails(n, closeTime, price)
	}

	form := broker.NewPostForm(option, details)
	if err := form.IsValid(nil); err != nil {
		return err
	}
	cmd.form = form

	symbol := extensioncurrency.ContractID(cmd.BSymbol)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}
	cmd.broker = symbol

	return nil
}

func (cmd *PostCommand) createOperation() (operation.Operation, error) {
	item := broker.NewPostItem(cmd.broker, cmd.form, cmd.Currency.CID)
	fact := broker.NewPostFact([]byte(cmd.Token), cmd.sender, []broker.PostItem{item})

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := broker.NewPost(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create post operation")
	}
	return op, nil
}
