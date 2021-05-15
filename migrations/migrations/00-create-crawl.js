const { getClient } = require("../pg_client");

module.exports = {
  up: async () => {
    const client = getClient();

    await client.connect();
    await client.query(`
      CREATE TABLE crawl
        (
           id            UUID,
           keyword       VARCHAR(40) NOT NULL,
           search_engine VARCHAR(40) NOT NULL,
           device        VARCHAR(40) NOT NULL,
           done          BOOLEAN NOT NULL DEFAULT false,
           created_at    TIMESTAMP NOT NULL,
           PRIMARY KEY(id)
        );
    `);

    await client.clean();
  },
  down: async () => {
    await client.connect();
    await client.query(`
      DROP TABLE crawl;
    `);

    await client.clean();
  },
}
