-- naive not optimized
SELECT T1.time, COUNT(*)
FROM (
	SELECT time, holder
	FROM analytics.token_holders
	WHERE token = '0x5a98fcbea516cf06857215779fd812ca3bef1b32' AND time >= 1633046400
) T1
INNER JOIN (
	SELECT time, holder
	FROM analytics.token_holders
	WHERE token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time >= 1633046400
) T2 ON (T1.holder = T2.holder AND T1.time = T2.time)
INNER JOIN (
	SELECT time, holder
	FROM analytics.token_holders
	WHERE token = '0xa5f1ea7df861952863df2e8d1312f7305dabf215' AND time >= 1633046400
) T3 ON (T1.holder = T3.holder AND T1.time = T3.time)
GROUP BY H0.time
ORDER BY H0.time;


-- GROUP((T1xT2)xT3)
-- 1) GROUP wrapper
-- 2) T builders
-- 3) T1xT2 builder
SELECT time, COUNT(*)
FROM (
  SELECT T12.time, T12.holder
  FROM (
    SELECT T1.time, T1.holder
    FROM (
      SELECT time, holder
      FROM analytics.token_holders
      WHERE network = 1 AND token = '0x5a98fcbea516cf06857215779fd812ca3bef1b32' AND time >= 1
    ) T1
    INNER JOIN (
      SELECT time, holder
      FROM analytics.token_holders
      WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time >= 1
    ) T2 ON (T1.holder = T2.holder AND T1.time = T2.time)
  ) T12
  INNER JOIN (
    SELECT time, holder
    FROM analytics.token_holders
    WHERE network = 1 AND token = '0xa5f1ea7df861952863df2e8d1312f7305dabf215' AND time >= 1
  ) T3 ON T12.time = T3.time AND T12.holder = T3.holder
)
GROUP BY time, holder
ORDER BY time;


-- generated
SELECT time, COUNT(*) 
FROM (
  SELECT T1.time, T1.holder 
  FROM (
    SELECT T1.time, T1.holder 
    FROM (
      SELECT time, holder 
      FROM analytics.token_holders 
      WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' AND time >= 1636156800 AND amount >= 1
    ) T1 
    INNER JOIN (
      SELECT time, holder 
      FROM analytics.token_holders 
      WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042' AND time >= 1636156800 AND amount >= 1
    ) T2 ON (T1.time = T2.time AND T1.holder = T2.holder)
  ) T1 
  INNER JOIN (
    SELECT time, holder 
    FROM analytics.token_holders 
    WHERE network = 1 AND token = '0x5a98fcbea516cf06857215779fd812ca3bef1b32' AND time >= 1636156800 AND amount >= 1
  ) T2 ON (T1.time = T2.time AND T1.holder = T2.holder)
) 
GROUP BY time 
ORDER BY time ASC