# Blockchain Registry Guide

The repository contains an **experimental Besu configuration, genesis and Solidity contract**, not a deployed national registry. The running ALS server does not connect to Besu. No bank participation, consensus deployment, indexing/synchronization, or security audit is evidenced here.

## Technical Stack
- **Engine**: Hyperledger Besu (IBFT 2.0 Consensus).
- **Network**: Example private-network configuration only (no consortium deployment).
- **ChainID**: Check the actual genesis configuration before running any node.

## Smart Contract: MWAliasRegistry
The contract is a design experiment; no synchronization from contract events into ALS is implemented. The JSON-backed server is currently the only running registry implementation. Before proposing on-chain identity records, review privacy, key custody, access control, contract upgrades, fees and data retention with local stakeholders. The sample Besu RPC binds to loopback by default; never expose privileged RPC APIs to the Internet.
