const ServerlessClient = require('serverless-postgres')

module.exports = {
  client: new ServerlessClient({
    connectionString: process.env.DB_CONNECTION_URL
  }),
  getClient: () => {
    return new ServerlessClient({
      connectionString: process.env.DB_CONNECTION_URL
    })
  }
};
