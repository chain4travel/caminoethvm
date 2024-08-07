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
	sunriseTimestampBig := big.NewInt(0).SetUint64(sunriseTimestamp)
	sunriseActiveHeader := &types.Header{
		Time: sunriseTimestamp,
	}

	adminCtrl := NewController(nil, &params.ChainConfig{
		SunrisePhase0BlockTimestamp: sunriseTimestampBig,
	})

	tests := map[string]struct {
		stateDB        func(c *gomock.Controller) *admin.MockStateDB
		header         *types.Header
		address        common.Address
		expectedResult bool
	}{
		"Not verified": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(NOT_VERIFIED)))
				return stateDB
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: false,
		},
		"KYC verified": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYC_VERIFIED)))
				return stateDB
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: true,
		},
		"KYB verified": {
			stateDB: func(c *gomock.Controller) *admin.MockStateDB {
				stateDB := admin.NewMockStateDB(c)
				stateDB.EXPECT().GetState(contractAddr, kycStoragePosition(address)).
					Return(common.BigToHash(big.NewInt(KYB_VERIFIED)))
				return stateDB
			},
			header:         sunriseActiveHeader,
			address:        address,
			expectedResult: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := gomock.NewController(t)
			result := adminCtrl.KycVerified(tt.header, tt.stateDB(c), tt.address)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
