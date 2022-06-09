package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft-market/digest"
	"github.com/ProtoconNet/mitum-nft-market/nft/broker"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	Hinters []hint.Hinter
	Types   []hint.Type
)

var types = []hint.Type{
	currency.AccountType,
	currency.AddressType,
	currency.AmountType,
	currency.CreateAccountsFactType,
	currency.CreateAccountsItemMultiAmountsType,
	currency.CreateAccountsItemSingleAmountType,
	currency.CreateAccountsType,
	currency.AccountKeyType,
	currency.KeyUpdaterFactType,
	currency.KeyUpdaterType,
	currency.AccountKeysType,
	currency.TransfersFactType,
	currency.TransfersItemMultiAmountsType,
	currency.TransfersItemSingleAmountType,
	currency.TransfersType,
	extensioncurrency.CurrencyDesignType,
	extensioncurrency.CurrencyPolicyType,
	extensioncurrency.CurrencyPolicyUpdaterFactType,
	extensioncurrency.CurrencyPolicyUpdaterType,
	extensioncurrency.CurrencyRegisterFactType,
	extensioncurrency.CurrencyRegisterType,
	extensioncurrency.FeeOperationFactType,
	extensioncurrency.FeeOperationType,
	extensioncurrency.FixedFeeerType,
	extensioncurrency.GenesisCurrenciesFactType,
	extensioncurrency.GenesisCurrenciesType,
	extensioncurrency.NilFeeerType,
	extensioncurrency.RatioFeeerType,
	extensioncurrency.SuffrageInflationFactType,
	extensioncurrency.SuffrageInflationType,
	extensioncurrency.ContractAccountKeysType,
	extensioncurrency.ContractAccountType,
	extensioncurrency.CreateContractAccountsFactType,
	extensioncurrency.CreateContractAccountsType,
	extensioncurrency.CreateContractAccountsItemMultiAmountsType,
	extensioncurrency.CreateContractAccountsItemSingleAmountType,
	extensioncurrency.WithdrawsFactType,
	extensioncurrency.WithdrawsType,
	extensioncurrency.WithdrawsItemMultiAmountsType,
	extensioncurrency.WithdrawsItemSingleAmountType,
	nft.NFTIDType,
	nft.NFTType,
	nft.DesignType,
	collection.CollectionPolicyType,
	collection.MintFormType,
	collection.DelegateFactType,
	collection.DelegateType,
	collection.DelegateItemType,
	collection.ApproveFactType,
	collection.ApproveType,
	collection.ApproveItemMultiNFTsType,
	collection.ApproveItemSingleNFTType,
	collection.BurnFactType,
	collection.BurnType,
	collection.BurnItemMultiNFTsType,
	collection.BurnItemSingleNFTType,
	collection.CollectionRegisterFactType,
	collection.CollectionRegisterType,
	collection.CollectionRegisterFormType,
	collection.MintFactType,
	collection.MintType,
	collection.MintItemMultiNFTsType,
	collection.MintItemSingleNFTType,
	collection.TransferFactType,
	collection.TransferType,
	collection.TransferItemMultiNFTsType,
	collection.TransferItemSingleNFTType,
	broker.BrokerPolicyType,
	broker.BrokerRegisterFactType,
	broker.BrokerRegisterFormType,
	broker.BrokerRegisterType,
	broker.PostFormType,
	broker.PostFactType,
	broker.PostItemMultiNFTsType,
	broker.PostItemSingleNFTType,
	broker.PostType,
	broker.PostingType,
	broker.UnpostItemMultiNFTsType,
	broker.UnpostItemSingleNFTType,
	broker.UnpostFactType,
	broker.UnpostType,
	broker.TradeItemType,
	broker.TradeFactType,
	broker.TradeType,
	broker.BiddingType,
	broker.BidItemType,
	broker.BidFactType,
	broker.BidType,
	broker.SettleAuctionItemMultiNFTsType,
	broker.SettleAuctionItemSingleNFTType,
	broker.SettleAuctionFactType,
	broker.SettleAuctionType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	digest.OperationValueType,
}

