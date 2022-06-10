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

var (
	TradeFactType   = hint.Type("mitum-nft-trade-operation-fact")
	TradeFactHint   = hint.NewHint(TradeFactType, "v0.0.1")
	TradeFactHinter = TradeFact{BaseHinter: hint.NewBaseHinter(TradeFactHint)}
	TradeType       = hint.Type("mitum-nft-trade-operation")
	TradeHint       = hint.NewHint(TradeType, "v0.0.1")
	TradeHinter     = Trade{BaseOperation: operationHinter(TradeHint)}
)

type TradeFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	nft    nft.NFTID
	cid    currency.CurrencyID
}

func NewTradeFact(token []byte, sender base.Address, nft nft.NFTID, cid currency.CurrencyID) TradeFact {
	fact := TradeFact{
		BaseHinter: hint.NewBaseHinter(TradeFactHint),
		token:      token,
		sender:     sender,
		nft:        nft,
		cid:        cid,
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
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.nft.Bytes(),
		fact.cid.Bytes(),
	)
}

func (fact TradeFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for TradeFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.nft,
		fact.cid); err != nil {
		return err
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

func (fact TradeFact) NFT() nft.NFTID {
	return fact.nft
}

func (fact TradeFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact TradeFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact TradeFact) Rebuild() TradeFact {
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
