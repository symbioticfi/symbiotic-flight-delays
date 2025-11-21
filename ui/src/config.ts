import { defineChain } from "viem";
import type { Address } from "viem";

const chainId = Number(import.meta.env.VITE_CHAIN_ID ?? 31337);
const chainName = import.meta.env.VITE_CHAIN_NAME ?? "Symbiotic Local";
export const rpcUrl = import.meta.env.VITE_RPC_URL ?? "http://127.0.0.1:8545";

export const flightDelaysAddress = (import.meta.env.VITE_FLIGHT_DELAYS_ADDRESS ??
  "0x0000000000000000000000000000000000000000") as Address;
export const flightsApiUrl = import.meta.env.VITE_FLIGHTS_API_URL ?? "http://127.0.0.1:8085";

export const chain = defineChain({
  id: chainId,
  name: chainName,
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: [rpcUrl] }, public: { http: [rpcUrl] } }
});
