import { formatUnits } from "viem";

export function formatAmount(value: bigint, decimals: number, precision = 2) {
  const formatted = Number(formatUnits(value, decimals));
  return formatted.toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: precision,
  });
}

export function formatTimestamp(timestamp: number) {
  const date = new Date(timestamp * 1000);
  return date.toLocaleString();
}
