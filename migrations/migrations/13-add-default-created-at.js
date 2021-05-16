const { getClient } = require("../pg_client");

module.exports = {
  up: async () => {
    const client = getClient();

    await client.connect();
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN created_at SET DEFAULT now();
      ALTER TABLE result ALTER COLUMN created_at SET DEFAULT now();
    `);

    await client.clean();
  },
  down: async () => {
    await client.connect();
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN creatd_at DROP DEFAULT;
      ALTER TABLE result ALTER COLUMN creatd_at DROP DEFAULT;
    `);

    await client.clean();
  },
}
