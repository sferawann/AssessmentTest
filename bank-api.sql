CREATE DATABASE IF NOT EXISTS bank-api;

CREATE TYPE jenis_transaksi AS ENUM ('tabung', 'tarik');

CREATE TABLE IF NOT EXISTS nasabah (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255),
    nik VARCHAR(50) UNIQUE,
    no_hp VARCHAR(20) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel rekening
CREATE TABLE IF NOT EXISTS rekening (
    id SERIAL PRIMARY KEY,
    nasabah_id INTEGER NOT NULL,
    no_rekening VARCHAR(50) UNIQUE,
    saldo DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (nasabah_id) REFERENCES nasabah(id)
);

-- Membuat tabel transaksi
CREATE TABLE IF NOT EXISTS transaksi (
    id SERIAL PRIMARY KEY,
    rekening_id INTEGER NOT NULL,
    nominal DECIMAL(15, 2),
    jenis_transaksi jenis_transaksi NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rekening_id) REFERENCES rekening(id)
);
