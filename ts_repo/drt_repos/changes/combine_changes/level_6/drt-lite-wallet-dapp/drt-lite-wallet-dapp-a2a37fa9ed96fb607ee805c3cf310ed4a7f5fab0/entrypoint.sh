#!/bin/sh

# List of environment variables and their default values
env_vars_with_defaults="START_NETWORK_ID_STOP=sovereign START_NETWORK_NAME_STOP=Sovereign START_API_ADDRESS_STOP=https://api-sovereign-test.numbat.ro START_GATEWAY_URL_STOP= START_EXTRAS_API_URL_STOP=https://extras-api-sovereign-test.numbat.ro START_SAMPLE_AUTHENTICATED_DOMAINS_STOP=https://api-sovereign-test.numbat.ro START_SOVEREIGN_CONTRACT_ADDR_STOP=drt1qqqqqqqqqqqqqpgqkh3lkj9dznw7awmulw2xcfzkael83jaflrhsnv2qvy START_WALLET_ADDRESS_STOP=https://wallet-sovereign-test.numbat.ro START_WREWA_ID_STOP=WREWA-bd4d79 START2_NETWORK_ID_STOP2=Testnet START2_NETWORK_NAME_STOP2=Testnet START2_API_ADDRESS_STOP2=https://testnet-api.dharitri.org START2_GATEWAY_URL_STOP2= START2_EXTRAS_API_URL_STOP2=https://testnet-extras-api.dharitri.org START2_SAMPLE_AUTHENTICATED_DOMAINS_STOP2=https://testnet-api.dharitri.org START2_SOVEREIGN_CONTRACT_ADDR_STOP2=drt1qqqqqqqqqqqqqpgqkhqeu7e2t62pwuadcshlrmxcharcstkhlrhs6ylhvm START2_WALLET_ADDRESS_STOP2=https://testnet-wallet.dharitri.org START2_WREWA_ID_STOP2="

replace_placeholder() {
  local var_name=$1
  local var_value=$2

  echo "Var ${var_name} defined, replacing ${var_value} in config"
  find /usr/share/nginx/html/ -type f -exec sed -i 's|'${var_name}'|'${var_value}'|g' {} +
}

# Loop through each environment variable
for entry in $env_vars_with_defaults; do
  # Split the entry into name and value
  var_name=$(echo $entry | cut -d= -f1)
  default_value=$(echo $entry | cut -d= -f2)

  # Use the environment variable value if defined; otherwise, use the default
  eval "value=\${$var_name:-$default_value}"

  # Execute the function with the variable name and value
  replace_placeholder "$var_name" "$value"
done

exec nginx -g 'daemon off;'
