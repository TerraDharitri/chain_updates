package blockAPI

import (
	"bytes"
	"encoding/hex"
	"errors"

	"math/big"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data"
	"github.com/TerraDharitri/drt-go-chain-core/data/api"
	"github.com/TerraDharitri/drt-go-chain-core/data/block"
	"github.com/TerraDharitri/drt-go-chain-core/data/receipt"
	"github.com/TerraDharitri/drt-go-chain-core/data/smartContractResult"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/holders"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/node/mock"
	"github.com/TerraDharitri/drt-go-chain/storage"
	"github.com/TerraDharitri/drt-go-chain/testscommon"
	dataRetrieverTestsCommon "github.com/TerraDharitri/drt-go-chain/testscommon/dataRetriever"
	"github.com/TerraDharitri/drt-go-chain/testscommon/dblookupext"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	"github.com/TerraDharitri/drt-go-chain/testscommon/genericMocks"
	"github.com/TerraDharitri/drt-go-chain/testscommon/hashingMocks"
	storageMocks "github.com/TerraDharitri/drt-go-chain/testscommon/storage"
	"github.com/stretchr/testify/require"
)

func createBaseBlockProcessor() *baseAPIBlockProcessor {
	return &baseAPIBlockProcessor{
		hasDbLookupExtensions:    true,
		selfShardID:              0,
		emptyReceiptsHash:        nil,
		store:                    &storageMocks.ChainStorerStub{},
		marshalizer:              &mock.MarshalizerFake{},
		uint64ByteSliceConverter: mock.NewNonceHashConverterMock(),
		historyRepo:              &dblookupext.HistoryRepositoryStub{},
		hasher:                   &hashingMocks.HasherMock{},
		addressPubKeyConverter:   testscommon.NewPubkeyConverterMock(32),
		txStatusComputer:         &mock.StatusComputerStub{},
		apiTransactionHandler:    &mock.TransactionAPIHandlerStub{},
		logsFacade:               &testscommon.LogsFacadeStub{},
		receiptsRepository:       &testscommon.ReceiptsRepositoryStub{},
		enableEpochsHandler:      &enableEpochsHandlerMock.EnableEpochsHandlerStub{},
	}
}

func TestBaseBlockGetIntraMiniblocksSCRS(t *testing.T) {
	t.Parallel()

	baseAPIBlockProc := createBaseBlockProcessor()

	scrHash := []byte("scr1")
	miniblock := &block.MiniBlock{
		Type:     block.SmartContractResultBlock,
		TxHashes: [][]byte{scrHash},
	}

	scResult := &smartContractResult.SmartContractResult{
		SndAddr: []byte("snd"),
		RcvAddr: []byte("rcv"),
		Data:    []byte("doSomething"),
	}
	scResultBytes, _ := baseAPIBlockProc.marshalizer.Marshal(scResult)

	baseAPIBlockProc.store = genericMocks.NewChainStorerMock(0)
	storer, _ := baseAPIBlockProc.store.GetStorer(dataRetriever.UnsignedTransactionUnit)
	_ = storer.Put(scrHash, scResultBytes)

	baseAPIBlockProc.receiptsRepository = &testscommon.ReceiptsRepositoryStub{
		LoadReceiptsCalled: func(header data.HeaderHandler, headerHash []byte) (common.ReceiptsHolder, error) {
			return holders.NewReceiptsHolder([]*block.MiniBlock{miniblock}), nil
		},
	}

	baseAPIBlockProc.apiTransactionHandler = &mock.TransactionAPIHandlerStub{
		UnmarshalTransactionCalled: func(txBytes []byte, txType transaction.TxType) (*transaction.ApiTransactionResult, error) {
			return &transaction.ApiTransactionResult{
				Sender:   hex.EncodeToString(scResult.SndAddr),
				Receiver: hex.EncodeToString(scResult.RcvAddr),
				Data:     scResult.Data,
			}, nil
		},
	}

	blockHeader := &block.Header{ReceiptsHash: []byte("aaaa"), Epoch: 0}
	intraMbs, err := baseAPIBlockProc.getIntrashardMiniblocksFromReceiptsStorage(blockHeader, []byte{}, api.BlockQueryOptions{WithTransactions: true})
	require.Nil(t, err)
	require.Equal(t, &api.MiniBlock{
		Hash: "f4add7b23eb83cf290422b0f6b770e3007b8ed3cd9683797fc90c8b4881f27bd",
		Type: "SmartContractResultBlock",
		Transactions: []*transaction.ApiTransactionResult{
			{
				Hash:          "73637231",
				HashBytes:     []byte{0x73, 0x63, 0x72, 0x31},
				Sender:        "736e64",
				Receiver:      "726376",
				Data:          []byte("doSomething"),
				MiniBlockType: "SmartContractResultBlock",
				MiniBlockHash: "f4add7b23eb83cf290422b0f6b770e3007b8ed3cd9683797fc90c8b4881f27bd",
			},
		},
		ProcessingType:        block.Normal.String(),
		IsFromReceiptsStorage: true,
	}, intraMbs[0])
}

func TestBaseBlockGetIntraMiniblocksReceipts(t *testing.T) {
	t.Parallel()

	baseAPIBlockProc := createBaseBlockProcessor()

	receiptHash := []byte("rec1")
	miniblock := &block.MiniBlock{
		Type:     block.ReceiptBlock,
		TxHashes: [][]byte{receiptHash},
	}

	receiptObj := &receipt.Receipt{
		Value:   big.NewInt(1000),
		SndAddr: []byte("sndAddr"),
		Data:    []byte("refund"),
		TxHash:  []byte("hash"),
	}
	receiptBytes, _ := baseAPIBlockProc.marshalizer.Marshal(receiptObj)

	baseAPIBlockProc.store = genericMocks.NewChainStorerMock(0)
	storer, _ := baseAPIBlockProc.store.GetStorer(dataRetriever.UnsignedTransactionUnit)
	_ = storer.Put(receiptHash, receiptBytes)

	baseAPIBlockProc.receiptsRepository = &testscommon.ReceiptsRepositoryStub{
		LoadReceiptsCalled: func(header data.HeaderHandler, headerHash []byte) (common.ReceiptsHolder, error) {
			return holders.NewReceiptsHolder([]*block.MiniBlock{miniblock}), nil
		},
	}

	baseAPIBlockProc.apiTransactionHandler = &mock.TransactionAPIHandlerStub{
		UnmarshalReceiptCalled: func(receiptBytes []byte) (*transaction.ApiReceipt, error) {
			encodedSndAddrReceiptObj, err := baseAPIBlockProc.addressPubKeyConverter.Encode(receiptObj.SndAddr)
			require.NoError(t, err)

			return &transaction.ApiReceipt{
				Value:   receiptObj.Value,
				SndAddr: encodedSndAddrReceiptObj,
				Data:    string(receiptObj.Data),
				TxHash:  hex.EncodeToString(receiptObj.TxHash),
			}, nil
		},
	}

	blockHeader := &block.Header{ReceiptsHash: []byte("aaaa"), Epoch: 0}
	intraMbs, err := baseAPIBlockProc.getIntrashardMiniblocksFromReceiptsStorage(blockHeader, []byte{}, api.BlockQueryOptions{WithTransactions: true})
	require.Nil(t, err)
	require.Equal(t, &api.MiniBlock{
		Hash: "596545f64319f2fcf8e0ebae06f40f3353d603f6070255588a48018c7b30c951",
		Type: "ReceiptBlock",
		Receipts: []*transaction.ApiReceipt{
			{
				SndAddr: "736e6441646472",
				Data:    "refund",
				TxHash:  "68617368",
				Value:   big.NewInt(1000),
			},
		},
		ProcessingType:        block.Normal.String(),
		IsFromReceiptsStorage: true,
	}, intraMbs[0])
}

func TestBaseBlock_getAndAttachTxsToMb_MiniblockTxBlock(t *testing.T) {
	t.Parallel()

	baseAPIBlockProc := createBaseBlockProcessor()

	txHash := []byte("tx1")
	txMb := &block.MiniBlock{
		Type:     block.TxBlock,
		TxHashes: [][]byte{txHash},
	}
	txMbBytes, _ := baseAPIBlockProc.marshalizer.Marshal(txMb)

	mbStorer := genericMocks.NewStorerMock()
	mbHash := []byte("mbHash")
	_ = mbStorer.Put(mbHash, txMbBytes)

	unsignedStorer := genericMocks.NewStorerMock()
	tx := &transaction.Transaction{
		Value:   big.NewInt(1000),
		SndAddr: []byte("sndAddr"),
		RcvAddr: []byte("rcvAddr"),
		Data:    []byte("refund"),
		Nonce:   1,
	}
	txBytes, _ := baseAPIBlockProc.marshalizer.Marshal(tx)
	_ = unsignedStorer.Put(txHash, txBytes)

	baseAPIBlockProc.store = &storageMocks.ChainStorerStub{
		GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
			switch unitType {
			case dataRetriever.MiniBlockUnit:
				return mbStorer, nil
			case dataRetriever.TransactionUnit:
				return unsignedStorer, nil
			}

			return nil, storage.ErrKeyNotFound
		},
	}

	baseAPIBlockProc.apiTransactionHandler = &mock.TransactionAPIHandlerStub{
		UnmarshalTransactionCalled: func(txBytes []byte, txType transaction.TxType) (*transaction.ApiTransactionResult, error) {
			return &transaction.ApiTransactionResult{
				Sender:   hex.EncodeToString(tx.SndAddr),
				Receiver: hex.EncodeToString(tx.RcvAddr),
				Data:     tx.Data,
				Nonce:    tx.Nonce,
			}, nil
		},
	}

	mbhr := &block.MiniBlockHeaderReserved{
		IndexOfFirstTxProcessed: 0,
		IndexOfLastTxProcessed:  1,
	}
	marshalizer := testscommon.ProtobufMarshalizerMock{}
	mbhrBytes, _ := marshalizer.Marshal(mbhr)

	mbHeader := &block.MiniBlockHeader{
		Hash:     mbHash,
		Reserved: mbhrBytes,
	}

	testHeader := &block.Header{
		Epoch: 0,
		Round: 37,
	}
	apiMB := &api.MiniBlock{}
	err := baseAPIBlockProc.getAndAttachTxsToMb(mbHeader, testHeader, apiMB, api.BlockQueryOptions{})
	require.Nil(t, err)
	require.Equal(t, &api.MiniBlock{
		Transactions: []*transaction.ApiTransactionResult{
			{
				Nonce:         1,
				Hash:          "747831",
				HashBytes:     []byte{0x74, 0x78, 0x31},
				Sender:        "736e6441646472",
				Receiver:      "72637641646472",
				Data:          []byte("refund"),
				MiniBlockType: "TxBlock",
				MiniBlockHash: "6d6248617368",
				Epoch:         testHeader.GetEpoch(),
				Round:         testHeader.GetRound(),
			},
		},
	}, apiMB)
}

