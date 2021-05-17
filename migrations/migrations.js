const { client } = require("./pg_client");
const Umzug = require('umzug');
const fs = require('fs');

module.exports = {
  up: async () => {
    await client.connect();

    const migrationStorePath = '/tmp/migrations.json';
    const umzug = new Umzug({
      storage: 'json',
      storageOptions: {
        path: migrationStorePath,
      },
      migrations: {
        path: './migrations'
      }
    });

    const migrationTableQuery = `
      CREATE TABLE IF NOT EXISTS migration (
        created_at TIMESTAMP DEFAULT NOW(),
        name       VARCHAR(100) NOT NULL
      )
    `;
    await client.query(migrationTableQuery);

    const res = await client.query('SELECT * FROM migration ORDER BY created_at');
    const migrations = res.rows;

    // transform migrations into an umzug accepted format
    const umzugMigrations = migrations.map(migration => migration.name);
    // write into umzug migration store
    fs.writeFileSync(migrationStorePath, JSON.stringify(umzugMigrations), { flag: 'w' });

    // When a migration file is executed, add it to our migration table
    function addMigration() {
      return async (name) => {
        console.log(`${name} migrated`);
        try {
          const migrationQueryInsert = `INSERT INTO migration (name) VALUES('${name}.js')`;
          console.log(migrationQueryInsert);
          await client.query(migrationQueryInsert);
          console.log(`${name} inserted into migration table`);
        } catch (error) {
          console.log(error);
          console.error(`${name} could not be inserted to migration table`);
          throw new Error("failed!");
        }
      };
    }
    // When umzug finished to execute a migration file, call addMigration
    umzug.on('migrated', addMigration());

    try {
      const result = await umzug.up()
      console.log(`Migration completed: ${JSON.stringify(result)}`);
    } catch (err) {
      throw err;
    }

    // This is a hack to ensure the last migration has been inserted
    const delay = ms => new Promise(resolve => setTimeout(resolve, ms))
    await delay(3000) /// waiting 3 second.
    await client.end();
  }
}