var hinters = []hint.Hinter{
	currency.AccountHinter,
	currency.AddressHinter,
	currency.AmountHinter,
	currency.CreateAccountsFactHinter,
	currency.CreateAccountsItemMultiAmountsHinter,
	currency.CreateAccountsItemSingleAmountHinter,
	currency.CreateAccountsHinter,
	currency.KeyUpdaterFactHinter,
	currency.KeyUpdaterHinter,
	currency.AccountKeysHinter,
	currency.AccountKeyHinter,
	currency.TransfersFactHinter,
	currency.TransfersItemMultiAmountsHinter,
	currency.TransfersItemSingleAmountHinter,
	currency.TransfersHinter,
	extensioncurrency.CurrencyDesignHinter,
	extensioncurrency.CurrencyPolicyUpdaterFactHinter,
	extensioncurrency.CurrencyPolicyUpdaterHinter,
	extensioncurrency.CurrencyPolicyHinter,
	extensioncurrency.CurrencyRegisterFactHinter,
	extensioncurrency.CurrencyRegisterHinter,
	extensioncurrency.FeeOperationFactHinter,
	extensioncurrency.FeeOperationHinter,
	extensioncurrency.FixedFeeerHinter,
	extensioncurrency.GenesisCurrenciesFactHinter,
	extensioncurrency.GenesisCurrenciesHinter,
	extensioncurrency.NilFeeerHinter,
	extensioncurrency.RatioFeeerHinter,
	extensioncurrency.SuffrageInflationFactHinter,
	extensioncurrency.SuffrageInflationHinter,
	extensioncurrency.ContractAccountKeysHinter,
	extensioncurrency.ContractAccountHinter,
	extensioncurrency.CreateContractAccountsFactHinter,
	extensioncurrency.CreateContractAccountsHinter,
	extensioncurrency.CreateContractAccountsItemMultiAmountsHinter,
	extensioncurrency.CreateContractAccountsItemSingleAmountHinter,
	extensioncurrency.WithdrawsFactHinter,
	extensioncurrency.WithdrawsHinter,
	extensioncurrency.WithdrawsItemMultiAmountsHinter,
	extensioncurrency.WithdrawsItemSingleAmountHinter,
	nft.NFTIDHinter,
	nft.NFTHinter,
	nft.DesignHinter,
	collection.CollectionPolicyHinter,
	collection.MintFormHinter,
	collection.DelegateFactHinter,
	collection.DelegateHinter,
	collection.DelegateItemHinter,
	collection.ApproveFactHinter,
	collection.ApproveHinter,
	collection.ApproveItemMultiNFTsHinter,
	collection.ApproveItemSingleNFTHinter,
	collection.BurnFactHinter,
	collection.BurnHinter,
	collection.BurnItemMultiNFTsHinter,
	collection.BurnItemSingleNFTHinter,
	collection.CollectionRegisterFactHinter,
	collection.CollectionRegisterHinter,
	collection.CollectionRegisterFormHinter,
	collection.MintFactHinter,
	collection.MintHinter,
	collection.MintItemMultiNFTsHinter,
	collection.MintItemSingleNFTHinter,
	collection.TransferFactHinter,
	collection.TransferHinter,
	collection.TransferItemMultiNFTsHinter,
	collection.TransferItemSingleNFTHinter,
	broker.BrokerPolicyHinter,
	broker.BrokerRegisterFactHinter,
	broker.BrokerRegisterFormHinter,
	broker.BrokerRegisterHinter,
	broker.PostFormHinter,
	broker.PostFactHinter,
	broker.PostItemMultiNFTsHinter,
	broker.PostItemSingleNFTHinter,
	broker.PostHinter,
	broker.PostingHinter,
	broker.UnpostItemMultiNFTsHinter,
	broker.UnpostItemSingleNFTHinter,
	broker.UnpostFactHinter,
	broker.UnpostHinter,
	broker.TradeItemHinter,
	broker.TradeFactHinter,
	broker.TradeHinter,
	broker.BiddingHinter,
	broker.BidItemHinter,
	broker.BidFactHinter,
	broker.BidHinter,
	broker.SettleAuctionItemMultiNFTsHinter,
	broker.SettleAuctionItemSingleNFTHinter,
	broker.SettleAuctionFactHinter,
	broker.SettleAuctionHinter,
	digest.AccountValue{},
	digest.BaseHal{},
	digest.NodeInfo{},
	digest.OperationValue{},
	digest.Problem{},
}

func init() {
	Hinters = make([]hint.Hinter, len(launch.EncoderHinters)+len(hinters))
	copy(Hinters, launch.EncoderHinters)
	copy(Hinters[len(launch.EncoderHinters):], hinters)

	Types = make([]hint.Type, len(launch.EncoderTypes)+len(types))
	copy(Types, launch.EncoderTypes)
	copy(Types[len(launch.EncoderTypes):], types)
}
