require('dotenv').config()
const { fetchBalances } = require('../src/api.js')
const { saveBalances } = require('../src/db.js')

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