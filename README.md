# Trading Bot

## Features (WIP)

- [ ] Live and paper trading with [Alpaca](https://alpaca.markets/)
- [ ] Live crypto trading with [Binance](https://www.binance.com/en)

**Note:** Functionality requires setting up API keys with corresponding services.

## Getting Started

This repository is divided into the following sections:

- `adhoc`: Useful ad hoc scripts and experiments
- `dash`: UI for visualizing trades and account balances
- `jobs`: Periodically scheduled jobs
- `lib`: Shared Go libraries
- `quant`: Machine learning models and trading strategies
- `trader`: Backend server for UI and makes periodic trades
- `traderdb`: Database containing financial, ML and user info

This directory will be referred to by the `TRADING_BOT_REPO` environment variable.

## Technologies

- [Go 1.22](https://go.dev/)
- [Node LTS](https://nodejs.org/en/)
- [PostgreSQL 16](https://www.postgresql.org/)
- [Python 3.11](https://www.python.org/)
- [React 18](https://reactjs.org/)
- [Terraform](https://www.terraform.io/)
