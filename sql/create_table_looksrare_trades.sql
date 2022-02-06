CREATE TABLE analytics.looksrare_trades (
  network UInt32,
  block UInt64,
  time UInt64,
  tx String,
  maker String,
  taker String,
  currency String,
  price Int256
) ENGINE = MergeTree()
ORDER BY (network, time);