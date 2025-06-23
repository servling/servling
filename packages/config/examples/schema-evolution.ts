import * as fs from 'node:fs';
import { createConfig, z } from '../src';

console.log('Testing schema evolution with defaults...');

// Create a test configuration file with initial data
const initialConfig = {
  server: {
    port: 3000,
    host: 'localhost',
  },
  database: {
    url: 'postgres://localhost:5432/mydb',
  }
};

// Write the initial config file
fs.writeFileSync('./evolution-test-config.json', JSON.stringify(initialConfig, null, 2));
console.log('Created initial configuration file');

// Step 1: Load the configuration with the initial schema
const initialSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    host: z.string().default('localhost'),
  }),
  database: z.object({
    url: z.string().default('postgres://localhost:5432/mydb'),
  }),
});

// Create a configuration loader with the initial schema
const initialLoader = createConfig(initialSchema, {
  configPath: './evolution-test-config.json',
  createIfMissing: false, // File already exists
});

// Load the initial configuration
const initialResult = initialLoader();
console.log('\nInitial configuration:');
console.log(JSON.stringify(initialResult.config, null, 2));

// Step 2: "Evolve" the schema by adding new fields with defaults
const evolvedSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    host: z.string().default('localhost'),
    // New field with default
    timeout: z.number().default(30000),
  }),
  database: z.object({
    url: z.string().default('postgres://localhost:5432/mydb'),
    // New field with default
    maxConnections: z.number().default(10),
  }),
  // Entirely new section with defaults
  logging: z.object({
    level: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
    format: z.enum(['json', 'text']).default('json'),
  }),
});

// Create a configuration loader with the evolved schema
const evolvedLoader = createConfig(evolvedSchema, {
  configPath: './evolution-test-config.json',
  createIfMissing: false, // Don't create a new file, use the existing one
});

// Load the configuration with the evolved schema
// This should not error even though the file is missing the new fields
// Instead, it should use the default values for the new fields
const evolvedResult = evolvedLoader();
console.log('\nEvolved configuration (should include new fields with defaults):');
console.log(JSON.stringify(evolvedResult.config, null, 2));

// Clean up the test file
fs.unlinkSync('./evolution-test-config.json');
console.log('\nTest completed and test file removed.');