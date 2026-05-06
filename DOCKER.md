# Docker & Multi-Stage Build Analysis

## Dockerfile Übersicht

Die Datei startet eine App + Datenbank zusammen, verbindet sie miteinander und stellt sicher, dass alles in der richtigen Reihenfolge läuft.

Sie definiert zwei Services:
- db: Eine PostgreSQL-Datenbank (Version 16, leichtgewichtiges Alpine-Image)
    - Benutzer, Passwort und Datenbankname sind gesetzt
    - Port 5432 wird nach außen freigegeben
    - Daten werden persistent im Volume pgdata gespeichert
    - Healthcheck prüft, ob die DB bereit ist
- api: Deine Anwendung
    - Wird aus dem aktuellen Verzeichnis gebaut (build: .)
    - Läuft auf Port 8080
    - Bekommt DB-Zugangsdaten über Umgebungsvariablen
    - Startet erst, wenn die Datenbank gesund ist
- volumes:
    - pgdata sorgt dafür, dass Datenbankdaten beim Neustart erhalten bleiben

## CGO_ENABLED=0: Erklärung & Bedeutung

### Was ist CGO?

CGO = C Go - Mechanismus, um C/C++-Code in Go aufzurufen

```
CGO_ENABLED=0  →  Statische Binary (keine C-Abhängigkeiten)
CGO_ENABLED=1  →  Dynamische Binary (benötigt C-Abhängigkeiten, glibc)
```

### Warum ist CGO_ENABLED=0 wichtig?

- Statische Binaries:
   - Mit `CGO_ENABLED=0` wird eine statisch gelinkte Binary erstellt
   - Alle Abhängigkeiten sind in der Binary enthalten
   - Keine zusätzlichen Laufzeit-Libraries nötig

- Alpine Kompatibilität:
   - Dynamisch gelinkte Binaries brauchen glibc und laufen nicht auf Alpine
   - Mit `CGO_ENABLED=0` läuft die Binary auf jedem Alpine-Image

- Portabilität:
   - Binary funktioniert überall (lokal, Docker, Kubernetes, etc.)
   - Keine Dependency Hell


## Image-Größen: Multi-Stage vs. Single-Stage

- Single-stage Build
    - Enthält Build-Tools, Dependencies und Source Code
    - Führt zu einem größeren Image
    - Oft unnötiger Ballast im finalen Container  
    - Typische Größe: >500 MB

- Multi-stage Build
    - Trennt Build-Phase und Runtime-Phase
    - Nur das fertige Artefakt wird ins finale Image kopiert
    - Führt zu einem deutlich kleineren Image
    - Typische Größe: ~50–100 MB

- Vorteile kleinerer Images
    - Schnellere Builds und Deployments  
    - Weniger Speicherverbrauch  
    - Geringere Angriffsfläche (Sicherheit)

# Testing
## Create

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}'
```

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}'
{"id":1,"name":"Widget","price":9.99}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Gadget","price":19.99}'
{"id":2,"name":"Gadget","price":19.99}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Gizmo","price":29.99}'
{"id":3,"name":"Gizmo","price":29.99}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

## Read All
```bash
curl http://localhost:8080/products
```
```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl http://localhost:8080/products
[{"id":1,"name":"Widget","price":9.99},{"id":2,"name":"Gadget","price":19.99},{"id":3,"name":"Gizmo","price":29.99}]
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

## Read One (by ID)
```bash
curl http://localhost:8080/products/1
````

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl http://localhost:8080/products/1
{"id":1,"name":"Widget","price":9.99}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

## Read One - Error expected
```bash
curl http://localhost:8080/products/999
```
```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl http://localhost:8080/products/999
{"error":"Product not found"}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

## Update
```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget Pro","price":14.99}'
```

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget Pro","price":14.99}'
{"id":1,"name":"Widget Pro","price":14.99}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```
```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl http://localhost:8080/products
[{"id":1,"name":"Widget Pro","price":14.99},{"id":2,"name":"Gadget","price":19.99},{"id":3,"name":"Gizmo","price":29.99}]
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

## Delete
```bash
curl -X DELETE http://localhost:8080/products/3
```

```bash
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl -X DELETE http://localhost:8080/products/3
{"result":"success"}
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % curl http://localhost:8080/products            
[{"id":1,"name":"Widget Pro","price":14.99},{"id":2,"name":"Gadget","price":19.99}]
tobias@Tobiass-MacBook-Air-2 CI-CD-MCM-Vollmann % 
```

