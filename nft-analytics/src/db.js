import pg from 'pg'
import format from 'pg-format'


export const readHoldersWithoutBalance = async function (network, token) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  const res = await client.query(format(
    'SELECT holder FROM token_holders_last WHERE network = %L AND token = %L AND holder NOT IN (SELECT address FROM balance_updated) ORDER BY holder',
    network, token,
  ))

  await client.end()
  return res.rows
}

export const saveHolders = async function (network, token, holders) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  var count = 0
  try {
    await client.query('BEGIN')

    await client.query(format('DELETE FROM token_holders_last WHERE network = %L AND token = %L', network, token))

    const query = 'INSERT INTO token_holders_last (network, token, holder, amount) VALUES %L'
    const values = holders.map(h => [network, token, h.address, h.balance])
    const res = await client.query(format(query, values))

    await client.query('COMMIT')
    count = res.rowCount
  } catch (e) {
    console.error("Save error:", e)
    await client.query('ROLLBACK')
  }

  await client.end()
  return count
}

export const saveBalances = async function (network, address, balances) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  var count = 0
  try {
    await client.query('BEGIN')

    await client.query(format('DELETE FROM balances WHERE network = %L AND address = %L', network, address))

    const updatedAt = Date.now()
    await client.query(format(
      'INSERT INTO balance_updated (network, address, updated_at) VALUES (%L) ON CONFLICT (network, address) DO UPDATE SET updated_at = %L',
      [network, address, updatedAt], updatedAt,
    ))

    if (balances.length > 0) {
      const query = 'INSERT INTO balances (network, address, token, amount, amount_usd) VALUES %L'
      const values = balances.map(b => [network, address, b.contract_address, b.balance, b.quote])
      const res = await client.query(format(query, values))
      count = res.rowCount
    }

    await client.query('COMMIT')
  } catch (e) {
    console.error("Save error:", e)
    await client.query('ROLLBACK')
    throw e
  } finally {
    await client.end()
  }

  return count
}
