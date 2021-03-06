const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    `);
  },
  down: async () => {
    await client.query(`
      DROP EXTENSION "uuid-ossp";
    `)
  },
}
