const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE EXTENSION "uuid-ossp";
    `);
  },
  down: async () => {
    await client.query(`
      DROP EXTENSION "uuid-ossp";
    `)
  },
}
