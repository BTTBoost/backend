import 'dotenv/config'
import { fetchTxs } from '../src/api.js'
import { saveTxs } from '../src/db.js'

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and address')
  const network = parseInt(process.argv[2])
  const address = process.argv[3]

  return fetchTxs(network, address)
    .then(txs => {
      console.log(`Fetched ${txs.length} transactions`)
      return saveTxs(network, address, txs)
    })
    .then(count => console.log(`Saved ${count} txs to db`))
    .catch(error => console.error(`Failed to update txs: ${error}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })