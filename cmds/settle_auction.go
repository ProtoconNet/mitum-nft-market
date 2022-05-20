package cmds

import (
	"github.com/ProtoconNet/mitum-nft-market/nft"
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type SettleAuctionCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	NFT      NFTIDFlag                   `arg:"" name:"nft" help:"target nft for auction settlement (ex: \"<symbol>,<idx>\")" required:"true"`
	sender   base.Address
}

func NewSettleAuctionCommand() SettleAuctionCommand {
	return SettleAuctionCommand{
		BaseCommand: NewBaseCommand("settle-auction-operation"),
	}
}

func (cmd *SettleAuctionCommand) Run(version util.Version) error {
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

func (cmd *SettleAuctionCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	return nil

}

func (cmd *SettleAuctionCommand) createOperation() (operation.Operation, error) {
	fact := broker.NewSettleAuctionFact(
		[]byte(cmd.Token), cmd.sender,
		nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx), cmd.Currency.CID)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := broker.NewSettleAuction(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create settle-auction operation")
	}
	return op, nil
}
