const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE extract_info ADD COLUMN cleaned_text TEXT NOT NULL DEFAULT '';
    `);
  },
  down: async () => {
    await client.query(`
      ALTER TABLE extract_info DROP COLUMN cleaned_text;
    `);
  },
}
