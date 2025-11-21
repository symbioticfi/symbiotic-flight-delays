import { defineChain } from "viem";
import type { Address } from "viem";

const isBrowser = typeof window !== "undefined";
const fallbackHost = isBrowser ? window.location.hostname || "127.0.0.1" : "127.0.0.1";
const fallbackProtocol = isBrowser ? window.location.protocol : "http:";

const chainId = Number(import.meta.env.VITE_CHAIN_ID ?? 31337);
const chainName = import.meta.env.VITE_CHAIN_NAME ?? "Symbiotic Local";

const defaultRpcUrl = `${fallbackProtocol}//${fallbackHost}:8545`;
const configuredRpcUrl = import.meta.env.VITE_RPC_URL ?? defaultRpcUrl;

const defaultFlightsApiUrl = `${fallbackProtocol}//${fallbackHost}:8085`;
const configuredFlightsApiUrl = import.meta.env.VITE_FLIGHTS_API_URL ?? defaultFlightsApiUrl;

const dockerHostnames = new Set(["flights-api", "anvil"]);

function resolveUrl(raw: string, fallbackPort?: number) {
  if (!isBrowser) {
    return raw;
  }
  try {
    const parsed = new URL(raw, `${fallbackProtocol}//${fallbackHost}`);
    if (dockerHostnames.has(parsed.hostname)) {
      parsed.hostname = fallbackHost;
      if (fallbackPort && !parsed.port) {
        parsed.port = String(fallbackPort);
      }
      parsed.protocol = fallbackProtocol;
    }
    return parsed.toString();
  } catch {
    if (fallbackPort) {
      return `${fallbackProtocol}//${fallbackHost}:${fallbackPort}`;
    }
    return `${fallbackProtocol}//${fallbackHost}`;
  }
}

export const rpcUrl = resolveUrl(configuredRpcUrl, 8545);

export const flightDelaysAddress = (import.meta.env.VITE_FLIGHT_DELAYS_ADDRESS ??
  "0x0000000000000000000000000000000000000000") as Address;
export const flightsApiUrl = resolveUrl(configuredFlightsApiUrl, 8085);

export const chain = defineChain({
  id: chainId,
  name: chainName,
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: [rpcUrl] }, public: { http: [rpcUrl] } },
});
