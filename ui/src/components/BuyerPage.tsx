import type { ReactNode } from "react";

import { airlineImages } from "../constants/airlines";
import type { ChainFlightData, FlattenedFlight, PolicyMap, ProtocolInfo } from "../types";
import { formatTimestamp } from "../utils/format";
import { flightKey } from "../utils/hash";

export function BuyerPage({
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
              const buyDisabled = !isConnected || !allowanceEnough || buyWindow?.state !== "open" || hasPolicy;

              let actionNode: ReactNode = null;
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
