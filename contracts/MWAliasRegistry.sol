// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

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
        address owner;       // The specific consortium member that registered this alias
        uint256 attestationLevel;
        bool isActive;
        bool isPrivate;
        Endpoint[] endpoints;
    }

    // Mapping from normalized alias string to its record
    mapping(string => AliasRecord) private registry;
    
    // Authorization infrastructure
    mapping(address => bool) public isConsortiumMember;
    address public admin;

    event AliasRegistered(string indexed aliasName, address indexed owner);
    event AliasUpdated(string indexed aliasName, address indexed updater);
    event AliasDeactivated(string indexed aliasName, address indexed deactivator);
    event ConsortiumMemberAdded(address indexed member);
    event ConsortiumMemberRemoved(address indexed member);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can perform this action");
        _;
    }

    modifier onlyConsortium() {
        require(isConsortiumMember[msg.sender], "Only consortium members can access");
        _;
    }

    modifier onlyRecordOwner(string calldata _alias) {
        require(registry[_alias].owner == msg.sender || msg.sender == admin, "Not the owner of this alias record");
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

    function removeConsortiumMember(address member) external onlyAdmin {
        isConsortiumMember[member] = false;
        emit ConsortiumMemberRemoved(member);
    }

    /**
     * @dev Registers or updates an alias with secure ownership protection.
     */
    function registerAlias(
        string calldata _alias,
        string calldata _identityMask,
        uint256 _attestation,
        bool _isPrivate,
        Endpoint[] calldata _endpoints
    ) external onlyConsortium {
        require(_endpoints.length > 0, "At least one endpoint required");
        require(_endpoints.length <= 10, "Exceeds max alternative endpoint paths");

        AliasRecord storage record = registry[_alias];
        
        if (bytes(record.aliasName).length == 0) {
            // Fresh registration mapping
            record.aliasName = _alias;
            record.owner = msg.sender;
            record.isActive = true;
            emit AliasRegistered(_alias, msg.sender);
        } else {
            // Protect existing entries against multi-bank hijacking loops
            require(record.owner == msg.sender || msg.sender == admin, "Alias is owned by another bank");
        }

        record.identityMask = _identityMask;
        record.attestationLevel = _attestation;
        record.isPrivate = _isPrivate;

        // Clean out old elements and map replacements
        delete record.endpoints;
        for (uint256 i = 0; i < _endpoints.length; i++) {
            record.endpoints.push(_endpoints[i]);
        }

        emit AliasUpdated(_alias, msg.sender);
    }

    /**
     * @dev Explicit deactivation pathway for cleaning stale or retired aliases
     */
    function deactivateAlias(string calldata _alias) external onlyConsortium onlyRecordOwner(_alias) {
        AliasRecord storage record = registry[_alias];
        require(record.isActive, "Alias already inactive");
        
        record.isActive = false;
        emit AliasDeactivated(_alias, msg.sender);
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
        require(bytes(record.aliasName).length > 0, "Alias does not exist");
        require(record.isActive, "Alias is currently inactive");

        return (
            record.identityMask,
            record.attestationLevel,
            record.isPrivate,
            record.endpoints
        );
    }
}