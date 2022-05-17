const got = require('got')
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

exports.ZapperClient = ZapperClient
