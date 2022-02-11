CREATE TABLE analytics.opensea_trades (
  network UInt32,
  block UInt64,
  time UInt64,
  tx String,
  buy_hash String,
  sell_hash String,
  maker String,
  taker String,
  price Int256
) ENGINE = MergeTree()
ORDER BY (network, time);