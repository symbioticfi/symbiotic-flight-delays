import type { ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WagmiProvider } from "wagmi";
import { createAppKit } from "@reown/appkit/react";
import { WagmiAdapter } from "@reown/appkit-adapter-wagmi";
import { chain } from "./config";

const queryClient = new QueryClient();

const projectId =
  (import.meta.env.VITE_APPKIT_PROJECT_ID as string | undefined) ?? "43eff4b1fac476ffee43e467ab916f34";

const appUrl = typeof window !== "undefined" ? window.location.origin : "https://symbiotic.fi";

const metadata = {
  name: "Symbiotic Flight Delays",
  description: "Purchase policies, manage vaults, and claim rewards.",
  url: appUrl,
  icons: ["https://symbiotic.fi/favicon.ico"]
};

const networks = [chain];

export const wagmiAdapter = new WagmiAdapter({
  networks: networks as any,
  projectId,
  ssr: false
});

createAppKit({
  adapters: [wagmiAdapter],
  networks: networks as any,
  projectId,
  metadata,
  features: {
    analytics: false
  },
  themeMode: "light"
});

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <WagmiProvider config={wagmiAdapter.wagmiConfig}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </WagmiProvider>
  );
}
