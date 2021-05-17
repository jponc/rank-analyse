const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE result DROP COLUMN device;
    `);
  },
  down: async () => {
    console.log("can't rollback")
  },
}
