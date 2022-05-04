const pg = require('pg')
const format = require('pg-format')

exports.readHoldersWithoutBalance = async function (network, token) {
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

exports.readHoldersTxsWithoutTxs = async function (network, token) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  const res = await client.query(format(
    'SELECT holder FROM token_holders_last WHERE network = %L AND token = %L AND holder NOT IN (SELECT address FROM txs_updated) ORDER BY holder',
    network, token,
  ))

  await client.end()
  return res.rows
}

// TODO: rename to replace
exports.saveHolders = async function (network, token, holders) {
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

// TODO: rename to replace
exports.saveBalances = async function (network, address, balances) {
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
      const balanceQuery = 'INSERT INTO balances (network, address, token, amount, amount_usd) VALUES %L'
      const balanceValues = balances.map(b => [network, address, b.contract_address, b.balance, b.quote])
      const res = await client.query(format(balanceQuery, balanceValues))
      count = res.rowCount

      const tokenQuery = 'INSERT INTO tokens (network, address, name, symbol, logo, decimals) VALUES %L ON CONFLICT (network, address) DO NOTHING'
      const tokenValues = balances.map(b => [
        network,
        b.contract_address,
        b.contract_name ? b.contract_name.replaceAll(String.fromCharCode(0), '') : null,
        b.contract_ticker_symbol ? b.contract_ticker_symbol.replaceAll(String.fromCharCode(0), '') : null,
        b.logo_url ? b.logo_url.replaceAll(String.fromCharCode(0), '') : null,
        b.contract_decimals
      ])
      await client.query(format(tokenQuery, tokenValues))
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

exports.appendTokens = async function (tokens) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  const query = 'INSERT INTO tokens (network, address, name, symbol, logo) VALUES %L ON CONFLICT (network, address) DO NOTHING'
  const values = tokens.map(t => [t.network, t.address, t.name, t.symbol, t.logo])
  const resp = await client.query(format(query, values))

  await client.end()

  return resp.rowCount
}

exports.saveMetadata = async function (network, token, metadata) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  try {
    await client.query('BEGIN')
    await client.query(format('DELETE FROM tokens WHERE network = %L AND address = %L', network, token))

    const query = 'INSERT INTO tokens (network, address, name, symbol, logo, decimals) VALUES (%L)'
    const values = [network, token, metadata.name, metadata.symbol, metadata.logo, metadata.decimals ? metadata.decimals : 0]
    await client.query(format(query, values))

    await client.query('COMMIT')
  } catch (e) {
    console.error("Save error:", e)
    await client.query('ROLLBACK')
  }

  await client.end()
}

exports.saveTxs = async function (network, address, txs) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  var count = 0
  try {
    await client.query('BEGIN')

    const updatedAt = Date.now()
    await client.query(format(
      'INSERT INTO txs_updated (network, address, updated_at) VALUES (%L) ON CONFLICT (network, address) DO UPDATE SET updated_at = %L',
      [network, address, updatedAt], updatedAt,
    ))

    if (txs.length > 0) {
      const query = 'INSERT INTO txs (tx_hash, block, from_address, to_address, value, successful) VALUES %L ON CONFLICT DO NOTHING'
      const values = txs.map(tx => [tx.tx_hash, tx.block_height, tx.from_address, tx.to_address, tx.value, tx.successful])
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

exports.saveNFTTokenlist = async function (tokenlist) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  var count = 0
  try {
    await client.query('BEGIN')

    await client.query('DELETE FROM nft_tokenlist')

    const query = 'INSERT INTO nft_tokenlist (address) VALUES %L'
    const values = tokenlist.map(t => [t.collection_address])
    const res = await client.query(format(query, values))
    count = res.rowCount

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

exports.readNFTTokenlistWithoutMetadata = async function (network, token) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  const res = await client.query('SELECT address FROM nft_tokenlist WHERE address NOT IN (SELECT address FROM tokens)')

  await client.end()
  return res.rows
}

exports.readNFTTokenlistWithoutHolders = async function (network, token) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  const res = await client.query('SELECT address FROM nft_tokenlist WHERE address NOT IN (SELECT DISTINCT token FROM token_holders_last)')

  await client.end()
  return res.rows
}