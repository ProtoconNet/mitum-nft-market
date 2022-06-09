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

var MaxBidItems = 10

var (
	BidFactType   = hint.Type("mitum-nft-market-bid-operation-fact")
	BidFactHint   = hint.NewHint(BidFactType, "v0.0.1")
	BidFactHinter = BidFact{BaseHinter: hint.NewBaseHinter(BidFactHint)}
	BidType       = hint.Type("mitum-nft-market-bid-operation")
	BidHint       = hint.NewHint(BidType, "v0.0.1")
	BidHinter     = Bid{BaseOperation: operationHinter(BidHint)}
)

type BidFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []BidItem
}

func NewBidFact(token []byte, sender base.Address, items []BidItem) BidFact {
	fact := BidFact{
		BaseHinter: hint.NewBaseHinter(BidFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact BidFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact BidFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact BidFact) Bytes() []byte {
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

func (fact BidFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for BidFact")
	} else if n > int(MaxBidItems) {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxBidItems)
	}

	if err := fact.sender.IsValid(nil); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nft := fact.items[i].NFT()
		if _, found := foundNFT[nft]; found {
			return isvalid.InvalidError.Errorf("duplicated nft found; %s", nft)
		}
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact BidFact) Token() []byte {
	return fact.token
}

func (fact BidFact) Sender() base.Address {
	return fact.sender
}

func (fact BidFact) Items() []BidItem {
	return fact.items
}

func (fact BidFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)
	as[0] = fact.sender
	return as, nil
}

func (fact BidFact) Rebuild() BidFact {
	items := make([]BidItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Bid struct {
	currency.BaseOperation
}

func NewBid(fact BidFact, fs []base.FactSign, memo string) (Bid, error) {
	bo, err := currency.NewBaseOperationFromFact(BidHint, fact, fs, memo)
	if err != nil {
		return Bid{}, err
	}

	return Bid{BaseOperation: bo}, nil
}