func TestBaseBlock_getAndAttachTxsToMbShouldIncludeLogsAsSpecified(t *testing.T) {
	t.Parallel()

	testHeader := &block.Header{
		Epoch: 7,
		Round: 140,
	}

	marshalizer := &marshal.GogoProtoMarshalizer{}

	storageService := genericMocks.NewChainStorerMock(testHeader.GetEpoch())
	processor := createBaseBlockProcessor()
	processor.marshalizer = marshalizer
	processor.store = storageService

	// Set up a dummy transformer for "txBytes" -> "ApiTransactionResult" (only "Nonce" is handled)
	processor.apiTransactionHandler = &mock.TransactionAPIHandlerStub{
		UnmarshalTransactionCalled: func(txBytes []byte, txType transaction.TxType) (*transaction.ApiTransactionResult, error) {
			tx := &transaction.Transaction{}
			err := marshalizer.Unmarshal(tx, txBytes)
			if err != nil {
				return nil, err
			}

			return &transaction.ApiTransactionResult{Nonce: tx.Nonce}, nil
		},
	}

	// Set up a miniblock
	miniblockHash := []byte{0xff}
	miniblock := &block.MiniBlock{
		Type:     block.TxBlock,
		TxHashes: [][]byte{{0xaa}, {0xbb}, {0xcc}},
	}
	miniblockBytes, _ := processor.marshalizer.Marshal(miniblock)
	_ = storageService.Miniblocks.Put(miniblockHash, miniblockBytes)

	// Setup some transactions
	firstTx := &transaction.Transaction{Nonce: 42}
	secondTx := &transaction.Transaction{Nonce: 43}
	thirdTx := &transaction.Transaction{Nonce: 44}

	firstTxBytes, _ := marshalizer.Marshal(firstTx)
	secondTxBytes, _ := marshalizer.Marshal(secondTx)
	thirdTxBytes, _ := marshalizer.Marshal(thirdTx)

	_ = storageService.Transactions.Put([]byte{0xaa}, firstTxBytes)
	_ = storageService.Transactions.Put([]byte{0xbb}, secondTxBytes)
	_ = storageService.Transactions.Put([]byte{0xcc}, thirdTxBytes)

	// Setup some logs for 1st and 3rd transactions (none for 2nd)
	processor.logsFacade = &testscommon.LogsFacadeStub{
		IncludeLogsInTransactionsCalled: func(txs []*transaction.ApiTransactionResult, logsKeys [][]byte, epoch uint32) error {
			// Check the input arguments to match our scenario
			if len(txs) != 3 || len(logsKeys) != 3 {
				return nil
			}
			if !bytes.Equal(logsKeys[0], []byte{0xaa}) ||
				!bytes.Equal(logsKeys[1], []byte{0xbb}) ||
				!bytes.Equal(logsKeys[2], []byte{0xcc}) {
				return nil
			}
			if epoch != testHeader.GetEpoch() {
				return nil
			}

			txs[0].Logs = &transaction.ApiLogs{
				Events: []*transaction.Events{
					{Identifier: "first"},
				},
			}

			txs[2].Logs = &transaction.ApiLogs{
				Events: []*transaction.Events{
					{Identifier: "third"},
				},
			}

			return nil
		},
	}

	// Now let's test the loading of transaction and logs
	miniblockHeader := &block.MiniBlockHeader{Hash: miniblockHash}
	miniblockOnApi := &api.MiniBlock{}
	err := processor.getAndAttachTxsToMb(miniblockHeader, testHeader, miniblockOnApi, api.BlockQueryOptions{WithLogs: true})

	require.Nil(t, err)
	require.Len(t, miniblockOnApi.Transactions, 3)
	require.Equal(t, uint64(42), miniblockOnApi.Transactions[0].Nonce)
	require.Equal(t, uint64(43), miniblockOnApi.Transactions[1].Nonce)
	require.Equal(t, uint64(44), miniblockOnApi.Transactions[2].Nonce)
	require.Equal(t, "first", miniblockOnApi.Transactions[0].Logs.Events[0].Identifier)
	require.Nil(t, miniblockOnApi.Transactions[1].Logs)
	require.Equal(t, "third", miniblockOnApi.Transactions[2].Logs.Events[0].Identifier)
	for _, tx := range miniblockOnApi.Transactions {
		require.Equal(t, testHeader.GetRound(), tx.Round)
		require.Equal(t, testHeader.GetEpoch(), tx.Epoch)
	}
}

