const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE result ALTER COLUMN title DROP DEFAULT;
      ALTER TABLE result ALTER COLUMN description DROP DEFAULT;
    `);
  },
  down: async () => {
    console.log("can't rollback")
  },
}
