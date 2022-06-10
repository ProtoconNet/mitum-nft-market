package broker

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var MaxSettleAuctionItems = 10

type SettleAuctionItem interface {
	hint.Hinter
	isvalid.IsValider
	collection.NFTsItem
	Bytes() []byte
	Currency() currency.CurrencyID
	Rebuild() SettleAuctionItem
}

var (
	SettleAuctionFactType   = hint.Type("mitum-nft-market-settle-auction-operation-fact")
	SettleAuctionFactHint   = hint.NewHint(SettleAuctionFactType, "v0.0.1")
	SettleAuctionFactHinter = SettleAuctionFact{BaseHinter: hint.NewBaseHinter(SettleAuctionFactHint)}
	SettleAuctionType       = hint.Type("mitum-nft-market-settle-auction-operation")
	SettleAuctionHint       = hint.NewHint(SettleAuctionType, "v0.0.1")
	SettleAuctionHinter     = SettleAuction{BaseOperation: operationHinter(SettleAuctionHint)}
)

type SettleAuctionFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []SettleAuctionItem
}

func NewSettleAuctionFact(token []byte, sender base.Address, items []SettleAuctionItem) SettleAuctionFact {
	fact := SettleAuctionFact{
		BaseHinter: hint.NewBaseHinter(SettleAuctionFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact SettleAuctionFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact SettleAuctionFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SettleAuctionFact) Bytes() []byte {
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

func (fact SettleAuctionFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if l := len(fact.items); l < 1 {
		return isvalid.InvalidError.Errorf("empty items for SettleAuctionFact")
	} else if l > int(MaxSettleAuctionItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", l, MaxSettleAuctionItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nfts := fact.items[i].NFTs()

		for j := range nfts {
			if err := nfts[j].IsValid(nil); err != nil {
				return err
			}

			n := nfts[j]
			if _, found := foundNFT[n]; found {
				return isvalid.InvalidError.Errorf("duplicated nft found; %s", n)
			}

			foundNFT[n] = true
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact SettleAuctionFact) Token() []byte {
	return fact.token
}

func (fact SettleAuctionFact) Sender() base.Address {
	return fact.sender
}

func (fact SettleAuctionFact) Items() []SettleAuctionItem {
	return fact.items
}

func (fact SettleAuctionFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact SettleAuctionFact) Rebuild() SettleAuctionFact {
	items := make([]SettleAuctionItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type SettleAuction struct {
	currency.BaseOperation
}

func NewSettleAuction(fact SettleAuctionFact, fs []base.FactSign, memo string) (SettleAuction, error) {
	bo, err := currency.NewBaseOperationFromFact(SettleAuctionHint, fact, fs, memo)
	if err != nil {
		return SettleAuction{}, err
	}

	return SettleAuction{BaseOperation: bo}, nil
}
