import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { imagetools } from 'vite-imagetools';
import { fileURLToPath, URL } from 'node:url';
export default defineConfig({
    plugins: [vue(), imagetools()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    server: {
        port: 5173,
        host: '0.0.0.0',
        strictPort: true,
        hmr: {
            host: '195.201.136.53',
            port: 5173,
            protocol: 'ws'
        }
    },
    build: {
        outDir: 'dist',
        sourcemap: false
    }
});
