package cmds

import (
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type BidCommand struct {
	*BaseCommand
	OperationFlags
	Sender AddressFlag                     `arg:"" name:"sender" help:"sender address" required:"true"`
	NFT    NFTIDFlag                       `arg:"" name:"nft" help:"target nft to bid (ex: \"<symbol>,<idx>\")" required:"true"`
	Amount currencycmds.CurrencyAmountFlag `arg:"" name:"amount" help:"amount to bid (ex: \"<currency>,<amount>\")" required:"true"`
	sender base.Address
	amount currency.Amount
	n      nft.NFTID
}

func NewBidCommand() BidCommand {
	return BidCommand{
		BaseCommand: NewBaseCommand("bid-operation"),
	}
}

func (cmd *BidCommand) Run(version util.Version) error {
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

func (cmd *BidCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	amount := currency.NewAmount(cmd.Amount.Big, cmd.Amount.CID)
	if err := amount.IsValid(nil); err != nil {
		return err
	}
	cmd.amount = amount

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	cmd.n = n

	return nil
}

func (cmd *BidCommand) createOperation() (operation.Operation, error) {
	item := broker.NewBidItem(cmd.n, cmd.amount)
	fact := broker.NewBidFact([]byte(cmd.Token), cmd.sender, []broker.BidItem{item})

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := broker.NewBid(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bid operation")
	}
	return op, nil
}
