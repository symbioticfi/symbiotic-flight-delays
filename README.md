# Symbiotic Flight Delay Insurance

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/symbioticfi/symbiotic-flight-delays)

This repository hosts a full end-to-end example of a flight-delay insurance protocol built on top of Symbiotic Settlement:

- **Smart contracts** (Foundry) manage airlines, vault deployment, premium accounting, slashing, and claims via `FlightDelays.sol` and `VotingPowers.sol`.
- **Mock flights API** (Go) exposes airlines and their flights with status transitions so the rest of the stack can run locally without external data feeds.
- **Off-chain oracle node** (Go) polls the flights API, asks Settlement to sign flight events, and submits `createFlight`, `delayFlight`, and `departFlight` transactions.
- **React/wagmi UI** allows buyers to purchase/claim policies and liquidity providers to deposit into the automatically deployed Symbiotic vaults and claim rewards.

The original super-sum example has been entirely replaced by this flight-delay flow.

## Repository layout

```
├─ src/                       # Solidity contracts
├─ test/                      # Foundry tests
├─ off-chain/cmd/flights-api  # Mock airlines/flights HTTP server
├─ off-chain/cmd/node         # Off-chain oracle node
├─ off-chain/internal         # Shared Go packages
├─ ui/                        # Vite + React application for buyers/providers
└─ network-scripts/           # Helper scripts to spin up local Symbiotic networks (unchanged)
```

## Prerequisites

- Node.js 18+
- npm (used both at repo root and inside `ui/`)
- Go 1.21+
- Foundry (`forge`)
- Docker (only if you plan to run the full Symbiotic devnet via `generate_network.sh`)

## Installation

```bash
npm install
cd off-chain && go mod tidy
cd ../ui && npm install
cd ..
```

## Quick start

1. Generate the docker assets (choose operator counts when prompted):
   ```bash
   ./generate_network.sh
   ```
2. Boot the entire network (anvil, relays, Settlement, flights API, oracle node) in one command:
   ```bash
   docker compose --project-directory temp-network up -d
   ```
3. Inspect or tail logs as needed:
   ```bash
   docker compose --project-directory temp-network ps
   docker compose --project-directory temp-network logs -f flight-node
   ```

The compose stack now includes the mock flights API and the off-chain node, so flights are mirrored on-chain automatically without any extra processes.

## Services

- **anvil** – local execution + settlement chain seeded by the deployer.
- **relay-sidecar-\* containers** – operators/aggregators registered with Symbiotic Relay.
- **flights-api** – exposes airline & flight endpoints at `http://localhost:8085` (seed data plus POST routes for scheduling/delay/depart events).
- **flight-node** – polls the flights API, requests Settlement signatures, and submits `createFlight` / `delayFlight` / `departFlight` transactions.
- **symbiotic-deployer / genesis-generator** – deploy contracts and derive genesis data before the rest of the stack comes online.

Stop and clean everything with:

```bash
docker compose --project-directory temp-network down -v
```

## Optional UI

Run the React/wagmi UI locally (outside docker) for a wallet-connected view:

```bash
cd ui
cp .env.example .env   # adjust RPC URL / addresses if needed
npm run dev
```

The UI surfaces a compact list of flights plus basic provider controls (approve collateral, buy coverage, claim after delays, inspect vault balances, deposit/withdraw, claim rewards).

## Configuration notes

- `ui/.env.example` documents the environment variables (RPC URL, chain metadata, contract address, flights API URL).
- The flights API keeps everything in memory—restart the docker service if you want a clean slate.

## Testing

Use the following commands to verify changes:

```bash
forge test                           # Solidity unit tests for FlightDelays
cd off-chain && go test ./...        # flights API + node packages
cd ui && npm run build               # type-check + build the frontend
```
