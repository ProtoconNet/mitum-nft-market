package broker

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
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
	option  PostOption
	details PostDetails
}

func NewPostForm(option PostOption, details PostDetails) PostForm {
	return PostForm{
		BaseHinter: hint.NewBaseHinter(PostFormHint),
		option:     option,
		details:    details,
	}
}

func MustNewPostForm(option PostOption, details PostDetails) PostForm {
	form := NewPostForm(option, details)

	if err := form.IsValid(nil); err != nil {
		panic(err)
	}

	return form
}

func (form PostForm) Bytes() []byte {
	return util.ConcatBytesSlice(
		form.option.Bytes(),
		form.details.Bytes(),
	)
}

func (form PostForm) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		form.BaseHinter,
		form.option,
		form.details,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid PostForm; %w", err)
	}

	if form.option != form.details.Option() {
		return isvalid.InvalidError.Errorf("different option; %q != %q", form.option, form.details.Option())
	}

	return nil
}

func (form PostForm) Option() PostOption {
	return form.option
}

func (form PostForm) Details() PostDetails {
	return form.details
}

func (form PostForm) Rebuild() PostForm {
	return form
}

var (
	PostItemType   = hint.Type("mitum-nft-market-post-item")
	PostItemHint   = hint.NewHint(PostItemType, "v0.0.1")
	PostItemHinter = PostItem{
		BaseHinter: hint.NewBaseHinter(PostItemHint),
	}
)

type PostItem struct {
	hint.BaseHinter
	broker extensioncurrency.ContractID
	form   PostForm
	cid    currency.CurrencyID
}

func NewPostItem(broker extensioncurrency.ContractID, form PostForm, cid currency.CurrencyID) PostItem {
	return PostItem{
		BaseHinter: hint.NewBaseHinter(PostItemHint),
		broker:     broker,
		form:       form,
		cid:        cid,
	}
}

func (it PostItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.broker.Bytes(),
		it.form.Bytes(),
		it.cid.Bytes(),
	)
}

func (it PostItem) IsValid([]byte) error {
	return isvalid.Check(nil, false,
		it.BaseHinter,
		it.broker,
		it.form,
		it.cid)
}

func (it PostItem) Broker() extensioncurrency.ContractID {
	return it.broker
}

func (it PostItem) Form() PostForm {
	return it.form
}

func (it PostItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it PostItem) Rebuild() PostItem {
	return it
}
