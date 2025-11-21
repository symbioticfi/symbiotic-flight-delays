import { useMemo, useState } from "react";
import { useAccount, useConnect, useDisconnect, usePublicClient, useWriteContract } from "wagmi";
import { useQuery } from "@tanstack/react-query";
import { Address, Hex, encodeAbiParameters, maxUint256, parseUnits, zeroAddress } from "viem";
import type { Abi } from "viem";

import flightDelaysAbi from "./abi/FlightDelays.json";
import erc20Abi from "./abi/ERC20.json";
import vaultAbi from "./abi/Vault.json";
import rewardsAbi from "./abi/Rewards.json";
import { flightsApiUrl, flightDelaysAddress } from "./config";
import { fetchAirlinesWithFlights } from "./api/flights";
import type { AirlineWithFlights } from "./types";
import { flightKey, hashIdentifier } from "./utils/hash";
import { formatAmount, formatTimestamp } from "./utils/format";

const flightDelaysABI = flightDelaysAbi as Abi;
const erc20ABI = erc20Abi as Abi;
const vaultABI = vaultAbi as Abi;
const rewardsABI = rewardsAbi as Abi;

interface ProtocolInfo {
  policyPremium: bigint;
  policyPayout: bigint;
  policyWindow: bigint;
  delayWindow: bigint;
  collateral: Address;
  collateralSymbol: string;
  collateralDecimals: number;
  network: Address;
}

interface ChainFlightData {
  [key: string]: {
    timestamp: bigint;
    status: number;
    policiesSold: bigint;
  };
}

interface PolicyMap {
  [key: string]: number;
}

const defaultRewardsToClaim = 5n;

