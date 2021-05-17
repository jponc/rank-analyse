const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE TABLE extract_info
        (
           id         UUID,
           result_id  UUID NOT NULL,
           title      VARCHAR(40) NOT NULL,
           content    TEXT NOT NULL,
           links      TEXT[] DEFAULT '{}',
           created_at TIMESTAMP NOT NULL DEFAULT NOW(),
           PRIMARY KEY(id),
           CONSTRAINT fk_result FOREIGN KEY(result_id) REFERENCES result(id)
        );
    `);
  },
  down: async () => {
    await client.query(`
      DROP TABLE extract_info;
    `);
  },
}
