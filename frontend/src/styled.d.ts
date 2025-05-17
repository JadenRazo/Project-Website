// src/styled.d.ts
import 'styled-components';
import { Theme as ImportedTheme } from './styles/theme.types'; // Renamed to avoid conflict if Theme is defined below

declare module 'styled-components' {
  // Extend the imported theme or define a new one that includes it
  export interface DefaultTheme extends ImportedTheme {
    // Ensure fonts is part of the theme, making it optional if it might not always be present
    fonts?: {
      primary?: string;
      mono?: string;
    };
    // You can add other properties here if DefaultTheme needs to be richer than ImportedTheme
  }
}
