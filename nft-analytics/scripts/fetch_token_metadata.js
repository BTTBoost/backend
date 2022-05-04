require('dotenv').config()
const { fetchTokenMetadata } = require('../src/api.js')
const { saveMetadata } = require('../src/db.js')

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token_address')
  const network = parseInt(process.argv[2])
  if (network != 1) {
    throw new Error('only ethereum mainnet is supported (network_id = 1)')
  }
  const token = process.argv[3]

  return fetchTokenMetadata(token)
    .then(metadata => {
      console.log(`Fetched metadata: ${JSON.stringify(metadata)}`)
      return saveMetadata(network, token, metadata)
    })
    .then(_ => console.log(`Saved metadata`))
    .catch(error => console.error(`Failed to update metadata: ${error}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })