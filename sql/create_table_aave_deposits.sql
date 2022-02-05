CREATE TABLE analytics.aave_deposits (
  network UInt32,
  block UInt64,
  time UInt64,
  tx String,
  user String,
  on_behalf_of String,
  reserve String,
  amount Int256
) ENGINE = MergeTree()
ORDER BY (network, time);