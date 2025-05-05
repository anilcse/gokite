# GoKite - Automated Trading Bot

GoKite is an automated trading bot built in Go that integrates with Zerodha's Kite Connect API. It allows you to define and execute trading strategies based on technical indicators and market conditions.

## Features

- Real-time market data processing via WebSocket
- Rule-based trading strategy execution
- Support for multiple timeframes and instruments
- PostgreSQL database for strategy persistence
- Hot-reloadable trading rules
- Configurable entry and exit parameters
- Support for various technical indicators (SMA, RSI)

## Prerequisites

- Go 1.20 or higher
- PostgreSQL database
- Zerodha Kite Connect API credentials

## Installation

1. Clone the repository:
```bash
git clone https://github.com/anilcse/gokite.git
cd gokite
```

2. Install dependencies:
```bash
go mod download
```

3. Create a configuration file:
```bash
cp configs/app.yaml.example configs/app.yaml
```

4. Update the configuration in `configs/app.yaml` with your:
   - Database connection string
   - Kite Connect API credentials
   - Trading instruments

## Project Structure

```
.
├── cmd/
│   └── server/         # Main application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── engine/        # Trading strategy engine
│   ├── kite/          # Kite Connect API client
│   ├── model/         # Data models
│   ├── scheduler/     # Job scheduling
│   └── store/         # Database operations
├── configs/           # Configuration files
└── scripts/          # Utility scripts
```

## Usage

1. Build the application:
```bash
make build
```

2. Run the server:
```bash
./bin/server
```

## Configuration

The application is configured through `configs/app.yaml`. Key configuration sections:

- Database connection settings
- Kite Connect API credentials
- Trading instruments to monitor
- Strategy rules and parameters

## Development

- Run tests:
```bash
make test
```

- Build the application:
```bash
make build
```

- Clean build artifacts:
```bash
make clean
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Disclaimer

This software is for educational purposes only. Use at your own risk. The authors are not responsible for any financial losses incurred through the use of this software.
