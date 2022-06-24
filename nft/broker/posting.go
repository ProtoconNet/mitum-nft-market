package broker

import (
	"regexp"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
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

type PostDetails interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	Option() PostOption
	NFT() nft.NFTID
	Price() currency.Amount
}

var (
	SellDetailsType   = hint.Type("mitum-nft-market-sell-details")
	SellDetailsHint   = hint.NewHint(SellDetailsType, "v0.0.1")
	SellDetailsHinter = SellDetails{BaseHinter: hint.NewBaseHinter(SellDetailsHint)}
)

type SellDetails struct {
	hint.BaseHinter
	nft   nft.NFTID
	price currency.Amount
}

func NewSellDetails(n nft.NFTID, price currency.Amount) SellDetails {
	return SellDetails{
		BaseHinter: hint.NewBaseHinter(SellDetailsHint),
		nft:        n,
		price:      price,
	}
}

func MustNewSellDetails(n nft.NFTID, price currency.Amount) SellDetails {
	details := NewSellDetails(n, price)

	if err := details.IsValid(nil); err != nil {
		panic(err)
	}

	return details
}

func (d SellDetails) Option() PostOption {
	return SellPostOption
}

func (d SellDetails) Bytes() []byte {
	return util.ConcatBytesSlice(d.nft.Bytes(), d.price.Bytes())
}

func (d SellDetails) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, d.nft, d.price); err != nil {
		return err
	}

	if !d.price.Big().OverZero() {
		return errors.Errorf("price must be over zero; %q", d.price.Big())
	}

	return nil
}

func (d SellDetails) NFT() nft.NFTID {
	return d.nft
}

func (d SellDetails) Price() currency.Amount {
	return d.price
}

var (
	AuctionDetailsType   = hint.Type("mitum-nft-market-auction-details")
	AuctionDetailsHint   = hint.NewHint(AuctionDetailsType, "v0.0.1")
	AuctionDetailsHinter = AuctionDetails{BaseHinter: hint.NewBaseHinter(AuctionDetailsHint)}
)

type AuctionDetails struct {
	hint.BaseHinter
	nft       nft.NFTID
	closeTime PostCloseTime
	price     currency.Amount
}

func NewAuctionDetails(n nft.NFTID, closeTime PostCloseTime, price currency.Amount) AuctionDetails {
	return AuctionDetails{
		BaseHinter: hint.NewBaseHinter(AuctionDetailsHint),
		nft:        n,
		closeTime:  closeTime,
		price:      price,
	}
}

func MustNewAuctionDetails(n nft.NFTID, closeTime PostCloseTime, price currency.Amount) AuctionDetails {
	details := NewAuctionDetails(n, closeTime, price)

	if err := details.IsValid(nil); err != nil {
		panic(err)
	}

	return details
}

func (d AuctionDetails) Option() PostOption {
	return AuctionPostOption
}

func (d AuctionDetails) Bytes() []byte {
	return util.ConcatBytesSlice(d.nft.Bytes(), d.closeTime.Bytes(), d.price.Bytes())
}

func (d AuctionDetails) IsValid([]byte) error {
	if err := isvalid.Check(nil, false, d.nft, d.closeTime, d.price); err != nil {
		return err
	}

	if !d.price.Big().OverZero() {
		return errors.Errorf("price must be over zero; %q", d.price.Big())
	}

	return nil
}

func (d AuctionDetails) CloseTime() PostCloseTime {
	return d.closeTime
}

func (d AuctionDetails) NFT() nft.NFTID {
	return d.nft
}

func (d AuctionDetails) Price() currency.Amount {
	return d.price
}

var (
	PostingType   = hint.Type("mitum-nft-market-posting")
	PostingHint   = hint.NewHint(PostingType, "v0.0.1")
	PostingHinter = Posting{BaseHinter: hint.NewBaseHinter(PostingHint)}
)

type Posting struct {
	hint.BaseHinter
	active  bool
	broker  extensioncurrency.ContractID
	option  PostOption
	details PostDetails
}

func NewPosting(active bool, broker extensioncurrency.ContractID, option PostOption, details PostDetails) Posting {
	return Posting{
		BaseHinter: hint.NewBaseHinter(PostingHint),
		active:     active,
		broker:     broker,
		option:     option,
		details:    details,
	}
}

func MustNewPosting(active bool, broker extensioncurrency.ContractID, option PostOption, details PostDetails) Posting {
	posting := NewPosting(active, broker, option, details)

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
		posting.details.Bytes(),
	)
}

func (posting Posting) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		posting.BaseHinter,
		posting.broker,
		posting.option,
		posting.details,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid Posting; %w", err)
	}

	if posting.option != posting.details.Option() {
		return isvalid.InvalidError.Errorf("different option; %q != %q", posting.option, posting.details.Option())
	}

	return nil
}

func (posting Posting) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(posting.Bytes())
}

func (posting Posting) Hash() valuehash.Hash {
	return posting.GenerateHash()
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

func (posting Posting) Details() PostDetails {
	return posting.details
}

func (posting Posting) Rebuild() Posting {
	return posting
}
