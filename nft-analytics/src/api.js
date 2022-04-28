import fetch from 'node-fetch'

export const fetchHolders = async function (network, token) {
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

export const fetchBalances = async function (network, address) {
  const pageSize = 100000
  const url = `https://api.covalenthq.com/v1/${network}/address/${address}/balances_v2/?` +
    `quote-currency=USD&format=JSON&nft=true&no-nft-fetch=false&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=${pageSize}`
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