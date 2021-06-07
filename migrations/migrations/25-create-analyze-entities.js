const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE TABLE analyze_entities
        (
           id               UUID DEFAULT uuid_generate_v4(),
           result_id        UUID NOT NULL,
           confidence_score FLOAT NOT NULL,
           relevance_score  FLOAT NOT NULL,
           entity           TEXT NOT NULL,
           matched_text     TEXT NOT NULL,
           PRIMARY KEY(id),
           CONSTRAINT fk_result FOREIGN KEY(result_id) REFERENCES result(id)
        );
    `);
  },
  down: async () => {
    await client.query(`
      DROP TABLE analyze_entities;
    `);
  },
}
