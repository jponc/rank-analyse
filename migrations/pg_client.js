const { Client } = require('pg')
const client = new Client({
  connectionString: process.env.DB_CONNECTION_URL
})

module.exports = {
  client: client,
};
