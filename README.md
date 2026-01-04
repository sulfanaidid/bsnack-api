# BSNACK API

Backend API untuk sistem penjualan keripik pangsit dengan fitur:
- Manajemen transaksi
- Perhitungan poin customer
- Laporan penjualan berdasarkan periode

## Tech Stack
- Go (Gin)
- PostgreSQL

## How to Run

```bash
cp .env.example .env
go run cmd/server/main.go