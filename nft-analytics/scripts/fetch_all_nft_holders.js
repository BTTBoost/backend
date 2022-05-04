require('dotenv').config()
const fastq = require('fastq')
const { fetchHolders } = require('../src/api.js')
const { readNFTTokenlistWithoutHolders, saveHolders } = require('../src/db.js')

const network = 1
const concurrency = 25

// fetch holders for all each nft token without holders
async function main() {
  var i = 0
  const worker = (address) => fetchHolders(network, address)
    .catch(e => {
      console.log(`[${i++}|${address}] failed to fetch holders: ${e}`)
      throw e
    })
    .then(holders => saveHolders(network, address, holders))
    .then(count => { console.log(`[${i}|${address}] ${count} holders saved`); i++ })
    .catch(e => console.log(`[${i}|${address}] failed to update holders: ${e}`))

  do {
    var tokenlist = await readNFTTokenlistWithoutHolders()
    console.log(`Found ${tokenlist.length} tokens without holders...`)

    const queue = fastq.promise(worker, concurrency)
    for (const t of tokenlist) {
      queue.push(t.address)
    }

    await queue.drain()
  } while (tokenlist.length > 0)
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })