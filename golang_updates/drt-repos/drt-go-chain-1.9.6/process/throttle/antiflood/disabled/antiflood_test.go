package disabled

import (
	"fmt"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestAntiFlood_ShouldNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		assert.Nil(t, r, "this shouldn't panic")
	}()

	daf := &AntiFlood{}
	assert.False(t, check.IfNil(daf))

	daf.SetMaxMessagesForTopic("test", 10)
	daf.ResetForTopic("test")
	daf.SetConsensusSizeNotifier(nil, 0)
	_ = daf.CanProcessMessagesOnTopic(core.PeerID(fmt.Sprint(1)), "test", 1, 0, nil)
	_ = daf.CanProcessMessage(nil, core.PeerID(fmt.Sprint(2)))
}
