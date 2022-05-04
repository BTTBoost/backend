require('dotenv').config()
const { fetchHolders } = require('../src/api.js')
const { readNFTTokenlistWithoutHolders, saveHolders } = require('../src/db.js')

const network = 1
const parallelLimit = 50

// fetch holders for all each nft token without holders
async function main() {
  var i = 0
  do {
    var tokenlist = await readNFTTokenlistWithoutHolders()
    console.log(`Found ${tokenlist.length} tokens without holders...`)

    var requests = []
    for (const t of tokenlist.slice(0, parallelLimit)) {
      const r = fetchHolders(network, t.address)
        .catch(e => {
          console.log(`[${i++}|${t.address}] failed to fetch holders: ${e}`)
          throw e
        })
        .then(holders => saveHolders(network, t.address, holders))
        .then(count => console.log(`[${i++}|${t.address}] ${count} holders saved`))
        .catch(e => console.log(`[${i++}|${t.address}] failed to update holders: ${e}`))
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