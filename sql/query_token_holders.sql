SELECT time, COUNT(*)
FROM analytics.token_holders
WHERE network = 1
  AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042'
  AND time >= 1636156800
GROUP BY time
ORDER BY time;