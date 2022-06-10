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
	BidFactType   = hint.Type("mitum-nft-bid-operation-fact")
	BidFactHint   = hint.NewHint(BidFactType, "v0.0.1")
	BidFactHinter = BidFact{BaseHinter: hint.NewBaseHinter(BidFactHint)}
	BidType       = hint.Type("mitum-nft-bid-operation")
	BidHint       = hint.NewHint(BidType, "v0.0.1")
	BidHinter     = Bid{BaseOperation: operationHinter(BidHint)}
)

type BidFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	nft    nft.NFTID
	amount currency.Amount
}

func NewBidFact(token []byte, sender base.Address, nft nft.NFTID, amount currency.Amount) BidFact {
	fact := BidFact{
		BaseHinter: hint.NewBaseHinter(BidFactHint),
		token:      token,
		sender:     sender,
		nft:        nft,
		amount:     amount,
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
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.nft.Bytes(),
		fact.amount.Bytes(),
	)
}

func (fact BidFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for BidFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.nft,
		fact.amount); err != nil {
		return err
	}

	if !fact.amount.Big().OverZero() {
		return isvalid.InvalidError.Errorf("amount should be over zero")
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

func (fact BidFact) NFT() nft.NFTID {
	return fact.nft
}

func (fact BidFact) Amount() currency.Amount {
	return fact.amount
}

func (fact BidFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact BidFact) Rebuild() BidFact {
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
