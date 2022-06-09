package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxTradeItems = 10

var (
	TradeFactType   = hint.Type("mitum-nft-market-trade-operation-fact")
	TradeFactHint   = hint.NewHint(TradeFactType, "v0.0.1")
	TradeFactHinter = TradeFact{BaseHinter: hint.NewBaseHinter(TradeFactHint)}
	TradeType       = hint.Type("mitum-nft-market-trade-operation")
	TradeHint       = hint.NewHint(TradeType, "v0.0.1")
	TradeHinter     = Trade{BaseOperation: operationHinter(TradeHint)}
)

type TradeFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []TradeItem
}

func NewTradeFact(token []byte, sender base.Address, items []TradeItem) TradeFact {
	fact := TradeFact{
		BaseHinter: hint.NewBaseHinter(TradeFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact TradeFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact TradeFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TradeFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact TradeFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for TradeFact")
	} else if n > int(MaxTradeItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxTradeItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		n := fact.items[i].NFT()

		if err := n.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundNFT[n]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found; %s", n)
		}

		foundNFT[n] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact TradeFact) Token() []byte {
	return fact.token
}

func (fact TradeFact) Sender() base.Address {
	return fact.sender
}

func (fact TradeFact) Items() []TradeItem {
	return fact.items
}

func (fact TradeFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact TradeFact) Rebuild() TradeFact {
	items := make([]TradeItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Trade struct {
	currency.BaseOperation
}

func NewTrade(fact TradeFact, fs []base.FactSign, memo string) (Trade, error) {
	bo, err := currency.NewBaseOperationFromFact(TradeHint, fact, fs, memo)
	if err != nil {
		return Trade{}, err
	}

	return Trade{BaseOperation: bo}, nil
}
