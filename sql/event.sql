SELECT T.time, D.time, T.holder, D.amount, D.tx
FROM (
  SELECT time, holder
  FROM analytics.token_holders
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time >= 1636156800
) T
INNER JOIN (
  SELECT time_day, time, user, amount, tx
  FROM analytics.aave_deposits
  WHERE time >= 1636156800
  ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day
) D ON T.time = D.time_day AND T.holder = D.user
ORDER BY T.time, D.time;

-- only non-0 time
SELECT T.time, SUM(D.amount), COUNT(*)
FROM (
  SELECT time, holder
  FROM analytics.token_holders
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time >= 1636156800
) T
INNER JOIN (
  SELECT time_day, time, user, amount, tx
  FROM analytics.aave_deposits
  WHERE time >= 1636156800
  ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day
) D ON T.time = D.time_day AND T.holder = D.user
GROUP BY T.time
ORDER BY T.time;

-- filled with 0 time
SELECT T.time, SUM(D.amount), COUNT(*)
FROM (
  SELECT time, holder
  FROM analytics.token_holders
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time >= 1636156800
) T
INNER JOIN (
  SELECT time_day, time, user, amount, tx
  FROM analytics.aave_deposits
  WHERE time >= 1636156800
  ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day
) D ON T.time = D.time_day AND T.holder = D.user
GROUP BY T.time
ORDER BY T.time WITH FILL FROM 1636156800 TO 1636156800+86400*90 STEP 86400;

-- generated
SELECT H.time, COUNT(*) 
FROM (
  SELECT time, holder 
  FROM analytics.token_holders 
  WHERE network = 1 AND token = '0x5a98fcbea516cf06857215779fd812ca3bef1b32' AND time >= 1636156800 AND amount >= 1
) H INNER JOIN (
  SELECT time_day as time, user, amount 
  FROM analytics.aave_deposits 
  WHERE time >= 1636156800 AND time < 1643932800 
  ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day
) E ON H.time = E.time AND H.holder = E.user 
GROUP BY H.time 
ORDER BY H.time WITH FILL FROM 1636156800 TO 1643932800 STEP 86400;

SELECT H.time, COUNT(*) 
FROM (
  SELECT time, holder 
  FROM analytics.token_holders 
  WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time >= 1636156800 AND time < 1643932800 AND amount >= 1
) H 
INNER JOIN (
  SELECT toUnixTimestamp(toStartOfDay(toDate(time))) as time, user 
  FROM (
    SELECT time, taker as user 
    FROM analytics.opensea_trades 
    WHERE network = 1 AND time >= 1636156800 AND time < 1643932800 
    UNION ALL 
    SELECT time, maker as user 
    FROM analytics.opensea_trades 
    WHERE network = 1 AND time >= 1636156800 AND time < 1643932800
  )
) E ON H.time = E.time AND H.holder = E.user 
GROUP BY H.time 
ORDER BY H.time WITH FILL FROM 1636156800 TO 1643932800 STEP 86400;