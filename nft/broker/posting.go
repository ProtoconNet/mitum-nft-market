package broker

import (
	"regexp"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type PostCloseTime string

var RevalidPostCloseTime = regexp.MustCompile(`^\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1-2]\d|3[0-1])T(?:[0-1]\d|2[0-3]):[0-5]\d:[0-5]\dZ$`)

func (pct PostCloseTime) Bytes() []byte {
	return []byte(pct)
}

func (pct PostCloseTime) String() string {
	return string(pct)
}

func (pct PostCloseTime) IsValid([]byte) error {
	if !RevalidPostCloseTime.Match([]byte(pct)) {
		return isvalid.InvalidError.Errorf("wrong post close time; %q", pct)
	}

	return nil
}

var (
	SellPostOption    = PostOption("sell")
	AuctionPostOption = PostOption("auction")
)

type PostOption string

func (po PostOption) Bytes() []byte {
	return []byte(po)
}

func (po PostOption) String() string {
	return string(po)
}

func (po PostOption) IsValid([]byte) error {
	if !(po == SellPostOption || po == AuctionPostOption) {
		return isvalid.InvalidError.Errorf("wrong post option; %q", po)
	}

	return nil
}

var (
	BiddingType   = hint.Type("mitum-nft-market-bidding")
	BiddingHint   = hint.NewHint(BiddingType, "v0.0.1")
	BiddingHinter = Bidding{BaseHinter: hint.NewBaseHinter(BiddingHint)}
)

type Bidding struct {
	hint.BaseHinter
	bidder base.Address
	amount currency.Amount
}

func NewBidding(bidder base.Address, amount currency.Amount) Bidding {
	return Bidding{
		bidder: bidder,
		amount: amount,
	}
}

func MustNewBidding(bidder base.Address, amount currency.Amount) Bidding {
	bidding := NewBidding(bidder, amount)

	if err := bidding.IsValid(nil); err != nil {
		panic(err)
	}

	return bidding
}

func (bidding Bidding) Bytes() []byte {
	return util.ConcatBytesSlice(
		bidding.bidder.Bytes(),
		bidding.amount.Bytes(),
	)
}

func (bidding Bidding) IsValid([]byte) error {

	if err := bidding.amount.IsValid(nil); err != nil {
		return err
	} else if !bidding.amount.Big().OverZero() {
		return isvalid.InvalidError.Errorf("amount should be over zero")
	}

	if err := bidding.bidder.IsValid(nil); err != nil {
		return isvalid.InvalidError.Errorf("invalid bidder for Bidding; %w", err)
	}
	return nil
}

func (bidding Bidding) Bidder() base.Address {
	return bidding.bidder
}

func (bidding Bidding) Amount() currency.Amount {
	return bidding.amount
}

var (
	PostingType   = hint.Type("mitum-nft-market-posting")
	PostingHint   = hint.NewHint(PostingType, "v0.0.1")
	PostingHinter = Posting{BaseHinter: hint.NewBaseHinter(PostingHint)}
)

type Posting struct {
	hint.BaseHinter
	active    bool
	broker    extensioncurrency.ContractID
	option    PostOption
	nft       nft.NFTID
	closeTime PostCloseTime
	price     currency.Amount
}

func NewPosting(active bool, broker extensioncurrency.ContractID, option PostOption, nft nft.NFTID, closeTime PostCloseTime, price currency.Amount) Posting {
	return Posting{
		broker:    broker,
		option:    option,
		nft:       nft,
		closeTime: closeTime,
		price:     price,
	}
}

func MustNewPosting(active bool, broker extensioncurrency.ContractID, option PostOption, nft nft.NFTID, closeTime PostCloseTime, price currency.Amount) Posting {
	posting := NewPosting(active, broker, option, nft, closeTime, price)

	if err := posting.IsValid(nil); err != nil {
		panic(err)
	}

	return posting
}

func (posting Posting) Bytes() []byte {
	ba := make([]byte, 1)
	if posting.active {
		ba[0] = 1
	} else {
		ba[0] = 0
	}

	return util.ConcatBytesSlice(
		ba,
		posting.broker.Bytes(),
		posting.option.Bytes(),
		posting.nft.Bytes(),
		posting.closeTime.Bytes(),
		posting.price.Bytes(),
	)
}

func (posting Posting) IsValid([]byte) error {
	if err := posting.price.IsValid(nil); err != nil {
		return err
	} else if !posting.price.Big().OverZero() {
		return isvalid.InvalidError.Errorf("price should be over zero")
	}

	if err := isvalid.Check(
		nil, false,
		posting.BaseHinter,
		posting.broker,
		posting.option,
		posting.nft,
		posting.closeTime,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid Posting; %w", err)
	}

	return nil
}

func (posting Posting) Active() bool {
	return posting.active
}

func (posting Posting) Broker() extensioncurrency.ContractID {
	return posting.broker
}

func (posting Posting) Option() PostOption {
	return posting.option
}

func (posting Posting) NFT() nft.NFTID {
	return posting.nft
}

func (posting Posting) CloseTime() PostCloseTime {
	return posting.closeTime
}

func (posting Posting) Price() currency.Amount {
	return posting.price
}

func (posting Posting) Rebuild() Posting {
	return posting
}
