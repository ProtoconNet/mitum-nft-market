package broker

import (
	"github.com/ProtoconNet/mitum-nft-market/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	UnpostFactType   = hint.Type("mitum-nft-unpost-operation-fact")
	UnpostFactHint   = hint.NewHint(UnpostFactType, "v0.0.1")
	UnpostFactHinter = UnpostFact{BaseHinter: hint.NewBaseHinter(UnpostFactHint)}
	UnpostType       = hint.Type("mitum-nft-unpost-operation")
	UnpostHint       = hint.NewHint(UnpostType, "v0.0.1")
	UnpostHinter     = Unpost{BaseOperation: operationHinter(UnpostHint)}
)

var MaxUnpostNFTs = 10

type UnpostFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	nfts   []nft.NFTID
	cid    currency.CurrencyID
}

func NewUnpostFact(token []byte, sender base.Address, nfts []nft.NFTID, cid currency.CurrencyID) UnpostFact {
	fact := UnpostFact{
		BaseHinter: hint.NewBaseHinter(UnpostFactHint),
		token:      token,
		sender:     sender,
		nfts:       nfts,
		cid:        cid,
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
	ns := make([][]byte, len(fact.nfts))

	for i := range fact.nfts {
		ns[i] = fact.nfts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.cid.Bytes(),
		util.ConcatBytesSlice(ns...),
	)
}

func (fact UnpostFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.token) < 1 {
		return isvalid.InvalidError.Errorf("empty token for UnpostFact")
	}

	if err := isvalid.Check(
		nil, false,
		fact.h,
		fact.sender,
		fact.cid); err != nil {
		return err
	}

	if n := len(fact.nfts); n < 1 {
		return isvalid.InvalidError.Errorf("empty nfts for UnpostFact")
	} else if n > MaxUnpostNFTs {
		return isvalid.InvalidError.Errorf("nfts over allowed; %d > %d", n, MaxUnpostNFTs)
	}

	foundNFT := map[string]bool{}
	for i := range fact.nfts {

		if err := fact.nfts[i].IsValid(nil); err != nil {
			return err
		}

		nft := fact.nfts[i].String()
		if _, found := foundNFT[nft]; found {
			return isvalid.InvalidError.Errorf("duplicate nft found; %s", nft)
		}

		foundNFT[nft] = true
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

func (fact UnpostFact) NFTs() []nft.NFTID {
	return fact.nfts
}

func (fact UnpostFact) Currency() currency.CurrencyID {
	return fact.cid
}

func (fact UnpostFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact UnpostFact) Rebuild() UnpostFact {
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
