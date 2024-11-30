# Simple Wallet App
A Simple Wallet Application to manager user wallet, including balance disbursement.

### Feature
- Disbursement: Enables users to transfer funds from their account

## Getting Started

### Prerequisite
- Go: Version 1.20 or higher
- Database: PostgreSQL

### Installation
- Clone the repository

```bash
git clone https://github.com/larif216/simple-wallet-app.git
cd simple-wallet-app
```

- Install dependency

```bash
go mod tidy
```

- Copy env.sample

```bash
cp env.sample .env
```

- Create database or just use docker to save your time.

```bash
cd dev
docker-compose up -d
```

#### Database Migration
- Download database migration tools

```bash
make tool-migrate
```

- Export the config, you can skip if you want to use default value (see: [Makefile](https://github.com/larif216/simple-wallet-app/blob/main/Makefile))

```bash
export POSTGRES_USER ?= postgres
export POSTGRES_PASS ?= postgres
export POSTGRES_DB_NAME ?= wallet
export POSTGRES_DB_HOST ?= 127.0.0.1
export POSTGRES_DB_PORT ?= 5432
```

- Run database migration

```bash
MIGRATE_ARGS=up make migrate
```

- (Optional) Populate database

```bash
make seed
```

#### Running in Local

- Start the application

```bash
go run cmd/main.go
```

## Usage

### API Endpoint

| Endpoint | Method    | Description    |
| :---:   | :---: | :---: |
| `/wallet/disburse` | `POST` | Processes a disbursement transaction |

### Example Request

#### Disbursement

```bash
curl -X POST http://localhost:8080/wallet/disburse \
-H "Content-Type: application/json" \
-d '{"user_id": 1, "amount": 100}'
```

### Example Response

```bash
{
  "disbursement_id": 1,
  "disbursement_status": "SUCCESS",
  "message": "Disbursement successfully processed"
}
```


