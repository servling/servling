import { createConfig, z } from '../src';

// Define a simple configuration schema
const configSchema = z.object({
  app: z.object({
    name: z.string().default('MyApp'),
    version: z.string().default('1.0.0')
  }),
  features: z.object({
    enableFeatureA: z.boolean().default(false),
    enableFeatureB: z.boolean().default(true)
  })
});

// Example of how to test with different configurations
function runTests() {
  console.log('Running tests with different configurations...\n');

  // Test with default configuration
  testConfig('Default Config', {});

  // Test with custom configuration
  testConfig('Custom Config', {
    defaults: {
      app: {
        name: 'TestApp',
        version: '2.0.0'
      },
      features: {
        enableFeatureA: true,
        enableFeatureB: false
      }
    }
  });

  // Test with partial overrides
  testConfig('Partial Overrides', {
    defaults: {
      features: {
        enableFeatureA: true
      }
    }
  });
}

// Helper function to test a configuration
function testConfig(name: string, options: any) {
  console.log(`=== ${name} ===`);
  
  const loadConfig = createConfig(configSchema, {
    configPath: `./test-${name.toLowerCase().replace(/\s+/g, '-')}.json`,
    createIfMissing: false,
    useEnvFallback: false,
    ...options
  });

  try {
    const { config } = loadConfig();
    console.log(JSON.stringify(config, null, 2));
  } catch (error) {
    console.error('Error:', error instanceof Error ? error.message : error);
  }
  
  console.log();
}

// Run the tests
runTests();