func TestExtractExecutedTxHashes(t *testing.T) {
	t.Parallel()

	array := make([][]byte, 10)
	res := extractExecutedTxHashes(array, 0, int32(len(array))-1)
	require.Len(t, res, 10)

	res = extractExecutedTxHashes(array, 0, int32(len(array)))
	require.Equal(t, res, array)

	res = extractExecutedTxHashes(array, -1, int32(len(array)))
	require.Equal(t, res, array)

	res = extractExecutedTxHashes(array, 20, int32(len(array)))
	require.Equal(t, res, array)

	res = extractExecutedTxHashes(array, 0, int32(len(array))+1)
	require.Equal(t, res, array)

	array = make([][]byte, 0, 10)
	for idx := 0; idx < 10; idx++ {
		array = append(array, []byte{byte(idx)})
	}
	res = extractExecutedTxHashes(array, 0, 5)
	require.Equal(t, res, [][]byte{{byte(0)}, {byte(1)}, {byte(2)}, {byte(3)}, {byte(4)}, {byte(5)}})
}

func TestAddScheduledInfoInBlock(t *testing.T) {
	t.Parallel()

	blockHeader := &block.HeaderV2{
		ScheduledRootHash:        []byte("srh"),
		ScheduledAccumulatedFees: big.NewInt(1000),
		ScheduledDeveloperFees:   big.NewInt(500),
		ScheduledGasProvided:     10,
		ScheduledGasPenalized:    1,
		ScheduledGasRefunded:     2,
	}

	apiBlock := &api.Block{}

	addScheduledInfoInBlock(blockHeader, apiBlock)
	require.Equal(t, &api.Block{
		ScheduledData: &api.ScheduledData{
			ScheduledRootHash:        "737268",
			ScheduledAccumulatedFees: "1000",
			ScheduledDeveloperFees:   "500",
			ScheduledGasProvided:     10,
			ScheduledGasPenalized:    1,
			ScheduledGasRefunded:     2,
		},
	}, apiBlock)
}

func TestProofToAPIProof(t *testing.T) {
	t.Parallel()

	headerProof := &block.HeaderProof{
		PubKeysBitmap:       []byte("bitmap"),
		AggregatedSignature: []byte("sig"),
		HeaderHash:          []byte("hash"),
		HeaderEpoch:         1,
		HeaderNonce:         3,
		HeaderShardId:       2,
		HeaderRound:         4,
		IsStartOfEpoch:      true,
	}

	proofToAPIProof(headerProof)
	require.Equal(t, &api.HeaderProof{
		PubKeysBitmap:       hex.EncodeToString(headerProof.PubKeysBitmap),
		AggregatedSignature: hex.EncodeToString(headerProof.AggregatedSignature),
		HeaderHash:          hex.EncodeToString(headerProof.HeaderHash),
		HeaderEpoch:         1,
		HeaderNonce:         3,
		HeaderShardId:       2,
		HeaderRound:         4,
		IsStartOfEpoch:      true,
	}, proofToAPIProof(headerProof))
}

func TestAddProof(t *testing.T) {
	t.Parallel()

	t.Run("no proof for required block should error", func(t *testing.T) {
		t.Parallel()

		baseAPIBlockProc := createBaseBlockProcessor()
		baseAPIBlockProc.proofsPool = &dataRetrieverTestsCommon.ProofsPoolMock{
			GetProofCalled: func(shardID uint32, headerHash []byte) (data.HeaderProofHandler, error) {
				return nil, errors.New("error")
			},
		}
		baseAPIBlockProc.store = &storageMocks.ChainStorerStub{
			GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
				return nil, errors.New("error")
			},
		}
		baseAPIBlockProc.enableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}

		header := &block.HeaderV2{}

		err := baseAPIBlockProc.addProof([]byte("hash"), header, &api.Block{})
		require.Equal(t, errCannotFindBlockProof, err)
	})

	t.Run("proof for current block returned from pool", func(t *testing.T) {
		t.Parallel()

		baseAPIBlockProc := createBaseBlockProcessor()
		baseAPIBlockProc.proofsPool = &dataRetrieverTestsCommon.ProofsPoolMock{
			GetProofCalled: func(shardID uint32, headerHash []byte) (data.HeaderProofHandler, error) {
				return &block.HeaderProof{
					HeaderHash: []byte("hash2"),
				}, nil
			},
		}
		baseAPIBlockProc.enableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}

		header := &block.HeaderV2{}

		apiBlock := &api.Block{}
		err := baseAPIBlockProc.addProof([]byte("hash"), header, apiBlock)
		require.Nil(t, err)

		require.Equal(t, &api.HeaderProof{
			HeaderHash: hex.EncodeToString([]byte("hash2")),
		}, apiBlock.Proof)
	})

	t.Run("no previous proof only current proof", func(t *testing.T) {
		t.Parallel()

		baseAPIBlockProc := createBaseBlockProcessor()
		baseAPIBlockProc.proofsPool = &dataRetrieverTestsCommon.ProofsPoolMock{
			GetProofCalled: func(shardID uint32, headerHash []byte) (data.HeaderProofHandler, error) {
				return &block.HeaderProof{
					HeaderHash: []byte("hash2"),
				}, nil
			},
		}
		baseAPIBlockProc.enableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}

		baseAPIBlockProc.enableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}
		header := &block.HeaderV2{}

		apiBlock := &api.Block{}
		err := baseAPIBlockProc.addProof([]byte("hash"), header, apiBlock)
		require.Nil(t, err)

		require.Equal(t, &api.HeaderProof{
			HeaderHash: hex.EncodeToString([]byte("hash2")),
		}, apiBlock.Proof)
	})

	t.Run("proof for block returned from storage", func(t *testing.T) {
		t.Parallel()

		baseAPIBlockProc := createBaseBlockProcessor()
		baseAPIBlockProc.proofsPool = &dataRetrieverTestsCommon.ProofsPoolMock{
			GetProofCalled: func(shardID uint32, headerHash []byte) (data.HeaderProofHandler, error) {
				return nil, errors.New("error")
			},
		}
		baseAPIBlockProc.enableEpochsHandler = &enableEpochsHandlerMock.EnableEpochsHandlerStub{
			IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
				return true
			},
		}

		proof := &block.HeaderProof{
			HeaderHash: []byte("hash2"),
		}
		proofBytes, err := baseAPIBlockProc.marshalizer.Marshal(proof)
		require.Nil(t, err)

		baseAPIBlockProc.store = &storageMocks.ChainStorerStub{
			GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
				return &storageMocks.StorerStub{
					GetFromEpochCalled: func(key []byte, epoch uint32) ([]byte, error) {
						return proofBytes, nil
					},
				}, nil
			},
		}

		header := &block.HeaderV2{}

		apiBlock := &api.Block{}
		err = baseAPIBlockProc.addProof([]byte("hash"), header, apiBlock)
		require.Nil(t, err)

		require.Equal(t, &api.HeaderProof{
			HeaderHash: hex.EncodeToString([]byte("hash2")),
		}, apiBlock.Proof)
	})
}

func TestBigInToString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "0", bigIntToStr(big.NewInt(0)))
	require.Equal(t, "0", bigIntToStr(nil))
	require.Equal(t, "15", bigIntToStr(big.NewInt(15)))
	require.Equal(t, "100", bigIntToStr(big.NewInt(100)))
}

func TestBaseBlock_getAndAttachTxsToMb_MiniblockTxBlockgetFromStore(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	baseAPIBlockProc := createBaseBlockProcessor()
	baseAPIBlockProc.store = &storageMocks.ChainStorerStub{
		GetStorerCalled: func(unitType dataRetriever.UnitType) (storage.Storer, error) {
			return nil, expectedErr
		},
	}

	resp, err := baseAPIBlockProc.getFromStorer(dataRetriever.BlockHeaderUnit, nil)
	require.Nil(t, resp)
	require.Equal(t, expectedErr, err)
}
