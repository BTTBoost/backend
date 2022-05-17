require('dotenv').config()

const pg = require('pg')
const got = require('got')
const { using } = require('using-statement')
const { Storage } = require('../src/storage')
const { promise: fastq } = require('fastq')

class ZapperClient {
  constructor(apiKey) {
    let authorization = `Basic ${Buffer.from(apiKey, 'utf8').toString('base64')}`
    this.client = got.extend({
      https: { rejectUnauthorized: false },
      prefixUrl: 'https://api.zapper.fi/v2/',
      headers: { authorization },
    })
  }

  async getTransactionsFor(address, network = undefined) {
    let response = await this.client.get({
      url: 'transactions',
      searchParams: {
        address,
        'addresses[]': address,
        network,
      },
    }).json()
    return response?.data
  }
}

async function main() {
  let nfts = [
    '0x026224a2940bfe258d0dbe947919b62fe321f042',
    '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d',
    '0x1a92f7381b9f03921564a437210bb9396471050c',
    '0x938e5ed128458139a9c3306ace87c60bcba9c067',
    '0x23581767a106ae21c074b2276D25e5C3e136a68b',
  ]

  let zapper = new ZapperClient(process.env.ZAPPER_API_KEY)

  let storage = new Storage(new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  }))

  let createRetryableWorker = (timeout, totalAttempts, asyncFunc) => {
    let worker = async (arg, attempts = 0) => {
      if (attempts > totalAttempts) return
      try {
        let result = await asyncFunc(arg)
        console.log(arg)
        return result
      } catch (err) {
        console.error('retry')
        console.dir(err, { depth: null });
        // sleep 5 min
        await new Promise((res) => setTimeout(res, timeout))
        return await worker(arg, attempts + 1)
      }
    }
    return worker
  }

  await using(await storage.connect(), async () => {
    let holders = await storage.getHoldersOfTokens(1, nfts)
    console.log({ holders: holders.length })

    let worker = createRetryableWorker(5 * 60 * 1000, 40,
      address => zapper.getTransactionsFor(address))

    let cnt = 0, now = Date.now()
    for (let holder of holders) {
      cnt += 1
      let list = await worker(holder)
      console.log(`txs: ${list?.length}, progress: ${cnt} / ${holders.length}, timeLeft: ${((Date.now() - now) / cnt * (holders.length - cnt) / 1000 / 60) | 0} minutes`)
      await storage.dumpZapperTransactions(list)
    }
  })
}

main().catch(error => {
  console.log(error)
  process.exit(1)
})
