const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN id SET DEFAULT uuid_generate_v4();
      ALTER TABLE result ALTER COLUMN id SET DEFAULT uuid_generate_v4();
      ALTER TABLE extract_info ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    `);
  },
  down: async () => {
    await client.query(`
      ALTER TABLE crawl ALTER COLUMN id DROP DEFAULT;
      ALTER TABLE result ALTER COLUMN id DROP DEFAULT;
      ALTER TABLE extract_info ALTER COLUMN id DROP DEFAULT;
    `);
  },
}
