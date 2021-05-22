const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE crawl ADD COLUMN email VARCHAR(100) NOT NULL;
    `);
  },
  down: async () => {
    await client.query(`
      ALTER TABLE crawl DROP COLUMN email;
    `);
  },
}
