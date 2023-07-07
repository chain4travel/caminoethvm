package commands

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	_ ExternalCommandVisitor = (*ExternalCommandVerifier)(nil)
)

type ExternalCommandVerifier struct {
}

func NewExternalCommandVerifier() (ExternalCommandVerifier, error) {
	return ExternalCommandVerifier{}, nil
}

func (ece ExternalCommandVerifier) ExecuteSetBaseFeeCommand(cmd *ExternalCommandSetBaseFee) error {

	switch {
	case cmd.NewBaseFee.Cmp(common.Big0) != 1:
		return errors.New("new base fee must be greater than zero")
	}

	return nil
}

func (ece ExternalCommandVerifier) ExecuteSetKYCStateCommand(cmd *ExternalCommandSetKYCState) error {
	set := make(map[common.Address]struct{})
	for _, update := range cmd.KYCUpdates {
		set[update.Address] = struct{}{}
	}

	if len(set) != len(cmd.KYCUpdates) {
		return errors.New("multiple updates not allowed")
	}

	return nil
}
