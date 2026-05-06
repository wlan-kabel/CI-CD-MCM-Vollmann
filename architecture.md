# Architecture Documentation

## Übersicht

Die Product Catalog API ist eine REST-basierte Produktverwaltungs-Anwendung in Go. Sie ermöglicht es, Produkte (mit Name und Preis) zu erstellen, anzuzeigen, zu aktualisieren und zu löschen (CRUD-Operationen).

### Funktionsweise:
- Client sendet HTTP-Requests an die API (Port 8080)
- Handler verarbeitet Requests, validiert Daten und speichert sie
- Daten können wahlweise im RAM (MemoryStore) oder in PostgreSQL gespeichert werden
- Die App antwortet mit JSON

### Verwendung:
- Lokal: `go run ./cmd/api` (MemoryStore)
- Docker: `docker compose up` (PostgreSQL)

### REST Endpoints

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | /health | Health-Check / Status |
| GET | /products | Alle Produkte |
| POST | /products | Neues Produkt |
| GET | /products/{id} | Einzelnes Produkt |
| PUT | /products/{id} | Produkt aktualisieren |
| DELETE | /products/{id} | Produkt löschen |

## Request Flow

HTTP Request → Router → Handler → Store → Database

- HTTP Request: Client sendet HTTP Request (z.B. `POST /products`)
- Router (mux):
   - `mux.Router` (gorilla/mux) matched die URL zum entsprechenden Handler
   - Beispiel: `POST /products` → `CreateProduct()` Handler
   - Unterstützt URL-Parameter wie `/products/{id}`
- Handler (handler.go):
   - Empfängt den Request
   - Parst JSON Body und validiert Daten
   - Ruft entsprechende Store-Methode auf (Create, Update, Delete, etc.)
   - Konvertiert Response zu JSON
- Store (memory.go oder postgres.go):
   - **MemoryStore:** Speichert in Go-Map im RAM
   - **PostgresStore:** Speichert via SQL in PostgreSQL-Datenbank
   - Beide implementieren die gleiche Schnittstelle (GetAll, GetByID, Create, Update, Delete)
- Response: Store gibt Ergebnis an Handler zurück, Handler sendet JSON-Response zum Client

## MemoryStore vs. PostgresStore

| Aspekt | MemoryStore | PostgresStore |
|--------|------------|---------------|
| **Speicherung** | RAM (Go-Map) | PostgreSQL-DB |
| **Persistenz** | Nein | Ja (auf Disk) |
| **App-Restart** | Daten verloren | Daten bleiben |
| **Performance** | Schnell | Langsamer |
| **Skalierbarkeit** | Begrenzt | Unbegrenzt |
| **Use Case** | Dev & Tests | Production |


