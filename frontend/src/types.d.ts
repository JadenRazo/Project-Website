/// <reference types="react" />
/// <reference types="react-dom" />

import React from 'react';

declare global {
  namespace JSX {
    interface IntrinsicElements {
      div: React.DetailedHTMLProps<React.HTMLAttributes<HTMLDivElement>, HTMLDivElement>;
      section: React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
      [elemName: string]: any;
    }
  }
}

declare module 'react' {
  interface SuspenseProps {
    fallback: React.ReactNode;
    children: React.ReactNode;
  }

  interface FunctionComponent<P = {}> {
    (props: P, context?: any): React.ReactElement<any, any> | null;
    displayName?: string;
  }

  interface LazyExoticComponent<T extends ComponentType<any>> {
    (props: ComponentProps<T>): React.ReactElement | null;
    _result: T;
  }
}

export {}; 