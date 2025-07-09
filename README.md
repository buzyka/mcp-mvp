# MCP Proxy Service

This service ("MCP") is a proxy for routing requests to a Shopware instance with OAuth2 authentication. It exposes a single endpoint that forwards API calls to your Shopware platform.

## Features

- Reads configuration from `.env` (`.env.example` contains all adjustments)

- Obtains OAuth2 token using Shopware client credentials

- Proxies any HTTP method and path under `/mcp/:chat_id/*proxyPath`

- Streams request and response bodies, preserving headers

## Prerequisites

- Go 1.24+

- Git

- A running Shopware 6 instance with prepared integration credentials:

    - Client Access Key ID

    - Client Secret Access Key

    - API domain URL

- `curl` (for testing)

## Setup

1. **Clone the repository:**
```bash
    git clone https://github.com/buzyka/mcp-mvp.git
    cd mcp-mvp
```

2. **Copy `.env.example` to `.env` and update values:**
``` bash
    cp .env.example .env
```

3. **Edit `.env` and set your values**

   ```dotenv
   # .env
   MCP_ACCESS_TOKEN=<any static token, not used by default>
   SHOPWARE_ACCESS_KEY_ID=<your client access key ID>
   SHOPWARE_SECRET_ACCESS_KEY=<your client secret>
   SHOPWARE_PLATFORM_DOMAIN=http://localhost:8000
   ```

---

## Configuration (.env)

| Key                          | Description                                          |
| ---------------------------- | ---------------------------------------------------- |
| `SHOPWARE_PLATFORM_DOMAIN`   | Base URL of your Shopware instance (no trailing `/`) |
| `SHOPWARE_ACCESS_KEY_ID`     | OAuth2 Client Access Key ID                          |
| `SHOPWARE_SECRET_ACCESS_KEY` | OAuth2 Client Secret Access Key                      |
| `MCP_ACCESS_TOKEN`           | authorisation token for MCP                          |

---

## Running

Run the service locally:

```bash
# Load env and start
go run main.go
```

By default, the server listens on port **8080**.

---

## Usage Examples

Assume your Shopware API is accessible at `http://localhost:8000`.

### 1. Get Shopware `_info version`

```bash
curl -v http://localhost:8080/mcp/12345/api/_info/version
```

* **12345** is an arbitrary chat ID for routing context.
* The proxy will fetch `/api/_info/version` from Shopware, adding `Authorization: Bearer <token>` `token` should be equal to the `MCP_ACCESS_TOKEN` .

