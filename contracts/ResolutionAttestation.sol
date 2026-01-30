// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ResolutionAttestation
 * @dev Minimal smart contract for recording tamper-proof attestations of resolved issues.
 * Designed for Coinbase's Base network (or Ethereum testnets).
 * 
 * This contract does NOT:
 * - Store raw data (privacy-preserving)
 * - Perform AI analysis
 * - Make decisions
 * 
 * This contract DOES:
 * - Accept evidence hashes
 * - Record timestamps
 * - Emit verifiable events
 * - Create an append-only audit trail
 */
contract ResolutionAttestation {
    
    // ============================================
    // STRUCTS
    // ============================================
    
    struct Attestation {
        bytes32 evidenceHash;      // Keccak256 hash of resolution evidence
        bytes32 previousHash;      // Hash of previous attestation (chain-of-custody)
        uint256 timestamp;         // Block timestamp when recorded
        uint256 blockNumber;       // Block number for verification
        string exchange;           // Exchange name (e.g., "coinbase", "kraken")
        string issueCategory;      // Category (e.g., "withdrawal_delays", "support_issues")
        address attestor;          // Address that submitted the attestation
    }
    
    // ============================================
    // STATE VARIABLES
    // ============================================
    
    // Mapping from attestation ID to Attestation data
    mapping(uint256 => Attestation) public attestations;
    
    // Total number of attestations
    uint256 public attestationCount;
    
    // Latest attestation hash per exchange+issue (for chain-of-custody)
    mapping(bytes32 => bytes32) public latestHashByIssue;
    
    // Owner for potential future governance
    address public owner;
    
    // ============================================
    // EVENTS
    // ============================================
    
    /**
     * @dev Emitted when a new resolution is attested on-chain
     * This is the primary event that external systems will monitor
     */
    event ResolutionRecorded(
        uint256 indexed attestationId,
        string indexed exchange,
        string issueCategory,
        bytes32 evidenceHash,
        bytes32 previousHash,
        uint256 timestamp,
        address attestor
    );
    
    /**
     * @dev Emitted when batch attestations are recorded (gas optimization)
     */
    event BatchRecorded(
        uint256 startId,
        uint256 endId,
        bytes32 merkleRoot,
        uint256 timestamp
    );
    
    // ============================================
    // MODIFIERS
    // ============================================
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not authorized");
        _;
    }
    
    // ============================================
    // CONSTRUCTOR
    // ============================================
    
    constructor() {
        owner = msg.sender;
        attestationCount = 0;
    }
    
    // ============================================
    // CORE FUNCTIONS
    // ============================================
    
    /**
     * @dev Record a single resolution attestation
     * @param exchange Name of the exchange (e.g., "coinbase")
     * @param issueCategory Category of the issue (e.g., "withdrawal_delays")
     * @param evidenceHash Keccak256 hash of the resolution evidence JSON
     * @return attestationId The ID of the newly created attestation
     */
    function recordResolution(
        string calldata exchange,
        string calldata issueCategory,
        bytes32 evidenceHash
    ) external returns (uint256 attestationId) {
        // Get the issue key for chain-of-custody tracking
        bytes32 issueKey = keccak256(abi.encodePacked(exchange, issueCategory));
        bytes32 previousHash = latestHashByIssue[issueKey];
        
        // Create new attestation
        attestationId = attestationCount;
        attestations[attestationId] = Attestation({
            evidenceHash: evidenceHash,
            previousHash: previousHash,
            timestamp: block.timestamp,
            blockNumber: block.number,
            exchange: exchange,
            issueCategory: issueCategory,
            attestor: msg.sender
        });
        
        // Update chain-of-custody
        latestHashByIssue[issueKey] = evidenceHash;
        attestationCount++;
        
        // Emit event for off-chain indexing
        emit ResolutionRecorded(
            attestationId,
            exchange,
            issueCategory,
            evidenceHash,
            previousHash,
            block.timestamp,
            msg.sender
        );
        
        return attestationId;
    }
    
    /**
     * @dev Record multiple attestations in a single transaction (gas efficient)
     * Uses Merkle root for batch verification
     * @param merkleRoot Root hash of all evidence hashes in the batch
     * @param exchanges Array of exchange names
     * @param categories Array of issue categories
     * @param evidenceHashes Array of evidence hashes
     */
    function recordBatch(
        bytes32 merkleRoot,
        string[] calldata exchanges,
        string[] calldata categories,
        bytes32[] calldata evidenceHashes
    ) external returns (uint256 startId, uint256 endId) {
        require(
            exchanges.length == categories.length && 
            categories.length == evidenceHashes.length,
            "Array length mismatch"
        );
        require(exchanges.length > 0, "Empty batch");
        
        startId = attestationCount;
        
        for (uint256 i = 0; i < exchanges.length; i++) {
            bytes32 issueKey = keccak256(abi.encodePacked(exchanges[i], categories[i]));
            bytes32 previousHash = latestHashByIssue[issueKey];
            
            attestations[attestationCount] = Attestation({
                evidenceHash: evidenceHashes[i],
                previousHash: previousHash,
                timestamp: block.timestamp,
                blockNumber: block.number,
                exchange: exchanges[i],
                issueCategory: categories[i],
                attestor: msg.sender
            });
            
            latestHashByIssue[issueKey] = evidenceHashes[i];
            attestationCount++;
        }
        
        endId = attestationCount - 1;
        
        emit BatchRecorded(startId, endId, merkleRoot, block.timestamp);
        
        return (startId, endId);
    }
    
    // ============================================
    // VIEW FUNCTIONS
    // ============================================
    
    /**
     * @dev Get attestation details by ID
     */
    function getAttestation(uint256 attestationId) external view returns (
        bytes32 evidenceHash,
        bytes32 previousHash,
        uint256 timestamp,
        uint256 blockNumber,
        string memory exchange,
        string memory issueCategory,
        address attestor
    ) {
        require(attestationId < attestationCount, "Attestation does not exist");
        Attestation storage a = attestations[attestationId];
        return (
            a.evidenceHash,
            a.previousHash,
            a.timestamp,
            a.blockNumber,
            a.exchange,
            a.issueCategory,
            a.attestor
        );
    }
    
    /**
     * @dev Verify if an evidence hash exists on-chain
     * @param evidenceHash The hash to verify
     * @return exists Whether the hash has been recorded
     * @return attestationId The ID of the attestation (if exists)
     */
    function verifyHash(bytes32 evidenceHash) external view returns (
        bool exists,
        uint256 attestationId
    ) {
        for (uint256 i = 0; i < attestationCount; i++) {
            if (attestations[i].evidenceHash == evidenceHash) {
                return (true, i);
            }
        }
        return (false, 0);
    }
    
    /**
     * @dev Get the latest attestation hash for a specific exchange+issue
     */
    function getLatestHash(
        string calldata exchange,
        string calldata issueCategory
    ) external view returns (bytes32) {
        bytes32 issueKey = keccak256(abi.encodePacked(exchange, issueCategory));
        return latestHashByIssue[issueKey];
    }
    
    /**
     * @dev Get all attestations for pagination
     */
    function getAttestationRange(
        uint256 startId,
        uint256 count
    ) external view returns (Attestation[] memory) {
        require(startId < attestationCount, "Start ID out of bounds");
        
        uint256 endId = startId + count;
        if (endId > attestationCount) {
            endId = attestationCount;
        }
        
        Attestation[] memory result = new Attestation[](endId - startId);
        for (uint256 i = startId; i < endId; i++) {
            result[i - startId] = attestations[i];
        }
        
        return result;
    }
    
    // ============================================
    // ADMIN FUNCTIONS
    // ============================================
    
    /**
     * @dev Transfer ownership
     */
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid address");
        owner = newOwner;
    }
}
