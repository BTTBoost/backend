CREATE TABLE analytics.token_holders (
  network UInt32,
  time UInt64,
  token String,
  holder String,
  amount Int256
) ENGINE = MergeTree()
ORDER BY (network, time, token);