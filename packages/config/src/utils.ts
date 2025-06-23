import * as fs from 'fs-extra';
import * as path from 'path';
import * as dotenv from 'dotenv';
import { ConfigSource, ConfigSources } from './types';

/**
 * Converts a camelCase or PascalCase string to SNAKE_CASE
 */
export function toSnakeCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1_$2')
    .replace(/([A-Z])([A-Z][a-z])/g, '$1_$2')
    .toUpperCase();
}

/**
 * Converts a nested object path to an environment variable name
 * e.g. "database.host" -> "DATABASE_HOST"
 */
export function toEnvVarName(path: string, prefix?: string): string {
  const envName = path.split('.').join('_').toUpperCase();
  return prefix ? `${prefix}_${envName}` : envName;
}

/**
 * Loads environment variables from .env file if it exists
 */
export function loadDotEnv(dotEnvPath?: string): void {
  const envPath = dotEnvPath || path.resolve(process.cwd(), '.env');
  if (fs.existsSync(envPath)) {
    dotenv.config({ path: envPath });
  }
}

/**
 * Gets a value from environment variables based on a path
 * e.g. "database.host" -> process.env.DATABASE_HOST
 */
export function getFromEnv<T>(
  path: string, 
  prefix?: string, 
  defaultValue?: T
): { value: T | undefined; source: ConfigSource } {
  const envName = toEnvVarName(path, prefix);
  const envValue = process.env[envName];
  
  if (envValue !== undefined) {
    // Try to parse the value if it looks like JSON
    if (
      (envValue.startsWith('{') && envValue.endsWith('}')) ||
      (envValue.startsWith('[') && envValue.endsWith(']')) ||
      envValue === 'true' ||
      envValue === 'false' ||
      !isNaN(Number(envValue))
    ) {
      try {
        return {
          value: JSON.parse(envValue) as T,
          source: ConfigSources.ENV
        };
      } catch {
        // If parsing fails, use the string value
        return {
          value: envValue as unknown as T,
          source: ConfigSources.ENV
        };
      }
    }
    
    return {
      value: envValue as unknown as T,
      source: ConfigSources.ENV
    };
  }
  
  return {
    value: defaultValue,
    source: defaultValue !== undefined ? ConfigSources.DEFAULT : ConfigSources.ENV
  };
}

/**
 * Ensures a directory exists
 */
export function ensureDirectoryExists(filePath: string): void {
  const dirname = path.dirname(filePath);
  if (!fs.existsSync(dirname)) {
    fs.mkdirpSync(dirname);
  }
}

/**
 * Flattens a nested object into a flat object with dot notation paths
 * e.g. { a: { b: 1 } } -> { "a.b": 1 }
 */
export function flattenObject(
  obj: Record<string, any>, 
  prefix = ''
): Record<string, any> {
  return Object.keys(obj).reduce((acc, key) => {
    const prefixedKey = prefix ? `${prefix}.${key}` : key;
    
    if (
      typeof obj[key] === 'object' && 
      obj[key] !== null && 
      !Array.isArray(obj[key])
    ) {
      Object.assign(acc, flattenObject(obj[key], prefixedKey));
    } else {
      acc[prefixedKey] = obj[key];
    }
    
    return acc;
  }, {} as Record<string, any>);
}

/**
 * Unflatten a flat object with dot notation paths into a nested object
 * e.g. { "a.b": 1 } -> { a: { b: 1 } }
 */
export function unflattenObject(
  obj: Record<string, any>
): Record<string, any> {
  const result: Record<string, any> = {};
  
  for (const key in obj) {
    const keys = key.split('.');
    let current = result;
    
    for (let i = 0; i < keys.length; i++) {
      const k = keys[i];
      if (i === keys.length - 1) {
        current[k] = obj[key];
      } else {
        current[k] = current[k] || {};
        current = current[k];
      }
    }
  }
  
  return result;
}