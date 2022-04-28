import 'dotenv/config'
import { fetchBalances } from '../src/api.js'
import { readHoldersWithoutBalance, saveBalances } from '../src/db.js'

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token')
  const network = parseInt(process.argv[2])
  const token = process.argv[3]

  try {
    const holders = await readHoldersWithoutBalance(network, token)
    console.log(`Found ${holders.length} holders without balances...`)

    var i = 0
    for (const h of holders) {
      try {
        const balances = await fetchBalances(network, h.holder)
        await saveBalances(network, h.holder, balances)
        console.log(`[${i++}|${h.holder}] update ${balances.length} balances`)
      } catch (e) {
        console.log(`[${i++}|${h.holder}] failed to update balances`)
      }
    }
  } catch (e) {
    throw e
  }
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })