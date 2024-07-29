package ethadmin

import (
	"math/big"
	"testing"

	"github.com/ava-labs/coreth/core/admin"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestKycVerified(t *testing.T) {
	address := common.Address{1}
	sunriseTimestamp := uint64(1000)
	sunriseActiveHeader := &types.Header{
		Time: sunriseTimestamp,
	}

	tests := map[string]struct {
		stateDB        func(c *gomock.Controller) *admin.MockStateDB
		config         *params.ChainConfig
		header         *types.Header
		address        common.Address
		expectedResult bool
	}{
		"Not verified: Before Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(NOT_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: false,
		},
		"Not verified: After Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(NOT_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
				BerlinBlockTimestamp:        &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: false,
		},
		"KYC verified: Before Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYC_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: true,
		},
		"KYC verified: After Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYC_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
				BerlinBlockTimestamp:        &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: true,
		},
		"KYB verified: Before Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYB_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: false,
		},
		"KYB verified: After Berlin": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYB_VERIFIED)))
				return stateDB
			},
			config: &params.ChainConfig{
				SunrisePhase0BlockTimestamp: &sunriseTimestamp,
				BerlinBlockTimestamp:        &sunriseTimestamp,
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			adminCtrl := NewController(nil, tt.config)
			result := adminCtrl.KycVerified(tt.header, tt.stateDB(c), tt.address)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
