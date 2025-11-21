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
