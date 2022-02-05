SELECT network, token, COUNT(DISTINCT time)
FROM analytics.token_holders
GROUP BY network, token;