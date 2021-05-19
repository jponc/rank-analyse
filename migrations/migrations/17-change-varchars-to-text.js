const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE result ALTER COLUMN title TYPE TEXT;
      ALTER TABLE result ALTER COLUMN link TYPE TEXT;
    `);
  },
  down: async () => {
    console.log("cannot rollback");
  },
}
