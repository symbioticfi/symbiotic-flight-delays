export type FlightAPIStatus = "SCHEDULED" | "DELAYED" | "DEPARTED";

export interface AirlineDTO {
  airlineId: string;
  name: string;
  code: string;
}

export interface FlightDTO {
  airlineId: string;
  flightId: string;
  departureTimestamp: number;
  status: FlightAPIStatus;
}

export interface AirlineWithFlights extends AirlineDTO {
  flights: FlightDTO[];
}

export interface FlightChainState {
  timestamp: bigint;
  status: number;
  policiesSold: bigint;
}

export interface PolicyState {
  status: number;
}

export interface ProtocolInfo {
  policyPremium: bigint;
  policyPayout: bigint;
  policyWindow: bigint;
  delayWindow: bigint;
  collateral: `0x${string}`;
  collateralSymbol: string;
  collateralDecimals: number;
  network: `0x${string}`;
}

export type ChainFlightData = Record<string, FlightChainState>;

export type PolicyMap = Record<string, number>;

export type FlattenedFlight = {
  airline: AirlineWithFlights;
  flight: FlightDTO;
};
