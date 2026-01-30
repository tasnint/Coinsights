// records proof that coinbase has resolved a specific issue on the blockchain.
package services

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tasnint/coinsights/internal/models"
	"golang.org/x/crypto/sha3"
)

// ============================================
// CONTRACT ABI (Minimal - only what we need)
// ============================================

const ResolutionAttestationABI = `[
	{
		"inputs": [
			{"internalType": "string", "name": "exchange", "type": "string"},
			{"internalType": "string", "name": "issueCategory", "type": "string"},
			{"internalType": "bytes32", "name": "evidenceHash", "type": "bytes32"}
		],
		"name": "recordResolution",
		"outputs": [{"internalType": "uint256", "name": "attestationId", "type": "uint256"}],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [{"internalType": "bytes32", "name": "evidenceHash", "type": "bytes32"}],
		"name": "verifyHash",
		"outputs": [
			{"internalType": "bool", "name": "exists", "type": "bool"},
			{"internalType": "uint256", "name": "attestationId", "type": "uint256"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [{"internalType": "uint256", "name": "attestationId", "type": "uint256"}],
		"name": "getAttestation",
		"outputs": [
			{"internalType": "bytes32", "name": "evidenceHash", "type": "bytes32"},
			{"internalType": "bytes32", "name": "previousHash", "type": "bytes32"},
			{"internalType": "uint256", "name": "timestamp", "type": "uint256"},
			{"internalType": "uint256", "name": "blockNumber", "type": "uint256"},
			{"internalType": "string", "name": "exchange", "type": "string"},
			{"internalType": "string", "name": "issueCategory", "type": "string"},
			{"internalType": "address", "name": "attestor", "type": "address"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "attestationCount",
		"outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "uint256", "name": "attestationId", "type": "uint256"},
			{"indexed": true, "internalType": "string", "name": "exchange", "type": "string"},
			{"indexed": false, "internalType": "string", "name": "issueCategory", "type": "string"},
			{"indexed": false, "internalType": "bytes32", "name": "evidenceHash", "type": "bytes32"},
			{"indexed": false, "internalType": "bytes32", "name": "previousHash", "type": "bytes32"},
			{"indexed": false, "internalType": "uint256", "name": "timestamp", "type": "uint256"},
			{"indexed": false, "internalType": "address", "name": "attestor", "type": "address"}
		],
		"name": "ResolutionRecorded",
		"type": "event"
	}
]`

// ============================================
// BLOCKCHAIN SERVICE
// ============================================

// BlockchainService handles all blockchain interactions
type BlockchainService struct {
	client          *ethclient.Client
	chainConfig     models.ChainConfig
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      *ecdsa.PrivateKey
	publicAddress   common.Address
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService() (*BlockchainService, error) {
	// Get chain configuration
	chainName := os.Getenv("BLOCKCHAIN_NETWORK")
	if chainName == "" {
		chainName = "base_sepolia" // Default to Base testnet (Coinbase-aligned)
	}

	chains := models.SupportedChains()
	chainConfig, ok := chains[chainName]
	if !ok {
		return nil, fmt.Errorf("unsupported blockchain network: %s", chainName)
	}

	// Override RPC URL if provided
	if rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL"); rpcURL != "" {
		chainConfig.RPCURL = rpcURL
	}

	// Get contract address
	contractAddr := os.Getenv("ATTESTATION_CONTRACT_ADDRESS")
	if contractAddr == "" {
		return nil, fmt.Errorf("ATTESTATION_CONTRACT_ADDRESS not set")
	}
	chainConfig.ContractAddress = contractAddr

	// Connect to blockchain
	client, err := ethclient.Dial(chainConfig.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	// Parse contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(ResolutionAttestationABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	// Load private key for signing transactions
	privateKeyHex := os.Getenv("BLOCKCHAIN_PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("BLOCKCHAIN_PRIVATE_KEY not set")
	}

	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &BlockchainService{
		client:          client,
		chainConfig:     chainConfig,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     parsedABI,
		privateKey:      privateKey,
		publicAddress:   publicAddress,
	}, nil
}

// Close closes the blockchain connection
func (bs *BlockchainService) Close() {
	if bs.client != nil {
		bs.client.Close()
	}
}

// GetChainInfo returns current chain configuration
func (bs *BlockchainService) GetChainInfo() models.ChainConfig {
	return bs.chainConfig
}

// GetWalletAddress returns the wallet address used for attestations
func (bs *BlockchainService) GetWalletAddress() string {
	return bs.publicAddress.Hex()
}

// ============================================
// HASHING FUNCTIONS
// ============================================

// HashEvidence creates a Keccak256 hash of the resolution evidence
// This is the hash that gets stored on-chain
func (bs *BlockchainService) HashEvidence(evidence *models.ResolutionEvidence) (string, error) {
	// Serialize evidence to canonical JSON
	jsonBytes, err := json.Marshal(evidence)
	if err != nil {
		return "", fmt.Errorf("failed to serialize evidence: %w", err)
	}

	// Compute Keccak256 hash (same as Solidity's keccak256)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(jsonBytes)
	hashBytes := hash.Sum(nil)

	return "0x" + hex.EncodeToString(hashBytes), nil
}

// HashEvidenceBytes returns the raw 32-byte hash
func (bs *BlockchainService) HashEvidenceBytes(evidence *models.ResolutionEvidence) ([32]byte, error) {
	var hashArray [32]byte

	jsonBytes, err := json.Marshal(evidence)
	if err != nil {
		return hashArray, fmt.Errorf("failed to serialize evidence: %w", err)
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(jsonBytes)
	copy(hashArray[:], hash.Sum(nil))

	return hashArray, nil
}

// ============================================
// ON-CHAIN OPERATIONS
// ============================================

// RecordAttestation records a resolution on the blockchain
func (bs *BlockchainService) RecordAttestation(
	ctx context.Context,
	resolution *models.Resolution,
) (*models.Attestation, error) {
	fmt.Printf("⛓️  Recording attestation for %s - %s\n", resolution.Exchange, resolution.IssueCategory)

	// Hash the evidence
	evidenceHash, err := bs.HashEvidenceBytes(&resolution.Evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to hash evidence: %w", err)
	}
	fmt.Printf("   Evidence hash: 0x%x\n", evidenceHash)

	// Get nonce
	nonce, err := bs.client.PendingNonceAt(ctx, bs.publicAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := bs.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Build transaction data
	txData, err := bs.contractABI.Pack(
		"recordResolution",
		resolution.Exchange,
		resolution.IssueCategory,
		evidenceHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transaction data: %w", err)
	}

	// Estimate gas
	gasLimit := uint64(150000) // Conservative estimate

	// Create transaction
	tx := types.NewTransaction(
		nonce,
		bs.contractAddress,
		big.NewInt(0), // No ETH value
		gasLimit,
		gasPrice,
		txData,
	)

	// Sign transaction
	chainID := big.NewInt(bs.chainConfig.ChainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), bs.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = bs.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("   Transaction sent: %s\n", txHash)

	// Wait for receipt
	receipt, err := bs.waitForReceipt(ctx, signedTx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction reverted")
	}

	// Get block timestamp
	block, err := bs.client.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	// Build attestation result
	attestation := &models.Attestation{
		TransactionHash: txHash,
		BlockNumber:     receipt.BlockNumber.Uint64(),
		BlockTimestamp:  time.Unix(int64(block.Time()), 0),
		ChainID:         bs.chainConfig.ChainID,
		ContractAddress: bs.contractAddress.Hex(),
		EvidenceHash:    "0x" + hex.EncodeToString(evidenceHash[:]),
		Attestor:        bs.publicAddress.Hex(),
		ExplorerURL:     fmt.Sprintf("%s/tx/%s", bs.chainConfig.ExplorerURL, txHash),
		Verified:        true,
	}

	// Try to get attestation ID from logs
	attestation.ID = bs.parseAttestationID(receipt.Logs)

	fmt.Printf("   ✅ Attestation recorded! Block: %d\n", attestation.BlockNumber)
	fmt.Printf("   🔗 Explorer: %s\n", attestation.ExplorerURL)

	return attestation, nil
}

// VerifyAttestation verifies an attestation exists on-chain
func (bs *BlockchainService) VerifyAttestation(
	ctx context.Context,
	evidenceHash string,
) (*models.VerificationResponse, error) {
	// Convert hex string to bytes32
	hashBytes, err := hex.DecodeString(strings.TrimPrefix(evidenceHash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid hash format: %w", err)
	}

	var hash32 [32]byte
	copy(hash32[:], hashBytes)

	// Call verifyHash on contract
	callData, err := bs.contractABI.Pack("verifyHash", hash32)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	result, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	// Unpack result
	var exists bool
	var attestationID *big.Int

	outputs, err := bs.contractABI.Unpack("verifyHash", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	exists = outputs[0].(bool)
	attestationID = outputs[1].(*big.Int)

	response := &models.VerificationResponse{
		OnChain:   exists,
		Verified:  exists,
		HashMatch: exists,
	}

	if exists {
		response.Message = fmt.Sprintf("Hash verified on-chain. Attestation ID: %d", attestationID.Uint64())

		// Get full attestation details
		attestation, err := bs.GetAttestationByID(ctx, attestationID.Uint64())
		if err == nil {
			response.Attestation = attestation
			response.TimestampValid = true
		}
	} else {
		response.Message = "Hash not found on-chain"
	}

	return response, nil
}

// GetAttestationByID retrieves an attestation by its on-chain ID
func (bs *BlockchainService) GetAttestationByID(
	ctx context.Context,
	attestationID uint64,
) (*models.Attestation, error) {
	callData, err := bs.contractABI.Pack("getAttestation", big.NewInt(int64(attestationID)))
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	result, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %w", err)
	}

	outputs, err := bs.contractABI.Unpack("getAttestation", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	evidenceHash := outputs[0].([32]byte)
	previousHash := outputs[1].([32]byte)
	timestamp := outputs[2].(*big.Int)
	blockNumber := outputs[3].(*big.Int)
	// exchange := outputs[4].(string) // We could use these if needed
	// issueCategory := outputs[5].(string)
	attestor := outputs[6].(common.Address)

	return &models.Attestation{
		ID:              attestationID,
		BlockNumber:     blockNumber.Uint64(),
		BlockTimestamp:  time.Unix(timestamp.Int64(), 0),
		ChainID:         bs.chainConfig.ChainID,
		ContractAddress: bs.contractAddress.Hex(),
		EvidenceHash:    "0x" + hex.EncodeToString(evidenceHash[:]),
		PreviousHash:    "0x" + hex.EncodeToString(previousHash[:]),
		Attestor:        attestor.Hex(),
		ExplorerURL:     fmt.Sprintf("%s/address/%s", bs.chainConfig.ExplorerURL, bs.contractAddress.Hex()),
		Verified:        true,
	}, nil
}

// GetAttestationCount returns the total number of attestations
func (bs *BlockchainService) GetAttestationCount(ctx context.Context) (uint64, error) {
	callData, err := bs.contractABI.Pack("attestationCount")
	if err != nil {
		return 0, fmt.Errorf("failed to pack call data: %w", err)
	}

	result, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("contract call failed: %w", err)
	}

	outputs, err := bs.contractABI.Unpack("attestationCount", result)
	if err != nil {
		return 0, fmt.Errorf("failed to unpack result: %w", err)
	}

	count := outputs[0].(*big.Int)
	return count.Uint64(), nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// waitForReceipt waits for a transaction receipt with timeout
func (bs *BlockchainService) waitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction receipt")
		case <-ticker.C:
			receipt, err := bs.client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			// Continue waiting if receipt not available yet
		}
	}
}

// parseAttestationID extracts the attestation ID from transaction logs
func (bs *BlockchainService) parseAttestationID(logs []*types.Log) uint64 {
	eventSig := bs.contractABI.Events["ResolutionRecorded"].ID

	for _, log := range logs {
		if len(log.Topics) > 0 && log.Topics[0] == eventSig {
			// The attestation ID is the first indexed parameter
			if len(log.Topics) > 1 {
				return new(big.Int).SetBytes(log.Topics[1].Bytes()).Uint64()
			}
		}
	}
	return 0
}
