package broker

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	PostFactType   = hint.Type("mitum-nft-post-operation-fact")
	PostFactHint   = hint.NewHint(PostFactType, "v0.0.1")
	PostFactHinter = PostFact{BaseHinter: hint.NewBaseHinter(PostFactHint)}
	PostType       = hint.Type("mitum-nft-post-operation")
	PostHint       = hint.NewHint(PostType, "v0.0.1")
	PostHinter     = Post{BaseOperation: operationHinter(PostHint)}
)

var MaxPostItems = 10

type PostFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []PostItem
}

func NewPostFact(token []byte, sender base.Address, items []PostItem) PostFact {
	fact := PostFact{
		BaseHinter: hint.NewBaseHinter(PostFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact PostFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact PostFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact PostFact) Bytes() []byte {
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

func (fact PostFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items for PostFact")
	} else if n > MaxPostItems {
		return isvalid.InvalidError.Errorf("items over allowed; %d > %d", n, MaxPostItems)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	foundNFT := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		nft := fact.items[i].NFT()

		if err := nft.IsValid(nil); err != nil {
			return err
		}

		if _, found := foundNFT[nft.String()]; found {
			return isvalid.InvalidError.Errorf("duplicate nft found; %s", nft.String())
		}

		foundNFT[nft.String()] = true
	}

	return nil
}

func (fact PostFact) Token() []byte {
	return fact.token
}

func (fact PostFact) Sender() base.Address {
	return fact.sender
}

func (fact PostFact) Items() []PostItem {
	return fact.items
}

func (fact PostFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact PostFact) Rebuild() PostFact {
	items := make([]PostItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type Post struct {
	currency.BaseOperation
}

func NewPost(fact PostFact, fs []base.FactSign, memo string) (Post, error) {
	bo, err := currency.NewBaseOperationFromFact(PostHint, fact, fs, memo)
	if err != nil {
		return Post{}, err
	}

	return Post{BaseOperation: bo}, nil
}
