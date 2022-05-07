const pg = require('pg')
const format = require('pg-format')

class Storage {
  constructor(client) {
    this.client = client
  }

  async connect() {
    await this.client.connect()
    return this
  }

  async close() {
    await this.client.end()
  }

  async getHoldersOfTokens(network, tokens) {
    let filter = tokens.map(t => `lower(token) = '${t}'`).join(' OR ')
    let query = `SELECT DISTINCT holder
                 FROM token_holders_last
                 WHERE ${filter}`
    let res = await this.client.query({ text: query, rowMode: 'array' })
    return res.rows.flat()
  }

  async dumpActivity(holdersActivityList) {
    const query = ([holder, day, week, month]) => `
      INSERT INTO holders_activity
      VALUES ('${holder}', ${day}, ${week}, ${month})
      ON CONFLICT (holder) DO UPDATE
      SET day = ${day}, week = ${week}, month = ${month};
    `;

    return await this.client.query(holdersActivityList.map(query).join(''))
  }

  async dumpTransactions(network, list) {
    if (list.length <= 0) return

    const query = 'INSERT INTO txs (tx_hash, block, from_address, to_address, value, successful) VALUES %L ON CONFLICT DO NOTHING'
    const values = list.map(t => [t.hash, t.block.height, t.from.address, t.to.address, t.value, t.successful])
    const res = await this.client.query(format(query, values))
    return res.rowCount
  }
}

exports.Storage = Storage
