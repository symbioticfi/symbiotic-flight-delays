import { Dispatch, SetStateAction, useCallback, useEffect, useMemo, useState } from "react";
import { NavLink, Route, Routes } from "react-router-dom";
import { useAccount, usePublicClient, useWriteContract, useReadContracts } from "wagmi";
import { useQuery } from "@tanstack/react-query";
import {
  Address,
  Hex,
  BaseError,
  encodeAbiParameters,
  maxUint256,
  parseUnits,
  zeroAddress,
  createPublicClient,
  http,
} from "viem";
import type { Abi, PublicClient } from "viem";
import { useAppKit } from "@reown/appkit/library/react";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

import flightDelaysAbi from "./abi/FlightDelays.json";
import erc20Abi from "./abi/ERC20.json";
import vaultAbi from "./abi/Vault.json";
import rewardsAbi from "./abi/Rewards.json";
import { chain, flightsApiUrl, flightDelaysAddress, rpcUrl } from "./config";
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

type ToastMessages = {
  pending: string;
  success: string;
  error: string;
};

function formatWriteError(error: unknown) {
  if (error && typeof error === "object") {
    const baseError = error as Partial<BaseError> & Partial<Error>;
    if (typeof baseError.shortMessage === "string") {
      return baseError.shortMessage;
    }
    if (typeof baseError.message === "string") {
      return baseError.message;
    }
  }
  return "Check console for more details.";
}

type MaybeContractResult<T> = { status: "success"; result: T } | { status: "failure"; error: Error } | T | undefined;

function unwrapResult<T>(entry: MaybeContractResult<T>): T | undefined {
  if (entry === undefined || entry === null) return undefined;
  if (typeof entry === "object" && entry !== null && "status" in entry) {
    const typed = entry as { status?: string; result?: T };
    return typed.status === "success" ? (typed["result"] as T | undefined) : undefined;
  }
  return entry as T;
}

const airlineImages: Record<string, string> = {
  ALPHA: "/alpha-air.png",
  BETA: "/beta-wings.png",
  GAMMA: "/gamma-connect.png",
};

type FlattenedFlight = {
  airline: AirlineWithFlights;
  flight: AirlineWithFlights["flights"][number];
};

