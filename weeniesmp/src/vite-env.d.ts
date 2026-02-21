/// <reference types="vite/client" />
/// <reference types="vite-imagetools/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '*&format=webp' {
  const src: string
  export default src
}

declare module '*&format=jpg' {
  const src: string
  export default src
}
