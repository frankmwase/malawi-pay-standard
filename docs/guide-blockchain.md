# Blockchain Registry Guide

For national-scale interoperability, the Malawi Pay Standard utilizes a **Decentralized Registry** built on Hyperledger Besu.

## Technical Stack
- **Engine**: Hyperledger Besu (IBFT 2.0 Consensus).
- **Network**: Private Malawian Sidechain (Consortium managed).
- **ChainID**: 650

## Smart Contract: MWAliasRegistry
Aliases are registered on-chain to ensure every bank in Malawi sees the same "Source of Truth".

### Registration Flow
1. Bank verifies customer identity (KYC).
2. Bank node calls `registerAlias()` on the smart contract.
3. All other banks' ALS nodes automatically sync the new alias.

## Advantages
- **No Single Owner**: No single entity can "shut down" or "censor" an alias.
- **Auditability**: Every alias change is recorded permanently on the ledger.
- **Transparency**: Real-time resolution regardless of which bank the sender uses.
