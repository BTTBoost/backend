-- insert token holdings to token_holdings table
INSERT INTO token_holdings(network, token, holding_token, holders, holders_share)
(
	SELECT 1, '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d', T2.token as token, COUNT(DISTINCT T2.address) as holders, 
	  COUNT(DISTINCT T2.address) * cast(100 as float) / (
	  	SELECT COUNT(*) FROM token_holders_last WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d'
	  ) as "holders_share"
	FROM (
	  SELECT holder FROM token_holders_last 
	  WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d'
	) T1 
	INNER JOIN (
	  SELECT *
	  FROM balances
	) T2 ON T1.holder = T2.address
	GROUP BY T2.network, T2.token
	ORDER BY COUNT(DISTINCT T2.address) DESC
);

-- query holding tokens
SELECT T2.network as network, T2.token as token, COUNT(DISTINCT T2.address) as holders, 
  COUNT(DISTINCT T2.address) * cast(100 as float) / (
  	SELECT COUNT(*) FROM token_holders_last WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042'
  ) as "holders_share"
FROM (
  SELECT holder FROM token_holders_last 
  WHERE network = 1 AND token = '0x026224a2940bfe258d0dbe947919b62fe321f042'
) T1 
INNER JOIN (
  SELECT *
  FROM balances
) T2 ON T1.holder = T2.address
GROUP BY T2.network, T2.token
ORDER BY COUNT(DISTINCT T2.address) DESC;

-- query holdings
SELECT holding_token, holders, holders_share 
FROM token_holdings 
WHERE network = 1 AND token = '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d' 
ORDER BY holders DESC;

-- clear token_holdings for token
DELETE FROM token_holdings WHERE token = '0x026224a2940bfe258d0dbe947919b62fe321f042';