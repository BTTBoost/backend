import 'dotenv/config'
import fetch from 'node-fetch'
import pg from 'pg'
import format from 'pg-format'

async function fetchBalances(network, address) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/address/${address}/balances_v2/?` +
    `quote-currency=USD&format=JSON&nft=true&no-nft-fetch=false&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many items')
      }
      return body.data.items
    })
}

async function saveBalances(network, address, balances) {
  const client = new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  })
  await client.connect()

  var count = 0
  try {
    await client.query('BEGIN')

    await client.query(format('DELETE FROM balances WHERE network = %L AND address = %L', network, address))

    const query = 'INSERT INTO balances (network, address, token, amount, amount_usd) VALUES %L'
    const values = balances.map(b => [network, address, b.contract_address, b.balance, b.quote])
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
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and address')
  const network = parseInt(process.argv[2])
  const address = process.argv[3]

  return fetchBalances(network, address)
    .then(balances => {
      console.log(`Fetched ${balances.length} assets`)
      return saveBalances(network, address, balances)
    })
    .then(count => console.log(`Saved ${count} assets to db`))
    .catch(error => console.error(`Failed to update assets: ${error}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })