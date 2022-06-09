package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	PostFormType   = hint.Type("mitum-nft-market-post-form")
	PostFormHint   = hint.NewHint(PostFormType, "v0.0.1")
	PostFormHinter = PostForm{BaseHinter: hint.NewBaseHinter(PostFormHint)}
)

type PostForm struct {
	hint.BaseHinter
	option    PostOption
	n         nft.NFTID
	closeTime PostCloseTime
	price     currency.Amount
}

func NewPostForm(option PostOption, n nft.NFTID, closeTime PostCloseTime, price currency.Amount) PostForm {
	return PostForm{
		BaseHinter: hint.NewBaseHinter(PostFormHint),
		option:     option,
		n:          n,
		closeTime:  closeTime,
		price:      price,
	}
}

func MustNewPostForm(broker extensioncurrency.ContractID, option PostOption, n nft.NFTID, closeTime PostCloseTime, price currency.Amount) PostForm {
	form := NewPostForm(option, n, closeTime, price)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form PostForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.option.Bytes(),
		form.n.Bytes(),
		form.closeTime.Bytes(),
		form.price.Bytes(),
	)
}

func (form PostForm) IsValid([]byte) error {
	if err := form.price.IsValid(nil); err != nil {
		return err
	} else if !form.price.Big().OverZero() {
		return isvalid.InvalidError.Errorf("price should be over zero")
	}

	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.option,
		form.n,
		form.closeTime,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid PostForm; %w", err)
	}

	return nil
}

func (form PostForm) Option() PostOption {
	return form.option
}

func (form PostForm) NFT() nft.NFTID {
	return form.n
}

func (form PostForm) CloseTime() PostCloseTime {
	return form.closeTime
}

func (form PostForm) Price() currency.Amount {
	return form.price
}

func (form PostForm) Rebuild() PostForm {
	return form
}

type BasePostItem struct {
	hint.BaseHinter
	broker extensioncurrency.ContractID
	forms  []PostForm
	cid    currency.CurrencyID
}

func NewBasePostItem(ht hint.Hint, broker extensioncurrency.ContractID, forms []PostForm, cid currency.CurrencyID) BasePostItem {
	return BasePostItem{
		BaseHinter: hint.NewBaseHinter(ht),
		broker:     broker,
		forms:      forms,
		cid:        cid,
	}
}

func (it BasePostItem) Bytes() []byte {
	bf := make([][]byte, len(it.forms))

	for i := range it.forms {
		bf[i] = it.forms[i].Bytes()
	}

	return util.ConcatBytesSlice(
		it.broker.Bytes(),
		it.cid.Bytes(),
		util.ConcatBytesSlice(bf...),
	)
}

func (it BasePostItem) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		it.BaseHinter,
		it.broker,
		it.cid); err != nil {
		return err
	}

	foundNFT := map[nft.NFTID]bool{}
	for i := range it.forms {
		if err := it.forms[i].IsValid(nil); err != nil {
			return err
		}

		if _, found := foundNFT[it.forms[i].NFT()]; found {
			return isvalid.InvalidError.Errorf("duplicate nft found; %q", it.forms[i].NFT())
		}

		foundNFT[it.forms[i].NFT()] = true
	}

	return nil
}

func (it BasePostItem) NFTs() []nft.NFTID {
	ns := make([]nft.NFTID, len(it.forms))

	for i := range it.forms {
		ns[i] = it.forms[i].NFT()
	}

	return ns
}

func (it BasePostItem) Broker() extensioncurrency.ContractID {
	return it.broker
}

func (it BasePostItem) Forms() []PostForm {
	return it.forms
}

func (it BasePostItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BasePostItem) Rebuild() PostItem {
	forms := []PostForm{}

	for i := range it.forms {
		forms = append(forms, it.forms[i].Rebuild())
	}
	it.forms = forms

	return it
}
