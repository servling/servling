import { z } from 'zod';

/**
 * Configuration source identifiers
 */
export const ConfigSources = {
  FILE: 'file',
  ENV: 'env',
  DEFAULT: 'default'
} as const;

export type ConfigSource = typeof ConfigSources[keyof typeof ConfigSources] | string;

/**
 * Configuration adapter interface
 */
export interface ConfigAdapter {
  /** Unique identifier for the adapter */
  id: string;
  
  /** File extensions this adapter can handle (without the dot) */
  extensions: string[];
  
  /** Parse a string into an object */
  parse: (content: string) => any;
  
  /** Stringify an object into a string */
  stringify: (obj: any) => string;
  
  /** Check if this adapter can handle the given file path */
  canHandle: (filePath: string) => boolean;
}

/**
 * Configuration options
 */
export interface ConfigOptions<T> {
  /** Path to the configuration file */
  configPath?: string;
  
  /** Adapter to use for parsing/stringifying (auto-detected from file extension if not specified) */
  adapter?: ConfigAdapter;
  
  /** Additional adapters to register */
  adapters?: ConfigAdapter[];
  
  /** Whether to create a config file with defaults if none exists */
  createIfMissing?: boolean;
  
  /** Whether to load from environment variables if config file is missing */
  useEnvFallback?: boolean;
  
  /** Environment variable prefix for loading from env vars */
  envPrefix?: string;
  
  /** Whether to load from .env file */
  useDotEnv?: boolean;
  
  /** Path to .env file (defaults to .env in process.cwd()) */
  dotEnvPath?: string;
}

/**
 * Schema factory function type
 * This allows for dynamic schema creation
 */
export type SchemaFactory<T> = () => ConfigSchema<T>;

/**
 * Configuration loader result
 */
export interface ConfigResult<T> {
  /** The loaded configuration */
  config: T;
  /** The source of each configuration value */
  sources: Record<keyof T, ConfigSource>;
  /** The path to the configuration file if loaded from file */
  configPath?: string;
}

/**
 * Configuration schema type
 */
export type ConfigSchema<T> = z.ZodType<T>;