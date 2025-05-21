# Billing MCP Server

A Model Context Protocol (MCP) server designed to manage billing operations.

## Overview

This project provides a backend service that exposes billing functionalities through the Model Context Protocol. It allows clients to interact with billing data, such as creating, retrieving, and managing invoices. The server is built using Go, following principles of Hexagonal Architecture and Domain-Driven Design, with GORM as the ORM for database interactions and UUIDs for primary entity identification.

## Features

- Retrieve invoice details by its UUID.
- Retrieve a list of invoices based on various criteria (e.g., status, issue date range).
- (Internally) Invoices are composed of line items aggregated from various movement sources.

## Getting Started

To get started with the Billing MCP Server, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/ricardogrande-masmovil/billing-mcp.git
   ```

2. Navigate to the project directory:
   ```bash
   cd billing-mcp
   ```

3. Install the required dependencies:
   ```bash
   go mod tidy
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## Configuration

The server is configured via a `.config.yaml` file in the root of the project. A sample configuration file named `.config.example.yaml` is provided. You can copy this file to `.config.yaml` and modify it to suit your environment.

### Running Seed Data

To populate the database with initial seed data for development or testing, you can enable the `runSeeds` option in your `.config.yaml` file:

```yaml
# .config.yaml
server:
  host: "localhost"
  port: "8080"
database:
  host: "localhost"
  port: 5432
  user: "billing_user"
  password: "yoursecurepassword"
  dbname: "billing_db"
  sslmode: "disable"
  maxRetries: 3
logLevel: "info"
version: "0.0.1"
runSeeds: true # Set to true to run seed data on startup
```

By default, `runSeeds` is `false`.

## Setup an MCP client

To set up an MCP client, you will need the following config:
```json
{
    "servers": {
        "billing": {
            "type": "sse",
            "url": "http://localhost:8080/sse"
        }
    }
}
```
You can change the URL to point to your server's address if it's not running locally. This configuration works directly in VSCode, enabling you to test the server's functionality by using the agent chat mode.