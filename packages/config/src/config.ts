import * as fs from 'fs-extra';
import * as path from 'path';
import { z } from 'zod';
import {
  ConfigAdapter,
  ConfigOptions,
  ConfigResult,
  ConfigSchema,
  ConfigSource,
  ConfigSources
} from './types';
import {
  AdapterRegistry,
  createRegistry
} from './parsers';
import {
  ensureDirectoryExists,
  flattenObject,
  getFromEnv,
  loadDotEnv,
  unflattenObject
} from './utils';

/**
 * Configuration loader class
 */
export class ConfigLoader<T extends Record<string, any>> {
  private schema: ConfigSchema<T>;
  private options: Required<ConfigOptions<T>>;
  private registry: AdapterRegistry;

  /**
   * Creates a new configuration loader
   *
   * @param schema Zod schema for the configuration
   * @param options Configuration options
   */
  constructor(schema: ConfigSchema<T>, options: ConfigOptions<T> = {}) {
    this.schema = schema;
    
    // Create adapter registry
    this.registry = createRegistry(options.adapters || []);
    
    // Determine the adapter to use
    let adapter: ConfigAdapter | undefined = options.adapter;
    
    if (!adapter && options.configPath) {
      adapter = this.registry.getForFile(options.configPath);
    }
    
    // Default to JSON if no adapter is found
    if (!adapter) {
      adapter = this.registry.getById('json');
    }
    
    if (!adapter) {
      throw new Error('No suitable adapter found for configuration');
    }
    
    // Set default options
    this.options = {
      configPath: options.configPath || path.resolve(process.cwd(), `config.${adapter.extensions[0]}`),
      adapter,
      adapters: options.adapters || [],
      createIfMissing: options.createIfMissing !== false,
      useEnvFallback: options.useEnvFallback !== false,
      envPrefix: options.envPrefix || '',
      useDotEnv: options.useDotEnv !== false,
      dotEnvPath: options.dotEnvPath || path.resolve(process.cwd(), '.env'),
    };
  }

  /**
   * Loads the configuration from the specified sources
   *
   * @returns The loaded configuration and metadata
   */
  public load(): ConfigResult<T> {
    // Load from .env file if enabled
    if (this.options.useDotEnv) {
      loadDotEnv(this.options.dotEnvPath);
    }

    // Try to load from file first
    const fileConfig = this.loadFromFile();
    
    // If file config exists, validate and return it
    if (fileConfig) {
      const validatedConfig = this.validateConfig(fileConfig);
      return {
        config: validatedConfig,
        sources: this.getConfigSources(validatedConfig, ConfigSources.FILE),
        configPath: this.options.configPath
      };
    }

    // If no file config and env fallback is enabled, try to load from env
    if (this.options.useEnvFallback) {
      const envConfig = this.loadFromEnv();
      const validatedConfig = this.validateConfig(envConfig.config);
      
      // Create config file if enabled
      if (this.options.createIfMissing) {
        this.saveConfig(validatedConfig);
      }
      
      return {
        config: validatedConfig,
        sources: envConfig.sources,
        configPath: this.options.createIfMissing ? this.options.configPath : undefined
      };
    }

    // If no file config and no env fallback, use schema defaults
    if (this.options.createIfMissing) {
      // Create an empty object and let Zod apply defaults
      const defaultConfig = this.validateConfig({});
      this.saveConfig(defaultConfig);
      
      return {
        config: defaultConfig,
        sources: this.getConfigSources(defaultConfig, ConfigSources.DEFAULT),
        configPath: this.options.configPath
      };
    }

    // If we get here, we couldn't load a config and couldn't create one
    throw new Error(
      `Could not load configuration from ${this.options.configPath} ` +
      `and no fallback options were available`
    );
  }

  /**
   * Loads configuration from a file
   *
   * @returns The loaded configuration or undefined if file doesn't exist
   */
  private loadFromFile(): T | undefined {
    try {
      if (fs.existsSync(this.options.configPath)) {
        const fileContent = fs.readFileSync(this.options.configPath, 'utf8');
        return this.options.adapter.parse(fileContent) as T;
      }
    } catch (error) {
      console.warn(
        `Failed to load configuration from ${this.options.configPath}:`,
        error
      );
    }
    
    return undefined;
  }

  /**
   * Loads configuration from environment variables
   *
   * @returns The loaded configuration and sources
   */
  private loadFromEnv(): {
    config: Partial<T>;
    sources: Record<keyof T, ConfigSource>
  } {
    // Create an empty object to collect environment variables
    const result: Record<string, any> = {};
    const sources: Record<string, ConfigSource> = {};
    
    // Get environment variables with the specified prefix
    for (const key in process.env) {
      if (key.startsWith(this.options.envPrefix)) {
        // Remove the prefix and convert to camelCase
        const configKey = key.substring(this.options.envPrefix.length)
          .toLowerCase()
          .replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
        
        // Parse the value
        let value: any = process.env[key];
        
        // Try to parse as JSON if it looks like JSON
        if (
          (value?.startsWith('{') && value?.endsWith('}')) ||
          (value?.startsWith('[') && value?.endsWith(']')) ||
          value === 'true' ||
          value === 'false' ||
          !isNaN(Number(value))
        ) {
          try {
            value = JSON.parse(value);
          } catch {
            // If parsing fails, use the string value
          }
        }
        
        result[configKey] = value;
        sources[configKey] = ConfigSources.ENV;
      }
    }
    
    // Unflatten the result
    return {
      config: unflattenObject(result) as Partial<T>,
      sources: sources as Record<keyof T, ConfigSource>
    };
  }

  /**
   * Validates a configuration against the schema
   *
   * @param config The configuration to validate
   * @returns The validated configuration
   */
  private validateConfig(config: Partial<T>): T {
    try {
      // First attempt: try to parse with the schema directly
      return this.schema.parse(config);
    } catch (error) {
      if (error instanceof z.ZodError) {
        // For schema evolution, we want to handle the case where new fields
        // with defaults have been added to the schema but don't exist in the config
        
        // Create a deep copy of the config to avoid modifying the original
        const enhancedConfig = JSON.parse(JSON.stringify(config || {}));
        
        // First, try to get default values by parsing an empty object
        // This will give us the schema's default values for all fields
        let defaultValues: any = {};
        try {
          defaultValues = this.schema.parse({});
        } catch {
          // If parsing an empty object fails, we'll continue without defaults
        }
        
        // Process each validation error
        for (const issue of error.issues) {
          // Only handle required field errors
          if (issue.code === 'invalid_type' && issue.message.includes('Required')) {
            const pathStr = issue.path.join('.');
            
            // Build the path to the missing field
            let current = enhancedConfig;
            let defaultCurrent = defaultValues;
            const path = [...issue.path];
            
            // Create all parent objects in the path
            for (let i = 0; i < path.length - 1; i++) {
              const segment = path[i] as string | number;
              
              // Create the path in the enhanced config if it doesn't exist
              if (!(segment in current) || current[segment] === null) {
                current[segment] = {};
              }
              current = current[segment];
              
              // Navigate the default values if they exist
              if (defaultCurrent && typeof defaultCurrent === 'object' && segment in defaultCurrent) {
                defaultCurrent = defaultCurrent[segment];
              } else {
                defaultCurrent = undefined;
              }
            }
            
            // Set the final property
            const lastSegment = path[path.length - 1] as string | number;
            if (!(lastSegment in current)) {
              // If we have default values for this path, use them
              if (defaultCurrent && typeof defaultCurrent === 'object' && lastSegment in defaultCurrent) {
                current[lastSegment] = defaultCurrent[lastSegment];
              } else {
                // Otherwise use an empty object
                current[lastSegment] = {};
              }
            }
          }
        }
        
        // Try to parse with the enhanced config
        try {
          return this.schema.parse(enhancedConfig);
        } catch (enhancedError) {
          // If we still have errors, report them
          if (enhancedError instanceof z.ZodError) {
            const issues = enhancedError.issues.map(issue =>
              `${issue.path.join('.')}: ${issue.message}`
            ).join('\n');
            
            throw new Error(
              `Configuration validation failed:\n${issues}`
            );
          }
          
          throw enhancedError;
        }
      }
      
      throw error;
    }
  }

  /**
   * Saves a configuration to a file
   *
   * @param config The configuration to save
   */
  private saveConfig(config: T): void {
    try {
      ensureDirectoryExists(this.options.configPath);
      fs.writeFileSync(
        this.options.configPath,
        this.options.adapter.stringify(config),
        'utf8'
      );
    } catch (error) {
      console.warn(
        `Failed to save configuration to ${this.options.configPath}:`,
        error
      );
    }
  }

  /**
   * Gets the sources for each configuration value
   *
   * @param config The configuration
   * @param defaultSource The default source to use
   * @returns The sources for each configuration value
   */
  private getConfigSources(
    config: T,
    defaultSource: ConfigSource
  ): Record<keyof T, ConfigSource> {
    const sources: Record<string, ConfigSource> = {};
    const flatConfig = flattenObject(config);
    
    for (const key in flatConfig) {
      sources[key] = defaultSource;
    }
    
    return unflattenObject(sources) as Record<keyof T, ConfigSource>;
  }
}