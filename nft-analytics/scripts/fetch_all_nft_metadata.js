require('dotenv').config()
const { fetchTokenMetadata } = require('../src/api.js')
const { readNFTTokenlistWithoutMetadata, saveMetadata } = require('../src/db.js')


const parallelLimit = 50

// fetch token metadata for all tokens from nft tokenlist without metadata
async function main() {
  var i = 0
  do {
    var tokenlist = await readNFTTokenlistWithoutMetadata()
    console.log(`Found ${tokenlist.length} tokens without metadata...`)

    var requests = []
    for (const t of tokenlist.slice(0, parallelLimit)) {
      const r = fetchTokenMetadata(t.address)
        .catch(e => {
          console.log(`[${i++}|${t.address}] failed to fetch metadata: ${e}`)
          throw e
        })
        .then(metadata => saveMetadata(1, t.address, metadata))
        .then(_ => console.log(`[${i++}|${t.address}] metadata saved`))
        .catch(e => console.log(`[${i++}|${t.address}] failed to update metadata: ${e}`))
      requests.push(r)
    }
    await Promise.all(requests)
  } while (tokenlist.length > 0)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })