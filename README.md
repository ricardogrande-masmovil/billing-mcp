# Billing MCP Server

A Model Context Protocol (MCP) server designed to manage billing operations.

## Overview

This project provides a backend service that exposes billing functionalities through the Model Context Protocol. It allows clients to interact with billing data, such as creating, retrieving, and managing invoices.

## Features

- Retrieve invoice details

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

## Setup an MCP client
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