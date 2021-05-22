const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE extract_info DROP COLUMN links;
    `);
  },
  down: async () => {
    console.log("can't rollback")
  },
}
