package cmds

import (
	"bytes"
	"context"

	"github.com/ProtoconNet/mitum-nft-market/digest"
	"github.com/pkg/errors"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/seal"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
	"github.com/spikeekips/mitum/launch/pm"
	"github.com/spikeekips/mitum/util/encoder"
)

type CleanByHeightStorageCommand struct {
	*mitumcmds.CleanByHeightStorageCommand
	*BaseNodeCommand
}

func newCleanByHeightStorageCommand() (CleanByHeightStorageCommand, error) {
	co := mitumcmds.NewCleanByHeightStorageCommand()
	cmd := CleanByHeightStorageCommand{
		CleanByHeightStorageCommand: &co,
		BaseNodeCommand:             NewBaseNodeCommand(co.Logging),
	}

	hooks := []pm.Hook{
		pm.NewHook(pm.HookPrefixPost, currencycmds.ProcessNameDigestDatabase,
			"set_digest_clean_storage_by_height", func(ctx context.Context) (context.Context, error) {
				var st *digest.Database
				if err := LoadDigestDatabaseContextValue(ctx, &st); err != nil {
					return ctx, err
				}

				return context.WithValue(ctx, mitumcmds.ContextValueCleanDatabaseByHeight,
					func(ctx context.Context, h base.Height) error {
						return st.CleanByHeight(ctx, h)
					}), nil
			}),
	}

	ps, err := cmd.BaseProcesses(co.Processes())
	if err != nil {
		return cmd, err
	}

	processes := []pm.Process{
		ProcessorDigestDatabase,
	}

	for i := range processes {
		if err := ps.AddProcess(processes[i], false); err != nil {
			return CleanByHeightStorageCommand{}, err
		}
	}

	for i := range hooks {
		if err := hooks[i].Add(ps); err != nil {
			return CleanByHeightStorageCommand{}, err
		}
	}

	_ = cmd.SetProcesses(ps)

	return cmd, nil
}

func LoadSeal(b []byte, networkID base.NetworkID) (seal.Seal, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, errors.Errorf("empty input")
	}

	var sl seal.Seal
	if err := encoder.Decode(b, jenc, &sl); err != nil {
		return nil, err
	}

	if err := sl.IsValid(networkID); err != nil {
		return nil, errors.Wrap(err, "invalid seal")
	}

	return sl, nil
}

func LoadSealAndAddOperation(
	b []byte,
	privatekey key.Privatekey,
	networkID base.NetworkID,
	op operation.Operation,
) (operation.Seal, error) {
	if b == nil {
		bs, err := operation.NewBaseSeal(
			privatekey,
			[]operation.Operation{op},
			networkID,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create operation.Seal")
		}
		return bs, nil
	}

	var sl operation.Seal
	if s, err := LoadSeal(b, networkID); err != nil {
		return nil, err
	} else if so, ok := s.(operation.Seal); !ok {
		return nil, errors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return nil, errors.Errorf("seal is not operation.SealUpdater, %T", s)
	} else {
		sl = so
	}

	// NOTE add operation to existing seal
	sl = sl.(operation.SealUpdater).SetOperations([]operation.Operation{op}).(operation.Seal)

	s, err := currencycmds.SignSeal(sl, privatekey, networkID)
	if err != nil {
		return nil, err
	}
	sl = s.(operation.Seal)

	return sl, nil
}
