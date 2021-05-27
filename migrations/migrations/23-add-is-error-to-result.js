const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE result ADD COLUMN is_error BOOLEAN NOT NULL DEFAULT false;
    `);
  },
  down: async () => {
    await client.query(`
      ALTER TABLE result DROP COLUMN is_error;
    `);
  },
}
