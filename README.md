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
- pnpm
- Go 1.21+
- Foundry (`forge`)
- Docker

## Installation

```bash
git clone --recurse-submodules https://github.com/symbioticfi/symbiotic-flight-delays.git
cd symbiotic-flight-delays
pnpm install
cd off-chain && go mod tidy
cd ../ui && pnpm install
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

## UI

The UI surfaces a compact list of flights plus basic provider controls (approve collateral, buy coverage, claim after delays, inspect vault balances, deposit/withdraw, claim rewards).

Available at `http://localhost:5173`

## Services

- **anvil** – local execution + settlement chain seeded by the deployer.
- **relay-sidecar-\* containers** – operators/aggregators registered with Symbiotic Relay.
- **flights-api** – exposes airline & flight endpoints at `http://localhost:8085` (seed data plus POST routes for scheduling/delay/depart events).
- **flight-node** – polls the flights API, requests Settlement signatures, and submits `createFlight` / `delayFlight` / `departFlight` transactions.
- **symbiotic-deployer / genesis-generator** – deploy contracts and derive genesis data before the rest of the stack comes online.

Stop and clean everything with:

```bash
docker compose --project-directory temp-network down -v
rm -rf temp-network
```

## Ports

| Service        | Default Port | Notes                                                                         |
| -------------- | ------------ | ----------------------------------------------------------------------------- |
| UI (Vite dev)  | 5173         | `pnpm --filter ui dev`; `pnpm preview` serves on 4173 by default.             |
| Flights API    | 8085         | Mock airline/flight data plus POST endpoints for create/delay/depart actions. |
| Off-chain node | 8080         | Relay-sidecar instances listen on 8080 inside the compose network.            |
| Anvil (L1/L2)  | 8545         | JSON-RPC endpoint used by contracts, oracle node, and the UI.                 |

## API Endpoints

All endpoints below are served by the mock Flights API (`http://localhost:8085` unless overridden):

- `GET /healthz` – readiness probe.
- `GET /airlines` – lists all airlines and their current metadata.
- `POST /airlines` – create an airline (`{ "airlineId": "...", "name": "...", "code": "ALP" }`).
- `GET /airlines/{airlineId}/flights` – list flights for an airline.
- `POST /airlines/{airlineId}/flights` – create/schedule a new flight (`flightId`, `departureTimestamp`).
- `POST /airlines/{airlineId}/flights/{flightId}/delay` – mark a flight as delayed.
- `POST /airlines/{airlineId}/flights/{flightId}/depart` – mark a flight as departed.

## Local Deployments

http://anvil:8545:

- `ValSetDriver`: 0x43C27243F96591892976FFf886511807B65a33d5
- `FlightDelays`: 0xA4b0f5eb09891c1538494c4989Eea0203b1153b1
- `VotingPowerProvider`: 0x369c72C823A4Fc8d2A3A5C3B15082fb34A342878
- `KeyRegistry`: 0xe1557A820E1f50dC962c3392b875Fe0449eb184F
- `Settlement`: 0x882B9439598239d9626164f7578F812Ef324F5Cb
- `Network`: 0xfdc4b2cA12dD7b1463CC01D8022a49BDcf5cFa24