export default function App() {
  const { address, isConnected } = useAccount();
  const wagmiPublicClient = usePublicClient({ chainId: chain.id });
  const fallbackClient = useMemo(() => createPublicClient({ chain, transport: http(rpcUrl) }), []);
  const publicClient: PublicClient = wagmiPublicClient ?? fallbackClient;
  const { writeContract } = useWriteContract();
  const { open } = useAppKit();

  const flightsQuery = useQuery({
    queryKey: ["flights", flightsApiUrl],
    queryFn: () => fetchAirlinesWithFlights(flightsApiUrl),
    refetchInterval: 10_000,
  });

  const flights = flightsQuery.data ?? [];
  const [now, setNow] = useState(() => Math.floor(Date.now() / 1000));

  useEffect(() => {
    const id = setInterval(() => setNow(Math.floor(Date.now() / 1000)), 1_000);
    return () => clearInterval(id);
  }, []);

  const submitWrite = useCallback(
    async (variables: Parameters<typeof writeContract>[0], messages: ToastMessages) => {
      const txPromise = new Promise<Hex>((resolve, reject) => {
        writeContract(variables, {
          onSuccess: (hash) => resolve(hash as Hex),
          onError: (error) => reject(error),
        });
      });
      toast.promise(txPromise, {
        pending: messages.pending,
        success: {
          render({ data }) {
            return `${messages.success}: ${data}`;
          },
        },
        error: {
          render({ data }) {
            return `${messages.error}: ${formatWriteError(data)}`;
          },
        },
      });
      return txPromise;
    },
    [writeContract],
  );

  const requestWalletConnection = async () => {
    if (address) {
      return true;
    }
    try {
      await open?.({ view: "Connect" });
    } catch (err) {
      console.error("wallet modal error", err);
    }
    return false;
  };

  const flattenedFlights = useMemo(
    () =>
      flights.flatMap((airline) =>
        airline.flights.map((flight) => ({
          airline,
          flight,
        })),
      ),
    [flights],
  );

  const protocolBaseContracts = useMemo(
    () => [
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyPremium" } as const,
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyPayout" } as const,
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "policyWindow" } as const,
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "delayWindow" } as const,
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "collateral" } as const,
      { address: flightDelaysAddress, abi: flightDelaysABI, functionName: "NETWORK" } as const,
    ],
    [],
  );

  const protocolBaseReads = useReadContracts({
    contracts: protocolBaseContracts,
    allowFailure: false,
    query: {
      refetchInterval: 15_000,
    },
  });

  const collateralAddress = useMemo(() => {
    const baseResults = protocolBaseReads.data as readonly unknown[] | undefined;
    return ((baseResults?.[4] as Address | undefined) ?? zeroAddress) as Address;
  }, [protocolBaseReads.data]);

  const collateralContracts = useMemo(() => {
    if (!collateralAddress) {
      return [];
    }
    return [
      { address: collateralAddress, abi: erc20ABI, functionName: "symbol" } as const,
      { address: collateralAddress, abi: erc20ABI, functionName: "decimals" } as const,
    ];
  }, [collateralAddress]);

  const collateralReads = useReadContracts({
    contracts: collateralContracts,
    allowFailure: false,
    query: {
      enabled: collateralContracts.length > 0,
      refetchInterval: 30_000,
    },
  });

  const protocol: ProtocolInfo | null = useMemo(() => {
    const base = protocolBaseReads.data as readonly unknown[] | undefined;
    if (!base) return null;
    const policyPremium = (base[0] as bigint | undefined) ?? 0n;
    const policyPayout = (base[1] as bigint | undefined) ?? 0n;
    const policyWindow = (base[2] as bigint | undefined) ?? 0n;
    const delayWindow = (base[3] as bigint | undefined) ?? 0n;
    const collateral = (base[4] as Address | undefined) ?? zeroAddress;
    const network = (base[5] as Address | undefined) ?? zeroAddress;
    const collateralData = collateralReads.data as readonly unknown[] | undefined;
    const symbol = (collateralData?.[0] as string | undefined) ?? "TOKEN";
    const decimals = Number((collateralData?.[1] as number | bigint | undefined) ?? 18);
    return {
      policyPremium,
      policyPayout,
      policyWindow,
      delayWindow,
      collateral,
      collateralSymbol: symbol,
      collateralDecimals: decimals,
      network,
    };
  }, [protocolBaseReads.data, collateralReads.data]);

  const flightContracts = useMemo(
    () =>
      flattenedFlights.map(({ airline, flight }) => ({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "flights",
        args: [hashIdentifier(airline.airlineId), hashIdentifier(flight.flightId)] as const,
      })),
    [flattenedFlights],
  );

  const flightReads = useReadContracts({
    contracts: flightContracts,
    allowFailure: true,
    query: {
      enabled: flightContracts.length > 0,
      refetchInterval: 6_000,
    },
  });

  const chainFlights = useMemo(() => {
    const map: ChainFlightData = {};
    const data = flightReads.data;
    if (!data) return map;
    data.forEach((resp, idx) => {
      const key = flightKey(flattenedFlights[idx].airline.airlineId, flattenedFlights[idx].flight.flightId);
      const result = unwrapResult<[bigint, number, bigint]>(resp as MaybeContractResult<[bigint, number, bigint]>);
      if (!result) return;
      const [timestamp, status, policiesSold] = result;
      map[key] = { timestamp, status: Number(status), policiesSold };
    });
    return map;
  }, [flightReads.data, flattenedFlights]);

  const policyContracts = useMemo(() => {
    if (!address) return [];
    return flattenedFlights.map(({ airline, flight }) => ({
      address: flightDelaysAddress,
      abi: flightDelaysABI,
      functionName: "policies",
      args: [hashIdentifier(airline.airlineId), hashIdentifier(flight.flightId), address] as const,
    }));
  }, [address, flattenedFlights]);

  const policyReads = useReadContracts({
    contracts: policyContracts,
    allowFailure: true,
    query: {
      enabled: policyContracts.length > 0,
      refetchInterval: 6_000,
    },
  });

  const policies = useMemo(() => {
    const map: PolicyMap = {};
    const data = policyReads.data;
    if (!data) return map;
    data.forEach((resp, idx) => {
      const contract = policyContracts[idx];
      if (!contract) return;
      const key = flightKey(flattenedFlights[idx].airline.airlineId, flattenedFlights[idx].flight.flightId);
      const policyState = unwrapResult<number | bigint>(resp as MaybeContractResult<number | bigint>);
      if (policyState === undefined) return;
      map[key] = Number(policyState);
    });
    return map;
  }, [policyReads.data, policyContracts, flattenedFlights]);

  const airlineContracts = useMemo(
    () =>
      flights.map((airline) => ({
        address: flightDelaysAddress,
        abi: flightDelaysABI,
        functionName: "airlines",
        args: [hashIdentifier(airline.airlineId)] as const,
      })),
    [flights],
  );

  const airlineReads = useReadContracts({
    contracts: airlineContracts,
    allowFailure: true,
    query: {
      enabled: airlineContracts.length > 0,
      refetchInterval: 12_000,
    },
  });

  const airlinesOnChain = useMemo(() => {
    const map = new Map<string, { vault: Address; rewards: Address }>();
    const data = airlineReads.data;
    if (!data) return map;
    data.forEach((resp, idx) => {
      const tuple = unwrapResult<[Address, Address, bigint, Hex]>(
        resp as MaybeContractResult<[Address, Address, bigint, Hex]>,
      );
      if (!tuple) return;
      const [vault, rewards] = tuple;
      map.set(flights[idx].airlineId, { vault, rewards });
    });
    return map;
  }, [airlineReads.data, flights]);

  const policyAllowanceContracts = useMemo(() => {
    if (!address || !protocol) return [];
    return [
      {
        address: protocol.collateral,
        abi: erc20ABI,
        functionName: "allowance",
        args: [address, flightDelaysAddress] as const,
      },
    ];
  }, [address, protocol]);

  const policyAllowanceReads = useReadContracts({
    contracts: policyAllowanceContracts,
    allowFailure: false,
    query: {
      enabled: policyAllowanceContracts.length > 0,
      refetchInterval: 8_000,
    },
  });

  const allowanceValue = useMemo(() => {
    const data = policyAllowanceReads.data as readonly unknown[] | undefined;
    return (data?.[0] as bigint | undefined) ?? 0n;
  }, [policyAllowanceReads.data]);

  const vaultAllowanceContracts = useMemo(() => {
    if (!address || !protocol || airlinesOnChain.size === 0) return [];
    const entries: { airlineId: string; contract: any }[] = [];
    airlinesOnChain.forEach((info, airlineId) => {
      entries.push({
        airlineId,
        contract: {
          address: protocol.collateral,
          abi: erc20ABI,
          functionName: "allowance",
          args: [address, info.vault] as const,
        },
      });
    });
    return entries;
  }, [address, protocol, airlinesOnChain]);

  const vaultAllowanceReads = useReadContracts({
    contracts: vaultAllowanceContracts.map((entry) => entry.contract),
    allowFailure: true,
    query: {
      enabled: vaultAllowanceContracts.length > 0,
      refetchInterval: 12_000,
    },
  });

  const vaultAllowances = useMemo(() => {
    const map = new Map<string, bigint>();
    const data = vaultAllowanceReads.data;
    if (!data) return map;
    data.forEach((resp, idx) => {
      const airlineId = vaultAllowanceContracts[idx]?.airlineId;
      if (!airlineId) return;
      const allowance = unwrapResult<bigint>(resp as MaybeContractResult<bigint>);
      if (allowance === undefined) return;
      map.set(airlineId, allowance);
    });
    return map;
  }, [vaultAllowanceReads.data, vaultAllowanceContracts]);

  const vaultBalanceContracts = useMemo(() => {
    if (!address || airlinesOnChain.size === 0) return [];
    const entries: { airlineId: string; contract: any }[] = [];
    airlinesOnChain.forEach((info, airlineId) => {
      entries.push({
        airlineId,
        contract: {
          address: info.vault,
          abi: vaultABI,
          functionName: "activeBalanceOf",
          args: [address] as const,
        },
      });
    });
    return entries;
  }, [address, airlinesOnChain]);

  const vaultBalanceReads = useReadContracts({
    contracts: vaultBalanceContracts.map((entry) => entry.contract),
    allowFailure: true,
    query: {
      enabled: vaultBalanceContracts.length > 0,
      refetchInterval: 10_000,
    },
  });

  const vaultBalances = useMemo(() => {
    const balances = new Map<string, bigint>();
    const data = vaultBalanceReads.data;
    if (!data) return balances;
    data.forEach((resp, idx) => {
      const airlineId = vaultBalanceContracts[idx]?.airlineId;
      if (!airlineId) return;
      const balance = unwrapResult<bigint>(resp as MaybeContractResult<bigint>);
      if (balance === undefined) return;
      balances.set(airlineId, balance);
    });
    return balances;
  }, [vaultBalanceReads.data, vaultBalanceContracts]);

  const rewardContracts = useMemo(() => {
    if (!address || !protocol || airlinesOnChain.size === 0) return [];
    const entries: { airlineId: string; contract: any }[] = [];
    airlinesOnChain.forEach((info, airlineId) => {
      const encoded = encodeAbiParameters(
        [
          { name: "network", type: "address" },
          { name: "maxRewards", type: "uint256" },
          { name: "hints", type: "bytes[]" },
        ],
        [protocol.network, defaultRewardsToClaim, []],
      );
      entries.push({
        airlineId,
        contract: {
          address: info.rewards,
          abi: rewardsABI,
          functionName: "claimable",
          args: [protocol.collateral, address!, encoded] as const,
        },
      });
    });
    return entries;
  }, [address, protocol, airlinesOnChain]);

  const rewardReads = useReadContracts({
    contracts: rewardContracts.map((entry) => entry.contract),
    allowFailure: true,
    query: {
      enabled: rewardContracts.length > 0,
      refetchInterval: 15_000,
    },
  });

  const rewardEstimates = useMemo(() => {
    const map = new Map<string, bigint>();
    const data = rewardReads.data;
    if (!data) return map;
    data.forEach((resp, idx) => {
      const airlineId = rewardContracts[idx]?.airlineId;
      if (!airlineId) return;
      const rewardValue = unwrapResult<bigint>(resp as MaybeContractResult<bigint>);
      if (rewardValue === undefined) return;
      map.set(airlineId, rewardValue);
    });
    return map;
  }, [rewardReads.data, rewardContracts]);

  const [depositInputs, setDepositInputs] = useState<Record<string, string>>({});
  const [withdrawInputs, setWithdrawInputs] = useState<Record<string, string>>({});
  const [maxRewardsInputs, setMaxRewardsInputs] = useState<Record<string, string>>({});

  const handleApprove = async () => {
    if (!protocol) return;
    if (!address) {
      await requestWalletConnection();
      return;
    }
    try {
      await submitWrite(
        {
          address: protocol.collateral,
          abi: erc20ABI,
          functionName: "approve",
          args: [flightDelaysAddress, maxUint256],
        },
        {
          pending: "Submitting collateral approval...",
          success: "Collateral approval sent",
          error: "Collateral approval failed",
        },
      );
    } catch (err) {
      console.error("policy approval error", err);
    }
  };

  const handleBuy = async (airlineId: string, flightId: string) => {
    if (!address) {
      await requestWalletConnection();
      return;
    }
    try {
      await submitWrite(
        {
          address: flightDelaysAddress,
          abi: flightDelaysABI,
          functionName: "buyInsurance",
          args: [hashIdentifier(airlineId), hashIdentifier(flightId)],
        },
        {
          pending: `Submitting buy for ${flightId}...`,
          success: "Buy transaction sent",
          error: "Buy failed",
        },
      );
    } catch (err) {
      console.error("buy error", err);
    }
  };

  const handleClaim = async (airlineId: string, flightId: string) => {
    if (!address) {
      await requestWalletConnection();
      return;
    }
    try {
      await submitWrite(
        {
          address: flightDelaysAddress,
          abi: flightDelaysABI,
          functionName: "claimInsurance",
          args: [hashIdentifier(airlineId), hashIdentifier(flightId)],
        },
        {
          pending: `Submitting claim for ${flightId}...`,
          success: "Claim transaction sent",
          error: "Claim failed",
        },
      );
    } catch (err) {
      console.error("claim error", err);
    }
  };

  const handleDeposit = async (airlineId: string) => {
    if (!protocol) return;
    if (!address) {
      await requestWalletConnection();
      return;
    }
    const amount = depositInputs[airlineId];
    if (!amount) return;
    try {
      const parsed = parseUnits(amount, protocol.collateralDecimals);
      const vault = airlinesOnChain.get(airlineId)?.vault;
      if (!vault) return;
      await submitWrite(
        {
          address: vault,
          abi: vaultABI,
          functionName: "deposit",
          args: [address as Address, parsed],
        },
        {
          pending: `Depositing to ${airlineId} vault...`,
          success: "Deposit transaction sent",
          error: "Deposit failed",
        },
      );
      setDepositInputs((prev) => ({ ...prev, [airlineId]: "" }));
    } catch (err) {
      console.error("deposit error", err);
    }
  };

  const handleWithdraw = async (airlineId: string) => {
    if (!protocol) return;
    if (!address) {
      await requestWalletConnection();
      return;
    }
    const amount = withdrawInputs[airlineId];
    if (!amount) return;
    try {
      const parsed = parseUnits(amount, protocol.collateralDecimals);
      const vault = airlinesOnChain.get(airlineId)?.vault;
      if (!vault) return;
      await submitWrite(
        {
          address: vault,
          abi: vaultABI,
          functionName: "withdraw",
          args: [address as Address, parsed],
        },
        {
          pending: `Withdrawing from ${airlineId} vault...`,
          success: "Withdraw transaction sent",
          error: "Withdraw failed",
        },
      );
      setWithdrawInputs((prev) => ({ ...prev, [airlineId]: "" }));
    } catch (err) {
      console.error("withdraw error", err);
    }
  };

  const handleClaimRewards = async (airlineId: string) => {
    if (!protocol) return;
    if (!address) {
      await requestWalletConnection();
      return;
    }
    const rewardsAddr = airlinesOnChain.get(airlineId)?.rewards;
    if (!rewardsAddr) return;
    const maxRewards = maxRewardsInputs[airlineId] ? BigInt(maxRewardsInputs[airlineId]) : defaultRewardsToClaim;
    try {
      const data = encodeAbiParameters(
        [
          { name: "network", type: "address" },
          { name: "maxRewards", type: "uint256" },
          { name: "hints", type: "bytes[]" },
        ],
        [protocol.network, maxRewards, []],
      );
      await submitWrite(
        {
          address: rewardsAddr,
          abi: rewardsABI,
          functionName: "claimRewards",
          args: [address as Address, protocol.collateral, data],
        },
        {
          pending: `Claiming rewards for ${airlineId}...`,
          success: "Rewards claim transaction sent",
          error: "Rewards claim failed",
        },
      );
    } catch (err) {
      console.error("claim rewards error", err);
    }
  };

  const handleApproveVault = async (airlineId: string) => {
    if (!protocol) return;
    if (!address) {
      await requestWalletConnection();
      return;
    }
    const vault = airlinesOnChain.get(airlineId)?.vault;
    if (!vault) return;
    try {
      await submitWrite(
        {
          address: protocol.collateral,
          abi: erc20ABI,
          functionName: "approve",
          args: [vault, maxUint256],
        },
        {
          pending: `Approving ${airlineId} vault...`,
          success: "Vault approval transaction sent",
          error: "Vault approval failed",
        },
      );
    } catch (err) {
      console.error("vault approval error", err);
    }
  };

  const allowanceEnough = protocol ? allowanceValue >= protocol.policyPremium : false;

  const collateralBalanceContracts = useMemo(() => {
    if (!address || !protocol) return [];
    return [
      {
        address: protocol.collateral,
        abi: erc20ABI,
        functionName: "balanceOf",
        args: [address] as const,
      },
    ];
  }, [address, protocol]);

  const collateralBalanceReads = useReadContracts({
    contracts: collateralBalanceContracts,
    allowFailure: false,
    query: {
      enabled: collateralBalanceContracts.length > 0,
      refetchInterval: 12_000,
    },
  });

  const collateralBalance = useMemo(() => {
    const data = collateralBalanceReads.data as readonly unknown[] | undefined;
    return (data?.[0] as bigint | undefined) ?? 0n;
  }, [collateralBalanceReads.data]);

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-title">
          <img src="/logo.png" alt="Symbiotic Flight Delay Insurance" className="app-logo" />
          <div>
            <h1>Symbiotic Flight Delay Insurance</h1>
            <p>Buy coverage, monitor flights, and manage vault liquidity powered by Settlement signatures.</p>
          </div>
        </div>
        <WalletStatus protocol={protocol} isConnected={isConnected} collateralBalance={collateralBalance} />
      </header>
      <nav className="main-nav">
        <NavLink to="/" end>
          Buy Coverage
        </NavLink>
        <NavLink to="/providers">Provide Coverage</NavLink>
      </nav>

      <Routes>
        <Route
          path="/"
          element={
            <BuyerPage
              flightsLoading={flightsQuery.isLoading}
              flattenedFlights={flattenedFlights}
              chainFlights={chainFlights}
              policies={policies}
              protocol={protocol}
              isConnected={isConnected}
              allowanceEnough={allowanceEnough}
              handleApprove={handleApprove}
              handleBuy={handleBuy}
              handleClaim={handleClaim}
              now={now}
            />
          }
        />
        <Route
          path="/providers"
          element={
            <ProviderPage
              flights={flights}
              isConnected={isConnected}
              protocol={protocol}
              airlinesOnChain={airlinesOnChain}
              vaultAllowances={vaultAllowances}
              vaultBalances={vaultBalances}
              rewardEstimates={rewardEstimates}
              depositInputs={depositInputs}
              setDepositInputs={setDepositInputs}
              withdrawInputs={withdrawInputs}
              setWithdrawInputs={setWithdrawInputs}
              maxRewardsInputs={maxRewardsInputs}
              setMaxRewardsInputs={setMaxRewardsInputs}
              handleDeposit={handleDeposit}
              handleWithdraw={handleWithdraw}
              handleApproveVault={handleApproveVault}
              handleClaimRewards={handleClaimRewards}
            />
          }
        />
      </Routes>
      <ToastContainer position="bottom-right" newestOnTop closeOnClick pauseOnFocusLoss={false} />
    </div>
  );
}

