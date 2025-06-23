import * as yaml from 'js-yaml';
import { ConfigAdapter } from './types';

/**
 * JSON adapter
 */
export const jsonAdapter: ConfigAdapter = {
  id: 'json',
  extensions: ['json'],
  parse: (content: string) => JSON.parse(content),
  stringify: (obj: any) => JSON.stringify(obj, null, 2),
  canHandle: (filePath: string) => filePath.toLowerCase().endsWith('.json')
};

/**
 * YAML adapter
 */
export const yamlAdapter: ConfigAdapter = {
  id: 'yaml',
  extensions: ['yaml', 'yml'],
  parse: (content: string) => yaml.load(content),
  stringify: (obj: any) => yaml.dump(obj, { indent: 2 }),
  canHandle: (filePath: string) => {
    const lower = filePath.toLowerCase();
    return lower.endsWith('.yaml') || lower.endsWith('.yml');
  }
};

/**
 * Default adapters
 */
export const defaultAdapters: ConfigAdapter[] = [
  jsonAdapter,
  yamlAdapter
];

/**
 * Registry of adapters
 */
export class AdapterRegistry {
  private adapters: ConfigAdapter[] = [];

  constructor(initialAdapters: ConfigAdapter[] = []) {
    this.registerAll(initialAdapters);
  }

  /**
   * Register a new adapter
   */
  register(adapter: ConfigAdapter): void {
    // Don't register the same adapter twice
    if (!this.adapters.some(a => a.id === adapter.id)) {
      this.adapters.push(adapter);
    }
  }

  /**
   * Register multiple adapters
   */
  registerAll(adapters: ConfigAdapter[]): void {
    adapters.forEach(adapter => this.register(adapter));
  }

  /**
   * Get an adapter by ID
   */
  getById(id: string): ConfigAdapter | undefined {
    return this.adapters.find(adapter => adapter.id === id);
  }

  /**
   * Get an adapter that can handle the given file path
   */
  getForFile(filePath: string): ConfigAdapter | undefined {
    return this.adapters.find(adapter => adapter.canHandle(filePath));
  }

  /**
   * Get an adapter for a file extension
   */
  getForExtension(extension: string): ConfigAdapter | undefined {
    const ext = extension.startsWith('.') ? extension.substring(1) : extension;
    return this.adapters.find(adapter =>
      adapter.extensions.includes(ext.toLowerCase())
    );
  }

  /**
   * Get all registered adapters
   */
  getAll(): ConfigAdapter[] {
    return [...this.adapters];
  }
}

/**
 * Create a new adapter registry with default adapters
 */
export function createRegistry(additionalAdapters: ConfigAdapter[] = []): AdapterRegistry {
  const registry = new AdapterRegistry(defaultAdapters);
  registry.registerAll(additionalAdapters);
  return registry;
}