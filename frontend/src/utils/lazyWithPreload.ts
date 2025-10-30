import { lazy, ComponentType, LazyExoticComponent } from 'react';

interface PreloadableComponent<T extends ComponentType<any>> extends LazyExoticComponent<T> {
  preload: () => Promise<{ default: T }>;
}

export function lazyWithPreload<T extends ComponentType<any>>(
  factory: () => Promise<{ default: T }>
): PreloadableComponent<T> {
  const Component = lazy(factory) as PreloadableComponent<T>;
  Component.preload = factory;
  return Component;
}

// Preload on hover or intersection
export function preloadComponent(Component: PreloadableComponent<any>) {
  if (Component.preload) {
    Component.preload();
  }
}