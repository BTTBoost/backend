require('dotenv').config()

const pg = require('pg')
const { using } = require('using-statement')
const { Storage } = require('../src/storage')
const { ZapperClient } = require('../src/zapper')

async function main() {
  let zapper = new ZapperClient(process.env.ZAPPER_API_KEY)

  let storage = new Storage(new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  }))

  await using(await storage.connect(), async () => {
    const apps = await zapper.getApps()
    console.log(apps)
    await storage.dumpZapperApps(apps)
  })
}

main().catch(error => {
  console.log(error)
  process.exit(1)
})
