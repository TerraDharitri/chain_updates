package iteratorChannelsProvider

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
)

const leavesChannelSize = 100

type userStateIteratorChannelsProvider struct {
}

// NewUserStateIteratorChannelsProvider creates a new instance of user state iterator channels provider
func NewUserStateIteratorChannelsProvider() *userStateIteratorChannelsProvider {
	return &userStateIteratorChannelsProvider{}
}

// GetIteratorChannels returns trie iterator channels
func (usicp *userStateIteratorChannelsProvider) GetIteratorChannels() *common.TrieIteratorChannels {
	return &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, leavesChannelSize),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (usicp *userStateIteratorChannelsProvider) IsInterfaceNil() bool {
	return usicp == nil
}
