package cmds

import (
	"github.com/ProtoconNet/mitum-nft-market/nft"
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
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
	BSymbol   string                          `arg:"" name:"symbol" help:"broker symbol" required:"true"`
	NFT       NFTIDFlag                       `arg:"" name:"nft" help:"target nft to bid (ex: \"<symbol>,<idx>\")" required:"true"`
	CloseTime PostCloseTimeFlag               `arg:"" name:"closetime" help:"post close time (ex: \"yyyy-MM-ddTHH:mm:ssZ\")"`
	Amount    currencycmds.CurrencyAmountFlag `arg:"" name:"price" help:"amount to bid (ex: \"<currency>,<amount>\")" required:"true"`
	Option    string                          `name:"option" help:"post option (sell|auction)" optional:""`
	sender    base.Address
	posting   broker.Posting
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
		option = broker.PostOption("")
	}
	if err := option.IsValid(nil); err != nil {
		return err
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}

	price := currency.NewAmount(cmd.Amount.Big, cmd.Amount.CID)
	if !price.Big().OverZero() {
		return errors.Errorf("price should be over zero")
	}

	posting := broker.NewPosting(
		nft.Symbol(cmd.BSymbol),
		option,
		n,
		broker.PostCloseTime(cmd.CloseTime.String()),
		price,
		[]broker.Bidding{},
	)
	cmd.posting = posting

	return nil
}

func (cmd *PostCommand) createOperation() (operation.Operation, error) {
	item := broker.NewPostItem(cmd.posting, cmd.Currency.CID)

	fact := broker.NewPostFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]broker.PostItem{item},
	)

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
