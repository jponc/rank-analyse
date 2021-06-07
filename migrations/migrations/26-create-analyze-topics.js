const { client } = require("../pg_client");

module.exports = {
  up: async () => {
    await client.query(`
      CREATE TABLE analyze_topics
        (
           id        UUID DEFAULT Uuid_generate_v4(),
           result_id UUID NOT NULL,
           label     TEXT NOT NULL,
           score     FLOAT NOT NULL,
           PRIMARY KEY(id),
           CONSTRAINT fk_result FOREIGN KEY(result_id) REFERENCES result(id)
        );
    `);
  },
  down: async () => {
    await client.query(`
      DROP TABLE analyze_topics;
    `);
  },
}
