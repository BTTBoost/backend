const fetch = require('node-fetch')
const { createAlchemyWeb3 } = require('@alch/alchemy-web3')

exports.fetchHolders = async function (network, token) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/tokens/${token}/token_holders/?` +
    `quote-currency=USD&format=JSON&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many holders')
      }
      return body.data.items
    })
}

exports.fetchBalances = async function (network, address) {
  const pageSize = 100000
  const withNFTs = false // API endpoint not working with NFTs
  const url = `https://api.covalenthq.com/v1/${network}/address/${address}/balances_v2/?` +
    `quote-currency=USD&format=JSON&nft=${withNFTs}&no-nft-fetch=false&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        // handle 'Endpoint will predictably time out ...'
        if (body.error_code == 406) {
          return []
        }
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many items')
      }
      return body.data.items
    })
}

// only on Ethereum mainnet now
exports.fetchTokenMetadata = async function (token) {
  const web3 = createAlchemyWeb3(
    `https://eth-mainnet.g.alchemy.com/v2/${process.env.ALCHEMY_API_KEY}`,
  )
  return await web3.alchemy.getTokenMetadata(token)
}

exports.fetchTxs = async function (network, address) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/address/${address}/transactions_v2/?` +
    `key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many items')
      }
      return body.data.items
    })
}

exports.fetchNFTMarket = async function (network) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/nft_market/?` +
    `key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      if (body.data.items.length >= pageSize) {
        throw new Error('too many items')
      }
      return body.data.items
    })
}