# Catalog

**A modern physical inventory-management platform connecting real-world assets with trusted digital records through NFC and QR.**

Catalog is a self-hosted inventory application for schools, music programs, clubs, labs, and organizations. It is especially comfortable with instruments and equipment, while its flexible locations, labels, custom fields, and borrower records keep it useful for any physical inventory.

Catalog is designed, substantially redesigned, extended, and maintained by **Shahin Madantschi — Student Developer**. The project preserves its AGPL license and required source history; it does not claim that every original line was written from scratch.

## Product overview

Catalog answers the practical questions around a real object: what it is, where it belongs, who has it, what condition it is in, whether it is missing, and what a staff member should do next. Secure physical tags make an item's record and permitted actions reachable at the item itself.

## Features

- Inventory records with asset numbers, manufacturer, model, serial number, condition, notes, photos, attachments, custom fields, purchase details, and warranties
- Nested locations and flexible labels for instruments, electronics, uniforms, accessories, tools, or any organization-defined category
- NFC and QR physical tags with secure public URLs, a guided iPhone/Android/manual programming flow, tap-to-verify setup, replacement, and revocation
- Fast staff actions from an authenticated NFC tap: checkout, return, move, report damage, add maintenance, mark missing/found, edit, and view history
- Separate staff accounts and borrower records, so students or members can receive inventory without administrative logins
- Checkout and return records with borrower, due date, location, condition, notes, and human-readable item history
- First-class missing/found state with last-known context and privacy-safe public tag behavior
- Owner, Administrator, Staff, and Viewer roles enforced by the API; staff invitations, account disabling, and explicit ownership transfer
- Maintenance logs, damage reporting, item activity, search/filtering, import/export, reports, and printable labels
- Responsive Nuxt interface, light/dark themes, and installable web-app support
- Self-hosted single-container deployment with persistent SQLite storage

## Screenshots

1. Dashboard
![Catalog dashboard](assets/screenshots/dashboard.png)
2. NFC Tag setup

## Physical NFC and QR setup

Catalog never treats a generated URL as proof that a physical tag was programmed.

1. Open an item and choose **Set up tag**.
2. Catalog creates a random, revocable public token and a QR fallback. Internal item IDs are not placed on the tag.
3. Choose iPhone, Android, or manual setup.
4. Write the displayed HTTPS URL to an NFC tag.
5. Tap or scan the physical tag. Catalog records verification only when that public token is actually opened.
6. The item changes to **Ready**. A tag can later be reprogrammed, replaced, revoked, copied, or printed as QR.

### iPhone

iOS does not provide a built-in general NFC-tag writer. Use a compatible application such as NFC Tools: add a URL record, paste the Catalog link, choose Write, and hold the iPhone near the tag. iPhone can read compatible NDEF URL tags from the system without opening Catalog first.

### Android

Use a compatible NFC-writing application and write the Catalog URL as an NDEF URL record. Catalog may offer direct Web NFC writing when the current Android browser supports it, but Web NFC is progressive enhancement and is never required.

### Tag and browser expectations

- Use writable NDEF-compatible NFC tags with enough capacity for the complete HTTPS URL.
- Public deployments should use HTTPS. Localhost is allowed for local development, but phones cannot reach a computer's `localhost`.
- QR always remains available when NFC is unavailable.
- Replacing or revoking a tag invalidates the old token. Revoked tokens do not disclose item data.
- Anonymous tag responses contain a deliberately small public data shape. Sensitive item, borrower, staff, attachment, financial, and location data stays server-side.

## Architecture

| Layer | Implementation |
| --- | --- |
| Frontend | Nuxt 3, Vue 3, TypeScript, Tailwind CSS, DaisyUI, Pinia |
| Backend | Go HTTP API using Chi and service/repository layers |
| Data | SQLite with ent ORM schemas and versioned Atlas migrations |
| Authentication | Stateful hashed bearer tokens, secure HTTP-only cookies, organization roles enforced server-side |
| Physical tags | Random public tokens, NDEF URL/QR output, verification timestamps, scan timestamps, replacement/revocation lifecycle |
| Production | Nuxt static output embedded into the Go binary, built and run as one Docker container |
| Deployment assets | Docker Compose plus existing rootless, Fly.io, and Kubernetes/Portainer-oriented project assets |

The frontend calls versioned `/api/v1` endpoints. Repository methods scope organization-owned records by group ID; public NFC resolution uses a separate minimal response rather than serializing an authenticated item object and hiding fields in Vue.

## Quick start with Docker

Docker Compose is the simplest supported path:

```bash
git clone https://github.com/sdmad737/catalog.git
cd catalog
docker compose up --build -d
```

Open [http://localhost:3100](http://localhost:3100), register the first account, and create your organization. The first account becomes the Owner. Application data is stored in the `catalog-data` volume mounted at `/data`.

Useful commands:

```bash
docker compose logs -f catalog
docker compose restart catalog
docker compose down
```

If Docker reports that `dockerDesktopLinuxEngine` cannot be found on Windows, start Docker Desktop and wait until its Linux engine reports **Running**, then retry.

## Local development

### Frontend

Requirements: Node.js 20 and pnpm 8.15.9.

```bash
cd frontend
corepack enable
corepack prepare pnpm@8.15.9 --activate
pnpm install --frozen-lockfile --shamefully-hoist
pnpm dev
```

The Nuxt development server proxies `/api` to the Go service on port `7745`.

### Backend

Requirements: Go 1.22 or newer and a C compiler for the SQLite test driver. Production Docker currently builds with Go 1.24.

```bash
cd backend
go test ./...
go run ./app/api
```

Versioned migrations are embedded into the API and applied on startup. Do not edit a migration that has already shipped; create a new migration for later schema changes.

## Backup and restore

SQLite can use WAL sidecar files while Catalog is running, so stop the service before copying `/data`.

Backup to a local `catalog-data-backup` directory:

```bash
docker compose stop catalog
docker compose cp catalog:/data ./catalog-data-backup
docker compose start catalog
```

Restore into a newly created, stopped container:

```bash
docker compose down
docker compose create catalog
docker compose cp ./catalog-data-backup/. catalog:/data
docker compose up -d
```

Restoring replaces the active application data. Keep an additional copy of the current `/data` directory before a restore, and protect backups because they contain account and inventory data.

## Verification commands

```bash
# Frontend
cd frontend
pnpm lint:ci
pnpm typecheck
pnpm test:ci
pnpm build

# Backend
cd ../backend
gofmt -w ./app ./internal ./pkgs
go test ./...
go build ./app/api

# Complete image
cd ..
docker compose build
docker compose up -d
```

## Feature status

**Stable:** core inventory, locations, labels, custom fields, attachments, import/export, maintenance records, organization accounts, staff permissions, borrower directory, checkout/return, missing/found, item history, NFC/QR tag lifecycle, Docker deployment.

**Progressive/experimental:** direct browser-based NFC writing on compatible Android devices. Manual NFC-app programming and QR fallback remain the supported cross-platform flow.

**Planned:** verified release screenshots and broader end-to-end coverage across physical iPhone/Android devices.

## Security

Do not post vulnerability details in a public issue. Follow [SECURITY.md](SECURITY.md) for private reporting options. Public NFC endpoints are intentionally minimal, and automated tests assert that private item fields are not serialized.

## Project structure

```text
backend/               Go API, services, ent schemas, repositories, migrations, and tests
frontend/              Nuxt/Vue application and frontend tests
docs/                  Project documentation
docker-compose.yml     Local/self-hosted deployment
Dockerfile             Production frontend + API image
Dockerfile.rootless    Non-root container variant
```

## License and attribution

Catalog is licensed under the [GNU Affero General Public License v3](LICENSE). Preserve the license, notices, and legally required source history when distributing or hosting modified versions.
