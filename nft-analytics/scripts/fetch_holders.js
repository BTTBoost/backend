import 'dotenv/config'
import fetch from 'node-fetch'

// https://api.covalenthq.com/v1/1/tokens/0x3883f5e181fccaf8410fa61e12b59bad963fb645/token_holders/?quote-currency=USD&format=JSON&key=ckey_docs
async function fetchHolders(network, token) {
  const url = `https://api.covalenthq.com/v1/${network}/tokens/${token}/token_holders/?` +
    `quote-currency=USD&format=JSON&key=${process.env.COVALENT_API_KEY}&page-number=0&page-size=100000`
  return fetch(url)
    .then(response => response.json())
    .then(body => {
      if (body.error || !body.data) {
        throw new Error(body.error_message)
      }
      return body.data.items
    })
}

async function main() {
  if (process.argv.length < 4) throw new Error('wrong arguments: pass network_id and token_address')
  const network = parseInt(process.argv[2])
  const token = process.argv[3]

  return fetchHolders(network, token)
    .then(holders => {
      console.log(`Fetched ${holders.length} holders!`)
    })
    .catch(error => {
      console.error(`Failed to fetch holders: ${error}`)
    })
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.log(error)
    process.exit(1)
  })