const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE extract_info ALTER COLUMN title TYPE TEXT;
    `);
  },
  down: async () => {
    console.log("cannot rollback");
  },
}
