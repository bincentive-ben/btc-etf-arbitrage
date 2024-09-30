# btc-etf-arbitrage
Application to find arbitrage opportunities between BTC and ETFs

# Paper account information
    Username: vnaphx534
    Account Number: DUA774243

# Environment setup
* https://interactivebrokers.github.io/cpwebapi/quickstart

## OAuth (not supported)
### Generation with openssl
1. openssl dhparam  -out dhparam.pem 2048
2. openssl genrsa -out private_signature.pem 2048
3. openssl genrsa -out private_encryption.pem 2048
4. openssl rsa -in private_signature.pem -outform PEM -pubout -out public_signature.pem
5. openssl rsa -in private_encryption.pem -outform PEM -pubout -out public_encryption.pem



