import { createConfig, z } from '../src';

// Define your configuration schema using Zod
const configSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    host: z.string().default('localhost')
  }),
  database: z.object({
    url: z.string().default('postgres://localhost:5432/mydb'),
    maxConnections: z.number().default(10)
  }),
  logging: z.object({
    level: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
    format: z.enum(['json', 'text']).default('json')
  })
});

// Create a configuration loader
const loadConfig = createConfig(configSchema, {
  configPath: './example-config.json',
  createIfMissing: true,
  useEnvFallback: true,
  envPrefix: 'APP',
  useDotEnv: true
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