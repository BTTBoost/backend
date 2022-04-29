import 'dotenv/config'
import { fetchBalances } from '../src/api.js'
import { readHoldersWithoutBalance, saveBalances } from '../src/db.js'

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token')
  const network = parseInt(process.argv[2])
  const token = process.argv[3]

  var i = 0
  do {
    var holders = await readHoldersWithoutBalance(network, token)
    console.log(`Found ${holders.length} holders without balances...`)

    var requests = []
    for (const h of holders.slice(0, 50)) {
      const r = fetchBalances(network, h.holder)
        .catch(e => {
          console.log(`[${i++}|${h.holder}] failed to fetch balances: ${e}`)
          throw e
        })
        .then(balances => saveBalances(network, h.holder, balances))
        .then(count => console.log(`[${i++}|${h.holder}] updated ${count} balances`))
        .catch(e => console.log(`[${i++}|${h.holder}] failed to update balances: ${e}`))
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