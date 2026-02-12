// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MWAliasRegistry
 * @dev A decentralized registry for financial aliases in Malawi.
 * Designed to run on a private Hyperledger Besu sidechain managed by Malawian banks.
 */
contract MWAliasRegistry {
    struct Endpoint {
        string provider;
        string destination; // Can be a signed token or encrypted MSISDN
        string endpointType;
    }

    struct AliasRecord {
        string aliasName;
        string identityMask;
        address owner;
        uint256 attestationLevel;
        bool isActive;
        bool isPrivate;
        Endpoint[] endpoints;
    }

    // Mapping from normalized alias string to its record
    mapping(string => AliasRecord) private registry;
    
    // Authorization: Only consortium members can perform critical updates
    mapping(address => bool) public isConsortiumMember;
    address public admin;

    event AliasRegistered(string indexed aliasName, address indexed owner);
    event AliasUpdated(string indexed aliasName);
    event ConsortiumMemberAdded(address indexed member);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can perform this action");
        _;
    }

    modifier onlyConsortium() {
        require(isConsortiumMember[msg.sender], "Only consortium members can update the registry");
        _;
    }

    constructor() {
        admin = msg.sender;
        isConsortiumMember[msg.sender] = true;
    }

    function addConsortiumMember(address member) external onlyAdmin {
        isConsortiumMember[member] = true;
        emit ConsortiumMemberAdded(member);
    }

    /**
     * @dev Registers or updates an alias. 
     * In a real system, banks would verify the identity before calling this.
     */
    function registerAlias(
        string calldata _alias,
        string calldata _identityMask,
        uint256 _attestation,
        bool _isPrivate,
        Endpoint[] calldata _endpoints
    ) external onlyConsortium {
        AliasRecord storage record = registry[_alias];
        
        // If it's a new registration, set the owner and name
        if (bytes(record.aliasName).length == 0) {
            record.aliasName = _alias;
            record.owner = msg.sender; // Placeholder: in reality, the individual or bank
            record.isActive = true;
        }

        record.identityMask = _identityMask;
        record.attestationLevel = _attestation;
        record.isPrivate = _isPrivate;

        // Update endpoints
        delete record.endpoints;
        for (uint i = 0; i < _endpoints.length; i++) {
            record.endpoints.push(_endpoints[i]);
        }

        emit AliasUpdated(_alias);
    }

    /**
     * @dev Resolves an alias. Returns all necessary data for the standard.
     */
    function resolve(string calldata _alias) external view returns (
        string memory identityMask,
        uint256 attestationLevel,
        bool isPrivate,
        Endpoint[] memory endpoints
    ) {
        AliasRecord storage record = registry[_alias];
        require(record.isActive, "Alias not found or inactive");

        return (
            record.identityMask,
            record.attestationLevel,
            record.isPrivate,
            record.endpoints
        );
    }
}
