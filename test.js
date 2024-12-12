const apiKey = ''
const prod = ''
const code = ''
fetch(`${prod}/auth?code=${code}`, {
  headers: {
    'Api-Key': apiKey
  }
}).then(r => r.text())
  .then(r => console.log(r))
  .catch(e => console.log('error', e))
