import { createConfig, z } from '../src';

// Define your configuration schema using Zod
const configSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    host: z.string().default('localhost')
  }),
  database: z.object({
    url: z.string(),
    maxConnections: z.number().default(10)
  }),
  api: z.object({
    key: z.string(),
    timeout: z.number().default(5000)
  })
});

// Set some environment variables for demonstration
process.env.MY_APP_SERVER_PORT = '4000';
process.env.MY_APP_DATABASE_URL = 'postgres://user:pass@remote-host:5432/prod-db';
process.env.MY_APP_API_KEY = 'secret-api-key';

// Create a configuration loader with environment variable prefix
const loadConfig = createConfig(configSchema, {
  configPath: './env-example-config.json',
  createIfMissing: true,
  useEnvFallback: true,
  envPrefix: 'MY_APP',
  useDotEnv: false // We're setting env vars manually for this example
});

// Load the configuration
const { config, sources, configPath } = loadConfig();

// Display the loaded configuration
console.log('Loaded configuration:');
console.log(JSON.stringify(config, null, 2));

// Display where each configuration value came from
console.log('\nConfiguration sources:');
console.log(JSON.stringify(sources, null, 2));

// Display the configuration file path if applicable
if (configPath) {
  console.log(`\nConfiguration file: ${configPath}`);
}

// You should see that some values came from environment variables
// and others from defaults