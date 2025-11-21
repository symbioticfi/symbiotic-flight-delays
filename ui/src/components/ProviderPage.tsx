import type { Dispatch, SetStateAction } from "react";
import type { Address } from "viem";
import { parseUnits } from "viem";

import { airlineImages } from "../constants/airlines";
import type { AirlineWithFlights, ProtocolInfo } from "../types";
import { formatAmount } from "../utils/format";

export function ProviderPage({
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
