DROP TABLE IF EXISTS zapper_transactions;
CREATE TABLE IF NOT EXISTS zapper_transactions
(
    network       TEXT   NOT NULL,
    hash          TEXT   NOT NULL,
    block_number  BIGINT NOT NULL,
    name          TEXT,
    direction     TEXT,
    time_stamp    TEXT,
    symbol        TEXT,
    address       TEXT,
    amount        TEXT,
    "from"        TEXT,
    destination   TEXT,
    contract      TEXT,
    nonce         TEXT,
    gas_price     REAL,
    gas_limit     REAL,
    input         TEXT,
    gas           REAL,
    tx_successful BOOL,
    account       TEXT,
    PRIMARY KEY (network, hash)
);
