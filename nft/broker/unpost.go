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

var MaxUnpostItems = 10

type UnpostItem interface {
	hint.Hinter
	isvalid.IsValider
	collection.NFTsItem
	Bytes() []byte
	Currency() currency.CurrencyID
	Rebuild() UnpostItem
}

var (
	UnpostFactType   = hint.Type("mitum-nft-market-unpost-operation-fact")
	UnpostFactHint   = hint.NewHint(UnpostFactType, "v0.0.1")
	UnpostFactHinter = UnpostFact{BaseHinter: hint.NewBaseHinter(UnpostFactHint)}
	UnpostType       = hint.Type("mitum-nft-market-unpost-operation")
	UnpostHint       = hint.NewHint(UnpostType, "v0.0.1")
	UnpostHinter     = Unpost{BaseOperation: operationHinter(UnpostHint)}
)

type UnpostFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []UnpostItem
}

func NewUnpostFact(token []byte, sender base.Address, items []UnpostItem) UnpostFact {
	fact := UnpostFact{
		BaseHinter: hint.NewBaseHinter(UnpostFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact UnpostFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact UnpostFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UnpostFact) Bytes() []byte {
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

func (fact UnpostFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for UnpostFact")
	} else if n > int(MaxUnpostItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxUnpostItems)
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

			nft := nfts[j]
			if _, found := foundNFT[nft]; found {
				return isvalid.InvalidError.Errorf("duplicated nft found; %s", nft)
			}

			foundNFT[nft] = true
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact UnpostFact) Token() []byte {
	return fact.token
}

func (fact UnpostFact) Sender() base.Address {
	return fact.sender
}

func (fact UnpostFact) Items() []UnpostItem {
	return fact.items
}

func (fact UnpostFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact UnpostFact) Rebuild() UnpostFact {
	items := make([]UnpostItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Unpost struct {
	currency.BaseOperation
}

func NewUnpost(fact UnpostFact, fs []base.FactSign, memo string) (Unpost, error) {
	bo, err := currency.NewBaseOperationFromFact(UnpostHint, fact, fs, memo)
	if err != nil {
		return Unpost{}, err
	}

	return Unpost{BaseOperation: bo}, nil
}