function WalletStatus({
  protocol,
  isConnected,
  collateralBalance,
}: {
  protocol: ProtocolInfo | null;
  isConnected: boolean;
  collateralBalance: bigint;
}) {
  const symbol = protocol?.collateralSymbol ?? "TOKEN";
  const decimals = protocol?.collateralDecimals ?? 18;
  const formattedBalance = (() => {
    if (!isConnected) return "Connect wallet";
    if (!protocol) return `-- ${symbol}`;
    return `${formatAmount(collateralBalance, decimals, 4)} ${symbol}`;
  })();

  return (
    <div className="wallet-status">
      <div className="balance-chip">
        <span className="balance-chip__label">{symbol} balance</span>
        <span className="balance-chip__value">{formattedBalance}</span>
      </div>
      <appkit-button />
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

interface BuyerPageProps {
  flightsLoading: boolean;
  flattenedFlights: FlattenedFlight[];
  chainFlights: ChainFlightData;
  policies: PolicyMap;
  protocol: ProtocolInfo | null;
  isConnected: boolean;
  allowanceEnough: boolean;
  handleApprove: () => Promise<void>;
  handleBuy: (airlineId: string, flightId: string) => Promise<void>;
  handleClaim: (airlineId: string, flightId: string) => Promise<void>;
  now: number;
}

function BuyerPage({
  flightsLoading,
  flattenedFlights,
  chainFlights,
  policies,
  protocol,
  isConnected,
  allowanceEnough,
  handleApprove,
  handleBuy,
  handleClaim,
  now,
}: BuyerPageProps) {
  return (
    <section className="panel">
      <h2>Available Flights</h2>
      {flightsLoading ? (
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
              const policyStatus = policies[key] ?? 0;
              const hasPolicy = policyStatus === 1 || policyStatus === 2;
              const chainStatus = chainData?.status ?? (flight.departureTimestamp > now ? 1 : 0);
              const onChainStatus = statusLabel(chainStatus);
              const buyWindow = getBuyWindowInfo(chainData, protocol, now, flight.departureTimestamp);
              const departureTs = chainData ? Number(chainData.timestamp ?? 0n) : flight.departureTimestamp;
              const departureCountdown = Math.max(0, departureTs - now);
              const canClaim = policyStatus === 1 && chainStatus === 2;
              const buyDisabled = !isConnected || !allowanceEnough || buyWindow?.state !== "open" || hasPolicy;

              let actionNode: React.ReactNode = null;
              let timerLabel: string | null = null;
              let timerClass = "after";

              if (chainStatus === 2) {
                if (policyStatus === 2) {
                  actionNode = <button disabled>Claimed</button>;
                } else if (policyStatus === 1) {
                  actionNode = (
                    <button onClick={() => handleClaim(airline.airlineId, flight.flightId)} disabled={!isConnected}>
                      Claim
                    </button>
                  );
                } else {
                  actionNode = <>Delayed</>;
                }
                timerLabel = "Flight delayed";
              } else if (chainStatus === 3) {
                actionNode = <>Departed</>;
              } else {
                if (!protocol) {
                  actionNode = <button disabled>Buy</button>;
                  timerLabel = `Departs in ${formatCountdown(departureCountdown)}`;
                } else if (buyWindow) {
                  switch (buyWindow.state) {
                    case "before":
                      actionNode = <button disabled>Buy</button>;
                      timerLabel = `Opens in ${formatCountdown(buyWindow.seconds)}`;
                      timerClass = "before";
                      break;
                    case "open":
                      if (hasPolicy) {
                        actionNode = <button disabled>Purchased</button>;
                        timerLabel = `Closes in ${formatCountdown(buyWindow.seconds)}`;
                        timerClass = "open";
                      } else {
                        actionNode = (
                          <button onClick={() => handleBuy(airline.airlineId, flight.flightId)} disabled={buyDisabled}>
                            Buy
                          </button>
                        );
                        timerLabel = `Closes in ${formatCountdown(buyWindow.seconds)}`;
                        timerClass = "open";
                      }
                      break;
                    case "after":
                      actionNode = <button disabled>Buy</button>;
                      timerLabel = `Departs in ${formatCountdown(buyWindow.seconds)}`;
                      timerClass = "after";
                      break;
                  }
                } else {
                  actionNode = <button disabled>Buy</button>;
                  timerLabel = `Departs in ${formatCountdown(departureCountdown)}`;
                }
              }
              return (
                <tr key={key}>
                  <td>
                    <div className="airline-cell">
                      <img
                        src={airlineImages[airline.airlineId.toUpperCase()] ?? "/logo.png"}
                        alt={airline.name}
                        className="airline-thumb"
                      />
                      <span>{airline.name}</span>
                    </div>
                  </td>
                  <td>{flight.flightId}</td>
                  <td>{formatTimestamp(flight.departureTimestamp)}</td>
                  <td>{flight.status}</td>
                  <td>{onChainStatus}</td>
                  <td>{policyStatusLabel(policyStatus)}</td>
                  <td className="actions">
                    {actionNode}
                    {timerLabel && <span className={`buy-window ${timerClass}`}>{timerLabel}</span>}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
      {isConnected && protocol && !allowanceEnough && (
        <div className="notice">
          <p>Approve the FlightDelays contract to spend your collateral before buying insurance.</p>
          <button onClick={handleApprove}>Approve {protocol.collateralSymbol}</button>
        </div>
      )}
    </section>
  );
}

interface ProviderPageProps {
  flights: AirlineWithFlights[];
  isConnected: boolean;
  protocol: ProtocolInfo | null;
  airlinesOnChain: Map<string, { vault: Address; rewards: Address }>;
  vaultAllowances: Map<string, bigint>;
  vaultBalances: Map<string, bigint>;
  rewardEstimates: Map<string, bigint>;
  depositInputs: Record<string, string>;
  withdrawInputs: Record<string, string>;
  maxRewardsInputs: Record<string, string>;
  setDepositInputs: Dispatch<SetStateAction<Record<string, string>>>;
  setWithdrawInputs: Dispatch<SetStateAction<Record<string, string>>>;
  setMaxRewardsInputs: Dispatch<SetStateAction<Record<string, string>>>;
  handleDeposit: (airlineId: string) => Promise<void>;
  handleWithdraw: (airlineId: string) => Promise<void>;
  handleApproveVault: (airlineId: string) => Promise<void>;
  handleClaimRewards: (airlineId: string) => Promise<void>;
}

function ProviderPage({
  flights,
  isConnected,
  protocol,
  airlinesOnChain,
  vaultAllowances,
  vaultBalances,
  rewardEstimates,
  depositInputs,
  withdrawInputs,
  maxRewardsInputs,
  setDepositInputs,
  setWithdrawInputs,
  setMaxRewardsInputs,
  handleDeposit,
  handleWithdraw,
  handleApproveVault,
  handleClaimRewards,
}: ProviderPageProps) {
  return (
    <section className="panel">
      <h2>Provider Tools</h2>
      {!isConnected ? (
        <p>Connect your wallet to manage vault liquidity and rewards.</p>
      ) : flights.length === 0 ? (
        <p>No airlines loaded yet.</p>
      ) : (
        flights.map((airline) => {
          const vaultInfo = airlinesOnChain.get(airline.airlineId);
          const balance = vaultBalances.get(airline.airlineId) ?? 0n;
          const claimable = rewardEstimates.get(airline.airlineId) ?? 0n;
          const depositValue = depositInputs[airline.airlineId] ?? "";
          const withdrawValue = withdrawInputs[airline.airlineId] ?? "";
          const maxRewardsValue = maxRewardsInputs[airline.airlineId] ?? "";
          const decimals = protocol?.collateralDecimals ?? 18;
          const vaultAllowance = vaultAllowances.get(airline.airlineId) ?? 0n;
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
                <div className="airline-card__brand">
                  <img
                    src={airlineImages[airline.airlineId.toUpperCase()] ?? "/logo.png"}
                    alt={airline.name}
                    className="airline-thumb large"
                  />
                  <div>
                    <h3>{airline.name}</h3>
                    <p className="muted">Vault: {vaultInfo?.vault ?? "-"}</p>
                    <p className="muted">Rewards: {vaultInfo?.rewards ?? "-"}</p>
                  </div>
                </div>
                <div>
                  <div>
                    Staked: {protocol ? formatAmount(balance, protocol.collateralDecimals) : "0"}{" "}
                    {protocol?.collateralSymbol}
                  </div>
                  <div>
                    Claimable rewards: {protocol ? formatAmount(claimable, protocol.collateralDecimals) : "0"}{" "}
                    {protocol?.collateralSymbol}
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
                      onChange={(e) => setDepositInputs((prev) => ({ ...prev, [airline.airlineId]: e.target.value }))}
                    />
                    <button onClick={() => handleDeposit(airline.airlineId)} disabled={!depositValue}>
                      Deposit
                    </button>
                  </div>
                  {protocol && (
                    <p className="muted small">
                      Allowance: {formatAmount(vaultAllowance, decimals)} {protocol.collateralSymbol}
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
                      onChange={(e) => setWithdrawInputs((prev) => ({ ...prev, [airline.airlineId]: e.target.value }))}
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
  );
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

function getBuyWindowInfo(
  chainData: { timestamp: bigint; status: number } | undefined,
  protocol: ProtocolInfo | null,
  now: number,
  fallbackDeparture?: number,
) {
  if (!protocol) return null;
  const status = chainData?.status ?? (fallbackDeparture ? 1 : 0);
  if (status !== 1) return null;
  const ts = chainData ? Number(chainData.timestamp ?? 0n) : fallbackDeparture ?? 0;
  if (!ts) return null;
  const openAt = ts - Number(protocol.policyWindow ?? 0n);
  const closeAt = ts - Number(protocol.delayWindow ?? 0n);
  if (closeAt <= openAt) return null;
  if (now < openAt) {
    return { state: "before" as const, seconds: openAt - now };
  }
  if (now >= openAt && now <= closeAt) {
    return { state: "open" as const, seconds: closeAt - now };
  }
  if (now < ts) {
    return { state: "after" as const, seconds: ts - now };
  }
  return { state: "after" as const, seconds: 0 };
}

function formatCountdown(seconds: number) {
  const clamped = Math.max(0, Math.floor(seconds));
  const minutes = Math.floor(clamped / 60);
  const secs = clamped % 60;
  if (minutes <= 0) {
    return `${secs}s`;
  }
  return `${minutes}m ${secs.toString().padStart(2, "0")}s`;
}
