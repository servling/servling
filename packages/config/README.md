# @servling/config

A configuration management library for Node.js applications with Zod schema validation.

## Features

- 📝 **Schema Validation**: Uses [Zod](https://github.com/colinhacks/zod) for schema validation
- 🔄 **Multiple Sources**: Load from files, environment variables, or defaults
- 📄 **Multiple Formats**: Supports JSON and YAML out of the box with an extensible adapter system
- 🌱 **Auto-creation**: Creates config files with defaults if they don't exist
- 🔍 **Traceability**: Tracks where each config value came from
- 🌿 **Dotenv Support**: Loads environment variables from .env files
- 🧩 **Extensible**: Add your own adapters for custom file formats
- 🔄 **Schema Evolution**: Gracefully handles schema changes by using defaults for new fields

## Installation

```bash
npm install @servling/config
# or
yarn add @servling/config
# or
pnpm add @servling/config
```

## Usage

### Basic Example

```typescript
import { createConfig, jsonAdapter, yamlAdapter, z } from '@servling/config';

// Define your configuration schema using Zod
const configSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    host: z.string().default('localhost'),
  }),
  database: z.object({
    url: z.string(),
    maxConnections: z.number().default(10),
  }),
  logging: z.object({
    level: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
    format: z.enum(['json', 'text']).default('json'),
  }),
});

// Create a configuration loader with support for both JSON and YAML
const loadConfig = createConfig(configSchema, {
  configPath: './config/app.json', // Will auto-detect format based on extension
  adapters: [jsonAdapter, yamlAdapter], // Register the adapters you want to use
  createIfMissing: true,
  useEnvFallback: true,
  envPrefix: 'APP',
  useDotEnv: true,
});

// Load the configuration
const { config, sources } = loadConfig();

// Use the configuration
const server = app.listen(config.server.port, config.server.host);
console.log(`Server started at http://${config.server.host}:${config.server.port}`);
```

### Using a Schema Factory Function

You can also use a factory function to create your schema dynamically:

```typescript
import { createConfig, z } from '@servling/config';
import { generateKeys } from 'crypto';

// Define a schema factory function that generates fresh values each time
function createConfigSchema() {
  // Generate keys or other dynamic values inside the factory
  const secretKey = generateKeys();
  
  return z.object({
    security: z.object({
      secretKey: z.string().default(secretKey),
    }),
    server: z.object({
      port: z.number().default(3000),
    }),
  });
}

// Infer the type from the schema factory
type Config = z.infer<ReturnType<typeof createConfigSchema>>;

// Create a configuration loader with the schema factory
const loadConfig = createConfig(createConfigSchema, {
  configPath: './config.json',
  createIfMissing: true,
});

// Load the configuration
const { config } = loadConfig();
```

### Configuration Sources

The library will try to load configuration in the following order:

1. From the JSON file specified in `configPath`
2. From environment variables (if `useEnvFallback` is true)
3. From default values specified in the schema

If `createIfMissing` is true and the config file doesn't exist, it will be created with the loaded configuration.

### Schema Evolution

The library gracefully handles schema evolution. When you add new fields to your schema with default values, the configuration system will automatically use those defaults for existing configuration files without throwing errors.

This is particularly useful for:
- Adding new features that require new configuration options
- Evolving your application over time without breaking existing deployments
- Providing sensible defaults for new configuration options

Example:

```typescript
// Original schema
const originalSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
  }),
});

// Evolved schema with new fields
const evolvedSchema = z.object({
  server: z.object({
    port: z.number().default(3000),
    // New field with default
    timeout: z.number().default(30000),
  }),
  // Entirely new section with defaults
  logging: z.object({
    level: z.enum(['debug', 'info']).default('info'),
  }),
});

// Even if your config file was created with the original schema,
// loading it with the evolved schema will work fine, using defaults
// for the new fields
```

### Environment Variables

Environment variables are mapped to configuration properties using snake case:

- `server.port` → `APP_SERVER_PORT`
- `database.url` → `APP_DATABASE_URL`
- `logging.level` → `APP_LOGGING_LEVEL`

The prefix can be customized using the `envPrefix` option.

### Advanced Usage

#### Custom Configuration Loader

```typescript
import { ConfigLoader, z } from '@servling/config';

const schema = z.object({
  // your schema
});

const loader = new ConfigLoader(schema, {
  configPath: './config/custom.json',
  // other options
});

const result = loader.load();
```

#### Type Inference

Zod provides automatic type inference:

```typescript
import { createConfig, z } from '@servling/config';

const configSchema = z.object({
  // your schema
});

// TypeScript will infer the correct type
type AppConfig = z.infer<typeof configSchema>;

const loadConfig = createConfig(configSchema);
const { config } = loadConfig();

// config is typed as AppConfig
```

## API Reference

### `createConfig(schema, options)`

Creates a configuration loader with the given schema and options.

#### Parameters

- `schema`: A Zod schema or schema factory function that defines the structure and validation rules for your configuration
- `options`: Configuration options (optional)

#### Options

- `configPath`: Path to the configuration file (default: `./config.json`)
- `adapter`: Specific adapter to use for parsing/stringifying (auto-detected from file extension if not specified)
- `adapters`: Additional adapters to register for auto-detection
- `createIfMissing`: Whether to create a config file with defaults if none exists (default: `true`)
- `useEnvFallback`: Whether to load from environment variables if config file is missing (default: `true`)
- `envPrefix`: Environment variable prefix for loading from env vars (default: `''`)
- `useDotEnv`: Whether to load from .env file (default: `true`)
- `dotEnvPath`: Path to .env file (default: `./.env`)

#### Returns

A function that, when called, loads and returns the configuration.

### Return Value

The loader function returns an object with the following properties:

- `config`: The loaded configuration
- `sources`: An object with the same structure as the config, indicating where each value came from
- `configPath`: The path to the configuration file if loaded from or saved to a file

## Creating Custom Adapters

You can create your own adapters to support additional file formats:

```typescript
import { ConfigAdapter, createConfig, z } from '@servling/config';
import * as toml from 'toml';

// Create a TOML adapter
const tomlAdapter: ConfigAdapter = {
  id: 'toml',
  extensions: ['toml'],
  parse: (content: string) => toml.parse(content),
  stringify: (obj: any) => toml.stringify(obj),
  canHandle: (filePath: string) => filePath.toLowerCase().endsWith('.toml'),
};

// Use the custom adapter
const configSchema = z.object({
  // your schema
});

const loadConfig = createConfig(configSchema, {
  configPath: 'config.toml',
  adapters: [tomlAdapter], // Register your custom adapter
});

const { config } = loadConfig();
```

## License

MIT