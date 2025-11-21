import type { AirlineDTO, AirlineWithFlights, FlightDTO } from "../types";

async function getJSON<T>(url: string): Promise<T> {
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(`Request failed: ${resp.status}`);
  }
  return resp.json();
}

function normalizeBase(url: string) {
  return url.replace(/\/+$/, "");
}

export async function fetchAirlinesWithFlights(baseUrl: string): Promise<AirlineWithFlights[]> {
  const normalized = normalizeBase(baseUrl);
  const airlinesResp = await getJSON<{ airlines: AirlineDTO[] }>(`${normalized}/airlines`);
  const airlines = airlinesResp.airlines ?? [];

  const results = await Promise.all(
    airlines.map(async (airline) => {
      const flightsResp = await getJSON<{ flights: FlightDTO[] }>(
        `${normalized}/airlines/${encodeURIComponent(airline.airlineId)}/flights`,
      );
      return { ...airline, flights: flightsResp.flights ?? [] };
    }),
  );
  return results;
}
