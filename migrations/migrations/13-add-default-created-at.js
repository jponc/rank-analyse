const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN created_at SET DEFAULT now();
      ALTER TABLE result ALTER COLUMN created_at SET DEFAULT now();
    `);
  },
  down: async () => {
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN created_at DROP DEFAULT;
      ALTER TABLE result ALTER COLUMN created_at DROP DEFAULT;
    `);
  },
}
