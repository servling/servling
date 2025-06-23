import { z } from 'zod';
import { ConfigLoader } from './config';
import {
  ConfigAdapter,
  ConfigOptions,
  ConfigResult,
  ConfigSchema,
  ConfigSource,
  ConfigSources,
  SchemaFactory
} from './types';
import {
  jsonAdapter,
  yamlAdapter,
  AdapterRegistry,
  createRegistry
} from './parsers';

/**
 * Creates a configuration loader with the given schema and options
 *
 * @param schema Zod schema or schema factory function for the configuration
 * @param options Configuration options
 * @returns A function that loads the configuration
 */
export function createConfig<T extends Record<string, any>>(
  schema: ConfigSchema<T> | SchemaFactory<T>,
  options: ConfigOptions<T> = {}
): () => ConfigResult<T> {
  // If schema is a factory function, call it to get the actual schema
  const actualSchema = typeof schema === 'function' ? schema() : schema;
  const loader = new ConfigLoader<T>(actualSchema, options);
  return () => loader.load();
}

// Export types and utilities
export {
  ConfigLoader,
  ConfigAdapter,
  ConfigOptions,
  ConfigResult,
  ConfigSchema,
  ConfigSource,
  ConfigSources,
  SchemaFactory,
  
  // Adapters
  jsonAdapter,
  yamlAdapter,
  AdapterRegistry,
  createRegistry
};

// Re-export zod for convenience
export { z };