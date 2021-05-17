const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE TABLE result
        (
           id          UUID,
           crawl_id    UUID NOT NULL,
           link        VARCHAR(40) NOT NULL,
           position    INTEGER NOT NULL,
           device      VARCHAR(40) NOT NULL,
           title       VARCHAR(40) NOT NULL DEFAULT false,
           description TEXT NOT NULL DEFAULT false,
           done        BOOLEAN NOT NULL DEFAULT false,
           created_at  TIMESTAMP NOT NULL,
           PRIMARY KEY(id),
           CONSTRAINT fk_crawl FOREIGN KEY(crawl_id) REFERENCES crawl(id)
        );
    `);
  },
  down: async () => {
    await client.query(`
      DROP TABLE result;
    `);
  },
}
