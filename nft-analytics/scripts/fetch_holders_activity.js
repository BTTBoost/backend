require('dotenv').config()

const pg = require('pg')
const got = require('got')
const { using } = require('using-statement')
const { Storage } = require('../src/storage')
const { subDays, subWeeks, subMonths } = require('date-fns')
const { promise: fastq } = require('fastq')


async function fetchTransactionsInEth(txTo, page, since) {
  const limit = 1000
  const offset = page * limit
  const query = `
    query ($offset: Int!, $txTo: [String!]!, $since: ISO8601DateTime) {
      ethereum(network: ethereum) {
        transactions(
          options: {desc: "block.timestamp.unixtime", limit: ${limit}, offset: $offset}
          date: {since: $since}
          txTo: {in: $txTo}
        ) {
          block {
            height
            timestamp { unixtime }
          }
          from: sender { address }
          to { address }
          value: amount
          successful: success
          hash
        }
      }
    }
  `

  const response = got
    .post('https://graphql.bitquery.io/', {
      json: { query, variables: { txTo, since, offset } },
      headers: { 'X-API-KEY': process.env.BITQUERY_API_KEY },
      https: { rejectUnauthorized: false },
    })
    .json()

  return response
}

async function fetchActivity(holder, now) {
  const query = `
    query ($holder: String!, $dayAgo: ISO8601DateTime!, $weekAgo: ISO8601DateTime!, $monthAgo: ISO8601DateTime!, $till: ISO8601DateTime!) {
      activity: ethereum(network: ethereum) {
        day: transactions(
          options: {limit: 1}
          date: {since: $dayAgo, till: $till}
          txSender: {is: $holder}
        ) {
          _: success
        }
        
        week: transactions(
          options: {limit: 1}
          date: {since: $weekAgo, till: $till}
          txSender: {is: $holder}
        ) {
          _: success
        }
        
        month: transactions(
          options: {limit: 1}
          date: {since: $monthAgo, till: $till}
          txSender: {is: $holder}
        ) {
          _: success
        }
      }
    }
  `

  const dayAgo = subDays(now, 1).toISOString()
  const weekAgo = subWeeks(now, 1).toISOString()
  const monthAgo = subMonths(now, 1).toISOString()
  const till = now.toISOString()

  const response = await got
    .post('https://graphql.bitquery.io/', {
      json: { query, variables: { holder, dayAgo, weekAgo, monthAgo, till } },
      headers: { 'X-API-KEY': process.env.BITQUERY_API_KEY },
      https: { rejectUnauthorized: false },
    })
    .json()

  let activity = response?.data?.activity
  return {
    day: activity?.day?.length > 0,
    week: activity?.week?.length > 0,
    month: activity?.month?.length > 0,
  }
}


async function main() {
  // if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and address')
  // const network = parseInt(process.argv[2])
  // const address = process.argv[3]

  let nfts = [
    '0x026224a2940bfe258d0dbe947919b62fe321f042',
    '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d',
    '0x1a92f7381b9f03921564a437210bb9396471050c',
    '0x938e5ed128458139a9c3306ace87c60bcba9c067',
    '0x23581767a106ae21c074b2276D25e5C3e136a68b',
  ]

  let storage = new Storage(new pg.Client({
    connectionString: process.env.DB_CONN_STRING,
    ssl: { rejectUnauthorized: false },
  }))

  await using(await storage.connect(), async () => {
    let now = new Date()
    let acc = []

    let cnt = 0
    let fetchAndSum = async (address) => {
      try {
        let activity = await fetchActivity(address, now)
        acc.push([address, activity.day, activity.week, activity.month])
        if (++cnt % 1000 === 0) console.log(`Fetched ${cnt} holders`)
      } catch (err) {
        console.error('retry')
        // sleep 1 min
        await new Promise((res) => setTimeout(res, 60_000))
        await fetchAndSum(address)
      }
    }

    let holders = await storage.getHoldersOfTokens(1, nfts)

    // let queue = fastq(fetchAndSum, 6)
    // await Promise.all(holders.map(queue.push))
    for (let holder of holders) {
      await fetchAndSum(holder)
      if (cnt % 100 === 0) {
        await storage.dumpActivity(acc)
        acc = []
      }
    }

    console.log({ holders: { total: holders.length, fetched: cnt } })

    let dbResult = await storage.dumpActivity(acc)
    // console.dir({ dbResult }, { depth: null })
  })
}

main().catch(error => {
  console.log(error)
  process.exit(1)
})
