import type { AirlineDTO, AirlineWithFlights, FlightDTO } from "../types";

async function getJSON<T>(url: string): Promise<T> {
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(`Request failed: ${resp.status}`);
  }
  return resp.json();
}

export async function fetchAirlinesWithFlights(baseUrl: string): Promise<AirlineWithFlights[]> {
  const airlinesResp = await getJSON<{ airlines: AirlineDTO[] }>(`${baseUrl}/airlines`);
  const airlines = airlinesResp.airlines ?? [];

  const results = await Promise.all(
    airlines.map(async (airline) => {
      const flightsResp = await getJSON<{ flights: FlightDTO[] }>(
        `${baseUrl}/airlines/${encodeURIComponent(airline.airlineId)}/flights`
      );
      return { ...airline, flights: flightsResp.flights ?? [] };
    })
  );
  return results;
}
