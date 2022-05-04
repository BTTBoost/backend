require('dotenv').config()
const { fetchNFTMarket } = require('../src/api.js')
const { saveNFTTokenlist } = require('../src/db.js')

async function main() {
  return fetchNFTMarket(1)
    .then(tokenlist => {
      console.log(`Fetched ${tokenlist.length} tokens`)
      return saveNFTTokenlist(tokenlist)
    })
    .then(count => console.log(`Saved ${count} tokens to db`))
    .catch(error => console.error(`Failed to update tokenlist: ${error}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })