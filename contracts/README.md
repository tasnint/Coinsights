# Coinsights Smart Contracts

## ResolutionAttestation.sol

A minimal smart contract for recording tamper-proof attestations of resolved cryptocurrency exchange issues.

### Design Philosophy

This contract follows the principle: **"Blockchain as a verification primitive, not a database."**

- ✅ Accepts evidence hashes (off-chain privacy)
- ✅ Records timestamps (immutable proof)
- ✅ Emits events (external indexing)
- ✅ Chain-of-custody (linked attestations)
- ❌ No raw data storage
- ❌ No AI logic
- ❌ No business rules

### Deployment

**Recommended Networks (Coinbase-aligned):**

1. **Base Sepolia** (testnet) - `chainId: 84532`
   - RPC: `https://sepolia.base.org`
   - Explorer: `https://sepolia.basescan.org`

2. **Base Mainnet** - `chainId: 8453`
   - RPC: `https://mainnet.base.org`
   - Explorer: `https://basescan.org`

**Alternative Networks:**

- Ethereum Sepolia: `chainId: 11155111`
- Polygon Mumbai: `chainId: 80001`

### Deployment with Foundry

```bash
# Install Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Compile
forge build

# Deploy to Base Sepolia
forge create --rpc-url https://sepolia.base.org \
  --private-key $PRIVATE_KEY \
  contracts/ResolutionAttestation.sol:ResolutionAttestation

# Verify on BaseScan
forge verify-contract <CONTRACT_ADDRESS> \
  contracts/ResolutionAttestation.sol:ResolutionAttestation \
  --chain-id 84532 \
  --etherscan-api-key $BASESCAN_API_KEY
```

### Contract Interface

```solidity
// Record a single resolution
function recordResolution(
    string calldata exchange,      // "coinbase"
    string calldata issueCategory, // "withdrawal_delays"
    bytes32 evidenceHash           // keccak256(jsonEvidence)
) external returns (uint256 attestationId);

// Batch record (gas optimized)
function recordBatch(
    bytes32 merkleRoot,
    string[] calldata exchanges,
    string[] calldata categories,
    bytes32[] calldata evidenceHashes
) external returns (uint256 startId, uint256 endId);

// Verify a hash exists
function verifyHash(bytes32 evidenceHash) external view returns (
    bool exists,
    uint256 attestationId
);
```

### Events

```solidity
event ResolutionRecorded(
    uint256 indexed attestationId,
    string indexed exchange,
    string issueCategory,
    bytes32 evidenceHash,
    bytes32 previousHash,
    uint256 timestamp,
    address attestor
);
```

### Gas Estimates

| Function | Estimated Gas |
|----------|---------------|
| `recordResolution` | ~80,000 |
| `recordBatch` (10 items) | ~400,000 |
| `verifyHash` | ~30,000 (view) |

### Security Considerations

1. **No sensitive data on-chain** - Only hashes stored
2. **Chain-of-custody** - Each attestation links to previous
3. **Event-based indexing** - Full history recoverable
4. **Owner controls** - Minimal, for future governance only