export default function App() {
  const { address, isConnected } = useAccount();
  const publicClient = usePublicClient();
  const { writeContractAsync } = useWriteContract();
  const [feedback, setFeedback] = useState<string | null>(null);

  const flightsQuery = useQuery({
    queryKey: ["flights", flightsApiUrl],
    queryFn: () => fetchAirlinesWithFlights(flightsApiUrl),
    refetchInterval: 10_000
  });

  const flights = flightsQuery.data ?? [];

  const flattenedFlights = useMemo(
    () =>
      flights.flatMap((airline) =>
        airline.flights.map((flight) => ({
          airline,
          flight
        }))
      ),
    [flights]
  );

  const protocolQuery = useQuery<ProtocolInfo>({
    queryKey: ["protocol-info", flightDelaysAddress],
    enabled: !!publicClient,
    queryFn: async () => {
      if (!publicClient) {
        throw new Error("Missing public client");
      }
      const [policyPremium, policyPayout, policyWindow, delayWindow, collateral, network] =
        await publicClient.multicall({
          contracts: [
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyPremium" },
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyPayout" },
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyWindow" },
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "delayWindow" },
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "collateral" },
            { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "NETWORK" }
          ]
        });
      const collateralAddress = (collateral.result ?? zeroAddress) as Address;
      const [symbol, decimals] = await publicClient.multicall({
        contracts: [
          { address: collateralAddress, abi: erc20ABI, functionName: "symbol" },
          { address: collateralAddress, abi: erc20ABI, functionName: "decimals" }
        ]
      });
      return {
        policyPremium: (policyPremium.result ?? 0n) as bigint,
        policyPayout: (policyPayout.result ?? 0n) as bigint,
        policyWindow: (policyWindow.result ?? 0n) as bigint,
        delayWindow: (delayWindow.result ?? 0n) as bigint,
        collateral: collateralAddress,
        collateralSymbol: (symbol.result ?? "TOKEN") as string,
        collateralDecimals: Number(decimals.result ?? 18),
        network: (network.result ?? zeroAddress) as Address
      };
    }
  });

  const chainFlightsQuery = useQuery<{ flights: ChainFlightData; policies: PolicyMap }>({
    queryKey: [
      "chain-flights",
      flattenedFlights.map(({ airline, flight }) => `${airline.airlineId}-${flight.flightId}`),
      address
    ],
    enabled: !!publicClient && flattenedFlights.length > 0,
    refetchInterval: 6_000,
    queryFn: async () => {
      if (!publicClient) {
        throw new Error("Missing client");
      }
      const flightCalls = flattenedFlights.map(({ airline, flight }) => ({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "flights",
        args: [hashIdentifier(airline.airlineId), hashIdentifier(flight.flightId)] as const
      }));
      const responses = await publicClient.multicall({ contracts: flightCalls, allowFailure: true });
      const map: ChainFlightData = {};
      responses.forEach((resp, idx) => {
        const key = flightKey(flattenedFlights[idx].airline.airlineId, flattenedFlights[idx].flight.flightId);
        if (resp.status === "success") {
          const [timestamp, status, policiesSold] = resp.result as unknown as [bigint, number, bigint];
          map[key] = { timestamp, status: Number(status), policiesSold };
        }
      });

      const policiesMap: PolicyMap = {};
      if (address) {
        const policyCalls = flattenedFlights.map(({ airline, flight }) => ({
          address: flightDelaysAddress,
          abi: flightDelaysABI,
          functionName: "policies",
          args: [hashIdentifier(airline.airlineId), hashIdentifier(flight.flightId), address] as const
        }));
        const policyResponses = await publicClient.multicall({ contracts: policyCalls, allowFailure: true });
        policyResponses.forEach((resp, idx) => {
          const key = flightKey(flattenedFlights[idx].airline.airlineId, flattenedFlights[idx].flight.flightId);
          if (resp.status === "success") {
            policiesMap[key] = Number(resp.result ?? 0);
          }
        });
      }

      return { flights: map, policies: policiesMap };
    }
  });

  const airlinesOnChainQuery = useQuery({
    queryKey: ["airlines-chain", flights.map((a) => a.airlineId)],
    enabled: !!publicClient && flights.length > 0,
    refetchInterval: 12_000,
    queryFn: async () => {
      if (!publicClient) throw new Error("Missing client");
      const calls = flights.map((airline) => ({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "airlines",
        args: [hashIdentifier(airline.airlineId)] as const
      }));
      const responses = await publicClient.multicall({ contracts: calls, allowFailure: true });
      const map = new Map<string, { vault: Address; rewards: Address }>();
      responses.forEach((resp, idx) => {
        if (resp.status === "success") {
          const [vault, rewards] = resp.result as unknown as [Address, Address, bigint, Hex];
          map.set(flights[idx].airlineId, { vault, rewards });
        }
      });
      return map;
    }
  });

  const allowanceQuery = useQuery({
    queryKey: ["allowance", address, protocolQuery.data?.collateral],
    enabled: !!publicClient && !!address && !!protocolQuery.data,
    refetchInterval: 8_000,
    queryFn: async () => {
      if (!publicClient || !address || !protocolQuery.data) throw new Error("Missing data");
      return (await publicClient.readContract({
        address: protocolQuery.data.collateral,
        abi: erc20ABI,
        functionName: "allowance",
        args: [address, flightDelaysAddress]
      })) as bigint;
    }
  });

  const vaultAllowancesQuery = useQuery({
    queryKey: ["vault-allowances", address, airlinesOnChainQuery.data?.size ?? 0],
    enabled:
      !!publicClient && !!address && !!protocolQuery.data && !!airlinesOnChainQuery.data?.size,
    refetchInterval: 12_000,
    queryFn: async () => {
      if (!publicClient || !address || !protocolQuery.data || !airlinesOnChainQuery.data) {
        throw new Error("Missing data");
      }
      const map = new Map<string, bigint>();
      for (const [airlineId, info] of airlinesOnChainQuery.data.entries()) {
        try {
          const allowance = (await publicClient.readContract({
            address: protocolQuery.data.collateral,
            abi: erc20ABI,
            functionName: "allowance",
            args: [address, info.vault]
          })) as bigint;
          map.set(airlineId, allowance);
        } catch (err) {
          console.warn("vault allowance read failed", err);
        }
      }
      return map;
    }
  });

  const vaultBalancesQuery = useQuery({
    queryKey: ["vault-balances", address, airlinesOnChainQuery.data?.size ?? 0],
    enabled: !!publicClient && !!address && !!airlinesOnChainQuery.data?.size,
    refetchInterval: 10_000,
    queryFn: async () => {
      if (!publicClient || !address || !airlinesOnChainQuery.data) throw new Error("Missing data");
      const balances = new Map<string, bigint>();
      for (const [airlineId, info] of airlinesOnChainQuery.data.entries()) {
        try {
          const bal = (await publicClient.readContract({
            address: info.vault,
            abi: vaultABI,
            functionName: "activeBalanceOf",
            args: [address]
          })) as bigint;
          balances.set(airlineId, bal);
        } catch (err) {
          console.warn("Failed to read vault balance", err);
        }
      }
      return balances;
    }
  });

  const rewardEstimatesQuery = useQuery({
    queryKey: ["rewards", address, airlinesOnChainQuery.data?.size ?? 0],
    enabled: !!publicClient && !!address && !!protocolQuery.data && !!airlinesOnChainQuery.data?.size,
    refetchInterval: 15_000,
    queryFn: async () => {
      if (!publicClient || !address || !protocolQuery.data || !airlinesOnChainQuery.data) {
        throw new Error("Missing data");
      }
      const map = new Map<string, bigint>();
      for (const [airlineId, info] of airlinesOnChainQuery.data.entries()) {
        try {
          const encoded = encodeAbiParameters(
            [
              { name: "network", type: "address" },
              { name: "maxRewards", type: "uint256" },
              { name: "hints", type: "bytes[]" }
            ],
            [protocolQuery.data.network, defaultRewardsToClaim, []]
          );
          const claimable = (await publicClient.readContract({
            address: info.rewards,
            abi: rewardsABI,
            functionName: "claimable",
            args: [protocolQuery.data.collateral, address, encoded]
          })) as bigint;
          map.set(airlineId, claimable);
        } catch (err) {
          console.warn("claimable read failed", err);
        }
      }
      return map;
    }
  });

  const [depositInputs, setDepositInputs] = useState<Record<string, string>>({});
  const [withdrawInputs, setWithdrawInputs] = useState<Record<string, string>>({});
  const [maxRewardsInputs, setMaxRewardsInputs] = useState<Record<string, string>>({});

  const chainFlights = chainFlightsQuery.data?.flights ?? {};
  const policies = chainFlightsQuery.data?.policies ?? {};

  const handleApprove = async () => {
    if (!protocolQuery.data || !isConnected || !address) return;
    try {
      const hash = await writeContractAsync({
        address: protocolQuery.data.collateral,
        abi: erc20ABI,
        functionName: "approve",
        args: [flightDelaysAddress, protocolQuery.data.policyPremium]
      });
      setFeedback(`Approve tx submitted: ${hash}`);
    } catch (err) {
      console.error(err);
      setFeedback("Approve failed. Check console for details.");
    }
  };

  const handleBuy = async (airlineId: string, flightId: string) => {
    if (!address) return;
    try {
      const hash = await writeContractAsync({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "buyInsurance",
        args: [hashIdentifier(airlineId), hashIdentifier(flightId)]
      });
      setFeedback(`Buy tx submitted: ${hash}`);
    } catch (err) {
      console.error(err);
      setFeedback("Buy failed. See console for details.");
    }
  };

  const handleClaim = async (airlineId: string, flightId: string) => {
    if (!address) return;
    try {
      const hash = await writeContractAsync({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "claimInsurance",
        args: [hashIdentifier(airlineId), hashIdentifier(flightId)]
      });
      setFeedback(`Claim tx submitted: ${hash}`);
    } catch (err) {
      console.error(err);
      setFeedback("Claim failed. See console for details.");
    }
  };

  const handleDeposit = async (airlineId: string) => {
    if (!protocolQuery.data || !airlinesOnChainQuery.data || !address) return;
    const amount = depositInputs[airlineId];
    if (!amount) return;
    try {
      const parsed = parseUnits(amount, protocolQuery.data.collateralDecimals);
      const vault = airlinesOnChainQuery.data.get(airlineId)?.vault;
      if (!vault) return;
      const hash = await writeContractAsync({
        address: vault,
        abi: vaultABI,
        functionName: "deposit",
        args: [address as Address, parsed]
      });
      setFeedback(`Deposit tx submitted: ${hash}`);
      setDepositInputs((prev) => ({ ...prev, [airlineId]: "" }));
    } catch (err) {
      console.error(err);
      setFeedback("Deposit failed.");
    }
  };

  const handleWithdraw = async (airlineId: string) => {
    if (!protocolQuery.data || !airlinesOnChainQuery.data || !address) return;
    const amount = withdrawInputs[airlineId];
    if (!amount) return;
    try {
      const parsed = parseUnits(amount, protocolQuery.data.collateralDecimals);
      const vault = airlinesOnChainQuery.data.get(airlineId)?.vault;
      if (!vault) return;
      const hash = await writeContractAsync({
        address: vault,
        abi: vaultABI,
        functionName: "withdraw",
        args: [address as Address, parsed]
      });
      setFeedback(`Withdraw tx submitted: ${hash}`);
      setWithdrawInputs((prev) => ({ ...prev, [airlineId]: "" }));
    } catch (err) {
      console.error(err);
      setFeedback("Withdraw failed.");
    }
  };

  const handleClaimRewards = async (airlineId: string) => {
    if (!protocolQuery.data || !airlinesOnChainQuery.data || !address) return;
    const rewardsAddr = airlinesOnChainQuery.data.get(airlineId)?.rewards;
    if (!rewardsAddr) return;
    const maxRewards = maxRewardsInputs[airlineId] ? BigInt(maxRewardsInputs[airlineId]) : defaultRewardsToClaim;
    try {
      const data = encodeAbiParameters(
        [
          { name: "network", type: "address" },
          { name: "maxRewards", type: "uint256" },
          { name: "hints", type: "bytes[]" }
        ],
        [protocolQuery.data.network, maxRewards, []]
      );
      const hash = await writeContractAsync({
        address: rewardsAddr,
        abi: rewardsABI,
        functionName: "claimRewards",
        args: [address as Address, protocolQuery.data.collateral, data]
      });
      setFeedback(`Rewards claim tx submitted: ${hash}`);
    } catch (err) {
      console.error(err);
      setFeedback("Rewards claim failed.");
    }
  };

  const handleApproveVault = async (airlineId: string) => {
    if (!protocolQuery.data || !airlinesOnChainQuery.data || !address) return;
    const vault = airlinesOnChainQuery.data.get(airlineId)?.vault;
    if (!vault) return;
    try {
      const hash = await writeContractAsync({
        address: protocolQuery.data.collateral,
        abi: erc20ABI,
        functionName: "approve",
        args: [vault, maxUint256]
      });
      setFeedback(`Vault approval tx submitted: ${hash}`);
    } catch (err) {
      console.error(err);
      setFeedback("Vault approval failed.");
    }
  };

  const allowanceEnough = protocolQuery.data && allowanceQuery.data !== undefined
    ? allowanceQuery.data >= protocolQuery.data.policyPremium
    : false;

  return (
    <div className="app-shell">
      <header className="app-header">
        <div>
          <h1>Symbiotic Flight Delay Insurance</h1>
          <p>Buy coverage, monitor flights, and manage vault liquidity powered by Settlement signatures.</p>
        </div>
        <WalletStatus />
      </header>

      {feedback && <div className="feedback">{feedback}</div>}

      <section className="panel">
        <h2>Available Flights</h2>
        {flightsQuery.isLoading ? (
          <p>Loading flights...</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Airline</th>
                <th>Flight</th>
                <th>Departure</th>
                <th>API Status</th>
                <th>On-chain</th>
                <th>Policies</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {flattenedFlights.map(({ airline, flight }) => {
                const key = flightKey(airline.airlineId, flight.flightId);
                const chainData = chainFlights[key];
                const onChainStatus = statusLabel(chainData?.status ?? 0);
                const canBuy = canBuyInsurance(chainData, protocolQuery.data);
                const policyStatus = policies[key] ?? 0;
                const canClaim = policyStatus === 1 && chainData?.status === 2;
                const buyDisabled = !isConnected || !allowanceEnough || !canBuy;
                return (
                  <tr key={key}>
                    <td>{airline.name}</td>
                    <td>{flight.flightId}</td>
                    <td>{formatTimestamp(flight.departureTimestamp)}</td>
                    <td>{flight.status}</td>
                    <td>{onChainStatus}</td>
                    <td>{policyStatusLabel(policyStatus)}</td>
                    <td className="actions">
                      {canClaim ? (
                        <button onClick={() => handleClaim(airline.airlineId, flight.flightId)} disabled={!isConnected}>
                          Claim
                        </button>
                      ) : (
                        <button
                          onClick={() => handleBuy(airline.airlineId, flight.flightId)}
                          disabled={buyDisabled}
                        >
                          Buy
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
        {isConnected && protocolQuery.data && !allowanceEnough && (
          <div className="notice">
            <p>Approve the FlightDelays contract to spend your collateral before buying insurance.</p>
            <button onClick={handleApprove}>Approve {protocolQuery.data.collateralSymbol}</button>
          </div>
        )}
      </section>

      <section className="panel">
        <h2>Provider Tools</h2>
        {!isConnected ? (
          <p>Connect your wallet to manage vault liquidity and rewards.</p>
        ) : flights.length === 0 ? (
          <p>No airlines loaded yet.</p>
        ) : (
          flights.map((airline) => {
            const vaultInfo = airlinesOnChainQuery.data?.get(airline.airlineId);
            const balance = vaultBalancesQuery.data?.get(airline.airlineId) ?? 0n;
            const claimable = rewardEstimatesQuery.data?.get(airline.airlineId) ?? 0n;
            const depositValue = depositInputs[airline.airlineId] ?? "";
            const withdrawValue = withdrawInputs[airline.airlineId] ?? "";
            const maxRewardsValue = maxRewardsInputs[airline.airlineId] ?? "";
            const decimals = protocolQuery.data?.collateralDecimals ?? 18;
            const vaultAllowance = vaultAllowancesQuery.data?.get(airline.airlineId) ?? 0n;
            const desiredDeposit = (() => {
              try {
                return depositValue ? parseUnits(depositValue, decimals) : 0n;
              } catch {
                return 0n;
              }
            })();
            const needsVaultApproval = desiredDeposit > 0n && vaultAllowance < desiredDeposit;
            return (
              <div key={airline.airlineId} className="airline-card">
                <div className="airline-card__header">
                  <div>
                    <h3>{airline.name}</h3>
                    <p className="muted">Vault: {vaultInfo?.vault ?? "-"}</p>
                    <p className="muted">Rewards: {vaultInfo?.rewards ?? "-"}</p>
                  </div>
                  <div>
                    <div>
                      Staked: {protocolQuery.data
                        ? formatAmount(balance, protocolQuery.data.collateralDecimals)
                        : "0"}{" "}
                      {protocolQuery.data?.collateralSymbol}
                    </div>
                    <div>
                      Claimable rewards: {protocolQuery.data
                        ? formatAmount(claimable, protocolQuery.data.collateralDecimals)
                        : "0"}{" "}
                      {protocolQuery.data?.collateralSymbol}
                    </div>
                  </div>
                </div>
                <div className="airline-card__actions">
                  <div>
                    <label>Deposit amount</label>
                    <div className="form-row">
                      <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={depositValue}
                        onChange={(e) =>
                          setDepositInputs((prev) => ({ ...prev, [airline.airlineId]: e.target.value }))
                        }
                      />
                      <button onClick={() => handleDeposit(airline.airlineId)} disabled={!depositValue}>
                        Deposit
                      </button>
                    </div>
                    {protocolQuery.data && (
                      <p className="muted small">
                        Allowance: {formatAmount(vaultAllowance, decimals)} {protocolQuery.data.collateralSymbol}
                      </p>
                    )}
                    {needsVaultApproval && (
                      <button className="ghost-btn" onClick={() => handleApproveVault(airline.airlineId)}>
                        Approve vault spending
                      </button>
                    )}
                  </div>
                  <div>
                    <label>Withdraw amount</label>
                    <div className="form-row">
                      <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={withdrawValue}
                        onChange={(e) =>
                          setWithdrawInputs((prev) => ({ ...prev, [airline.airlineId]: e.target.value }))
                        }
                      />
                      <button onClick={() => handleWithdraw(airline.airlineId)} disabled={!withdrawValue}>
                        Withdraw
                      </button>
                    </div>
                  </div>
                  <div>
                    <label>Rewards (max batches)</label>
                    <div className="form-row">
                      <input
                        type="number"
                        min="1"
                        step="1"
                        value={maxRewardsValue}
                        onChange={(e) =>
                          setMaxRewardsInputs((prev) => ({ ...prev, [airline.airlineId]: e.target.value }))
                        }
                      />
                      <button onClick={() => handleClaimRewards(airline.airlineId)}>Claim Rewards</button>
                    </div>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </section>
    </div>
  );
}

function WalletStatus() {
  const { address, isConnected } = useAccount();
  const { connectors, connect, error, isPending } = useConnect();
  const { disconnect } = useDisconnect();
  const truncated = address ? `${address.slice(0, 6)}…${address.slice(-4)}` : "";

  if (!isConnected) {
    return (
      <div className="wallet-status">
        {connectors.map((connector) => (
          <button
            key={connector.id}
            onClick={() => connect({ connector })}
            disabled={!connector.ready || isPending}
          >
            Connect {connector.name}
          </button>
        ))}
        {error && <span className="error-text">{error.message}</span>}
      </div>
    );
  }

  return (
    <div className="wallet-status">
      <span className="muted">{truncated}</span>
      <button onClick={() => disconnect()}>Disconnect</button>
    </div>
  );
}

function statusLabel(status: number) {
  switch (status) {
    case 1:
      return "Scheduled";
    case 2:
      return "Delayed";
    case 3:
      return "Departed";
    default:
      return "Not created";
  }
}

function policyStatusLabel(status: number) {
  switch (status) {
    case 1:
      return "Purchased";
    case 2:
      return "Claimed";
    default:
      return "None";
  }
}

function canBuyInsurance(chainData?: { timestamp: bigint; status: number }, protocol?: ProtocolInfo | null) {
  if (!chainData || !protocol) return false;
  if (chainData.status !== 1) return false;
  const now = Math.floor(Date.now() / 1000);
  const ts = Number(chainData.timestamp ?? 0n);
  const policyWindowSeconds = Number(protocol.policyWindow ?? 0n);
  const delayWindowSeconds = Number(protocol.delayWindow ?? 0n);
  const lowerBound = ts - policyWindowSeconds;
  const upperBound = ts - delayWindowSeconds;
  if (upperBound <= lowerBound) return false;
  return now > lowerBound && now <= upperBound;
}
