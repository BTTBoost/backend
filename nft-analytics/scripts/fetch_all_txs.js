import 'dotenv/config'
import { fetchTxs } from '../src/api.js'
import { readHoldersTxsWithoutTxs, saveTxs } from '../src/db.js'

const parallelLimit = 50

// fetch transasctions for each token holder
async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token')
  const network = parseInt(process.argv[2])
  const token = process.argv[3]

  var i = 0
  do {
    var holders = await readHoldersTxsWithoutTxs(network, token)
    console.log(`Found ${holders.length} holders without txs...`)

    var requests = []
    for (const h of holders.slice(0, parallelLimit)) {
      const r = fetchTxs(network, h.holder)
        .catch(e => {
          console.log(`[${i++}|${h.holder}] failed to fetch txs: ${e}`)
          throw e
        })
        .then(txs => saveTxs(network, h.holder, txs))
        .then(count => console.log(`[${i++}|${h.holder}] saved ${count} txs`))
        .catch(e => console.log(`[${i++}|${h.holder}] failed to update txs: ${e}`))
      requests.push(r)
    }
    await Promise.all(requests)
  } while (holders.length > 0)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })