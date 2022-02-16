-- cross-token holder count
SELECT T2.network as network, T2.token as token, holders, 
  holders * 100 / (
    SELECT COUNT(*)
    FROM analytics.token_holders
    WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time = 1643932800
  ) as "%"
FROM (
  SELECT *
  FROM analytics.token_holders
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time = 1643932800
) T1
INNER JOIN (
  SELECT *
  FROM analytics.token_holders 
  WHERE time = 1643932800
) T2 ON T1.holder = T2.holder
GROUP BY T2.network, T2.token 
ORDER BY COUNT(*) as holders DESC
OFFSET 1;

-- cross-token amount
SELECT T2.network as network, T2.token as token, SUM(T2.amount) as total_amount
FROM (
  SELECT *
  FROM analytics.token_holders
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time = 1643932800
) T1
INNER JOIN (
  SELECT *
  FROM analytics.token_holders 
  WHERE time = 1643932800
) T2 ON T1.holder = T2.holder
GROUP BY T2.network, T2.token;


-- generated
SELECT T2.network as network, T2.token as token, holders,
  holders * 100 / (
    SELECT COUNT(*) FROM (
      SELECT T1.holder 
      FROM (
        SELECT holder 
        FROM analytics.token_holders 
        WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time = 1643932800 AND amount >= 1
      ) T1 
      INNER JOIN (
        SELECT holder 
        FROM analytics.token_holders 
        WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time = 1643932800 AND amount >= 1
      ) T2 ON (T1.holder = T2.holder))
  ) as "holders_share"
FROM (
  SELECT T1.holder 
  FROM (
    SELECT holder 
    FROM analytics.token_holders 
    WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time = 1643932800 AND amount >= 1
  ) T1 
  INNER JOIN (
    SELECT holder 
    FROM analytics.token_holders 
    WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time = 1643932800 AND amount >= 1
  ) T2 ON (T1.holder = T2.holder)
) T1 
INNER JOIN (
  SELECT * 
  FROM analytics.token_holders 
  WHERE time = 1643932800
) T2 ON T1.holder = T2.holder 
GROUP BY T2.network, T2.token 
ORDER BY COUNT(*) as holders DESC