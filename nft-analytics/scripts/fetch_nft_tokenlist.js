import 'dotenv/config'
import { fetchNFTMarket } from '../src/api.js'
import { saveNFTTokenlist } from '../src/db.js'

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