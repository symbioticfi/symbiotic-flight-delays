#!/bin/sh
FLIGHT_DELAYS_ADDRESS=0xA4b0f5eb09891c1538494c4989Eea0203b1153b1

exec /app/flight-node --relay-api-url "$1" --evm-rpc-url http://anvil:8545 --flight-delays-address "$FLIGHT_DELAYS_ADDRESS" --flights-api-url http://flights-api:8085 --private-key "$2" --log-level info
