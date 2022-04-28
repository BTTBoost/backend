import 'dotenv/config'
import { fetchHolders } from '../src/api.js'
import { saveHolders } from '../src/db.js'

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