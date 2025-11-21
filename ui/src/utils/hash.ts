import { keccak256, stringToBytes } from "viem";

export function hashIdentifier(value: string) {
  const normalized = value.trim().toUpperCase();
  return keccak256(stringToBytes(normalized));
}

export function flightKey(airlineId: string, flightId: string) {
  return `${airlineId}::${flightId}`;
}
