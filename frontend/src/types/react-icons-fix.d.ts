declare module 'react-icons/fa' {
  import * as React from 'react';
  export interface IconProps extends React.SVGProps<SVGSVGElement> {
    size?: string | number;
  }

  export const FaLaptopCode: React.FC<IconProps>;
  export const FaGithub: React.FC<IconProps>;
  export const FaExternalLinkAlt: React.FC<IconProps>;
  export const FaCode: React.FC<IconProps>;
  export const FaWindows: React.FC<IconProps>;
}

declare module 'react-icons/fc' {
  import * as React from 'react';
  export interface IconProps extends React.SVGProps<SVGSVGElement> {
    size?: string | number;
  }

  export const FcGoogle: React.FC<IconProps>;
} 