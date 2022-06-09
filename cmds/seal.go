package cmds

import (
	extensioncmds "github.com/ProtoconNet/mitum-currency-extension/cmds"
	collectioncmds "github.com/ProtoconNet/mitum-nft/cmds"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
)

type SealCommand struct {
	Send                  SendCommand                                `cmd:"" name:"send" help:"send seal to remote mitum node"`
	CreateAccount         currencycmds.CreateAccountCommand          `cmd:"" name:"create-account" help:"create new account"`
	CreateContractAccount extensioncmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw              extensioncmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw contract account"`
	Transfer              currencycmds.TransferCommand               `cmd:"" name:"transfer" help:"transfer big"`
	KeyUpdater            currencycmds.KeyUpdaterCommand             `cmd:"" name:"key-updater" help:"update keys"`
	Delegate              collectioncmds.DelegateCommand             `cmd:"" name:"delegate" help:"delegate agent or cancel agent delegation"`
	Approve               collectioncmds.ApproveCommand              `cmd:"" name:"approve" help:"approve account for nft"`
	CollectionRegister    collectioncmds.CollectionRegisterCommand   `cmd:"" name:"collection-register" help:"register collection to contract account"`
	Mint                  collectioncmds.MintCommand                 `cmd:"" name:"mint" help:"mint nft to collection"`
	TransferNFTs          collectioncmds.TransferCommand             `cmd:"" name:"transfer-nfts" help:"transfer nfts"`
	Burn                  collectioncmds.BurnCommand                 `cmd:"" name:"burn" help:"burn nfts"`
	BrokerRegister        BrokerRegisterCommand                      `cmd:"" name:"broker-register" help:"register broker to contract account"`
	Post                  PostCommand                                `cmd:"" name:"post" help:"post nft on broker"`
	Unpost                UnpostCommand                              `cmd:"" name:"unpost" help:"unpost nft from broker"`
	Trade                 TradeCommand                               `cmd:"" name:"trade" help:"trade nft and ft tokens"`
	Bid                   BidCommand                                 `cmd:"" name:"bid" help:"bid on nft"`
	SettleAction          SettleAuctionCommand                       `cmd:"" name:"settle-auction" help:"settle posted auction"`
	CurrencyRegister      extensioncmds.CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater extensioncmds.CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`  // revive:disable-line:line-length-limit
	SuffrageInflation     currencycmds.SuffrageInflationCommand      `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"` // revive:disable-line:line-length-limit
	Sign                  currencycmds.SignSealCommand               `cmd:"" name:"sign" help:"sign seal"`
	SignFact              currencycmds.SignFactCommand               `cmd:"" name:"sign-fact" help:"sign facts of operation seal"`
}

func NewSealCommand() SealCommand {
	return SealCommand{
		Send:                  NewSendCommand(),
		CreateAccount:         currencycmds.NewCreateAccountCommand(),
		CreateContractAccount: extensioncmds.NewCreateContractAccountCommand(),
		Withdraw:              extensioncmds.NewWithdrawCommand(),
		Transfer:              currencycmds.NewTransferCommand(),
		KeyUpdater:            currencycmds.NewKeyUpdaterCommand(),
		Delegate:              collectioncmds.NewDelegateCommand(),
		Approve:               collectioncmds.NewApproveCommand(),
		CollectionRegister:    collectioncmds.NewCollectionRegisterCommand(),
		Mint:                  collectioncmds.NewMintCommand(),
		TransferNFTs:          collectioncmds.NewTransferCommand(),
		Burn:                  collectioncmds.NewBurnCommand(),
		BrokerRegister:        NewBrokerRegisterCommand(),
		Post:                  NewPostCommand(),
		Unpost:                NewUnpostCommand(),
		Trade:                 NewTradeCommand(),
		Bid:                   NewBidCommand(),
		SettleAction:          NewSettleAuctionCommand(),
		CurrencyRegister:      extensioncmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater: extensioncmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:     currencycmds.NewSuffrageInflationCommand(),
		Sign:                  currencycmds.NewSignSealCommand(),
		SignFact:              currencycmds.NewSignFactCommand(),
	}
}
