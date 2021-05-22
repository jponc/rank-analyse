const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE TABLE extract_links
        (
           id         UUID NOT NULL DEFAULT uuid_generate_v4(),
           result_id  UUID NOT NULL,
           text       TEXT,
           link_url   TEXT NOT NULL,
           created_at TIMESTAMP NOT NULL DEFAULT now(),
           PRIMARY KEY(id),
           CONSTRAINT fk_result FOREIGN KEY(result_id) REFERENCES result(id)
        );
    `);
  },
  down: async () => {
    await client.query(`
      DROP TABLE extract_links;
    `);
  },
}
