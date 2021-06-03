const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE crawl DROP COLUMN email;
    `);
  },
  down: async () => {
    console.log("can't rollback")
  },
}
