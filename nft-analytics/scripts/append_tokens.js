import 'dotenv/config'
import { readFile } from 'fs/promises'
import { appendTokens } from '../src/db.js'


async function main() {
  const nfts = JSON.parse(await readFile(new URL('../tokenlist_nfts.json', import.meta.url)));

  await appendTokens(nfts)
    .then(count => console.log(`Appended ${count} nfts to db`))
    .catch(e => console.error(`Failed to append tokens: ${e}`))
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })