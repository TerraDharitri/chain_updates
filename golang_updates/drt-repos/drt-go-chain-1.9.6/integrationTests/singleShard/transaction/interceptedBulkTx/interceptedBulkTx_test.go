package interceptedBulkTx

import (
	"crypto"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/integrationTests"
	"github.com/TerraDharitri/drt-go-chain/process"
	"github.com/stretchr/testify/assert"
)

var stepDelay = time.Second / 10

func TestNode_GenerateSendInterceptBulkTransactionsWithMessenger(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	startingNonce := uint64(6)
	var numOfShards uint32 = 1
	var shardID uint32 = 0
	var txSignPrivKeyShardId uint32 = 0

	n := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            numOfShards,
		NodeShardId:          shardID,
		TxSignPrivKeyShardId: txSignPrivKeyShardId,
	})

	defer func() {
		n.Close()
	}()

	time.Sleep(integrationTests.P2pBootstrapDelay)

	//set the account's nonce to startingNonce
	_ = n.SetAccountNonce(startingNonce)
	noOfTx := common.MaxTxNonceDeltaAllowed

	time.Sleep(stepDelay)

	wg := sync.WaitGroup{}
	wg.Add(noOfTx)

	chanDone := make(chan bool)

	go func() {
		wg.Wait()
		chanDone <- true
	}()

	mut := sync.Mutex{}
	txHashes := make([][]byte, 0)
	transactions := make([]data.TransactionHandler, 0)

	//wire up handler
	n.DataPool.Transactions().RegisterOnAdded(func(key []byte, value interface{}) {
		mut.Lock()
		defer mut.Unlock()

		txHashes = append(txHashes, key)

		dataStore := n.DataPool.Transactions().ShardDataStore(
			process.ShardCacherIdentifier(n.ShardCoordinator.SelfId(), n.ShardCoordinator.SelfId()),
		)
		val, _ := dataStore.Get(key)

		if val == nil {
			assert.Fail(t, fmt.Sprintf("key %s not in store?", base64.StdEncoding.EncodeToString(key)))
			return
		}

		transactions = append(transactions, val.(*transaction.Transaction))
		wg.Done()
	})

	shardId := uint32(0)
	senderPrivateKeys := []crypto.PrivateKey{n.OwnAccount.SkTxSign}
	mintingValue := big.NewInt(10000000000)
	integrationTests.CreateMintingForSenders([]*integrationTests.TestProcessorNode{n}, shardId, senderPrivateKeys, mintingValue)

	encodedReceiverAddr, err := integrationTests.TestAddressPubkeyConverter.Encode(integrationTests.CreateRandomBytes(32))
	assert.Nil(t, err)

	err = n.Node.GenerateAndSendBulkTransactions(
		encodedReceiverAddr,
		big.NewInt(1),
		uint64(noOfTx),
		n.OwnAccount.SkTxSign,
		nil,
		integrationTests.ChainID,
		integrationTests.MinTransactionVersion,
	)

	assert.Nil(t, err)

	select {
	case <-chanDone:
	case <-time.After(time.Second * 60):
		assert.Fail(t, "timeout")
		return
	}

	integrationTests.CheckTxPresentAndRightNonce(
		t,
		startingNonce,
		noOfTx,
		txHashes,
		transactions,
		n.DataPool.Transactions(),
		n.ShardCoordinator,
	)
}

func TestNode_SendTransactionFromAnUnmintedAccountShouldReturnErrorAtApiLevel(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	var numOfShards uint32 = 1
	var shardID uint32 = 0
	var txSignPrivKeyShardId uint32 = 0

	node := integrationTests.NewTestProcessorNode(integrationTests.ArgTestProcessorNode{
		MaxShards:            numOfShards,
		NodeShardId:          shardID,
		TxSignPrivKeyShardId: txSignPrivKeyShardId,
	})

	defer func() {
		node.Close()
	}()

	tx := &transaction.Transaction{
		Nonce:    node.OwnAccount.Nonce,
		Value:    big.NewInt(1),
		SndAddr:  node.OwnAccount.Address,
		RcvAddr:  integrationTests.CreateRandomBytes(32),
		GasPrice: integrationTests.MinTxGasPrice,
		GasLimit: integrationTests.MinTxGasLimit,
		ChainID:  integrationTests.ChainID,
		Version:  integrationTests.MinTransactionVersion,
	}

	txBuff, _ := tx.GetDataForSigning(integrationTests.TestAddressPubkeyConverter, integrationTests.TestTxSignMarshalizer, integrationTests.TestTxSignHasher)
	tx.Signature, _ = node.OwnAccount.SingleSigner.Sign(node.OwnAccount.SkTxSign, txBuff)

	err := node.Node.ValidateTransaction(tx)
	assert.True(t, errors.Is(err, process.ErrInsufficientFunds))
}
