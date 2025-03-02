// src/styled.d.ts
import 'styled-components';
import { Theme } from './styles/theme.types';

declare module 'styled-components' {
  export interface DefaultTheme extends Theme {}
}
