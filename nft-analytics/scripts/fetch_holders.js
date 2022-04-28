import 'dotenv/config'
import fetch from 'node-fetch'
import pg from 'pg'
import format from 'pg-format'

async function fetchHolders(network, token) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/tokens/${token}/token_holders/?` +
    `quote-currency=USD&format=JSON&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many holders')
      }
      return body.data.items
    })
}

async function saveHolders(network, token, holders) {
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

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token_address')
  const network = parseInt(process.argv[2])
  const token = process.argv[3]

  return fetchHolders(network, token)
    .then(holders => {
      console.log(`Fetched ${holders.length} holders`)
      return saveHolders(network, token, holders)
    })
    .then(count => console.log(`Saved ${count} holders to db`))
    .catch(error => console.error(`Failed to update holders: ${error}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })