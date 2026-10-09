import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { resolve } from 'path';

export default defineConfig({
  plugins: [
    tailwindcss(),
    svelte()
  ],
  resolve: {
    alias: {
      '$types': resolve(__dirname, './src/types'),
      '$components': resolve(__dirname, './src/components'),
      '$layouts': resolve(__dirname, './src/layouts'),
      '$views': resolve(__dirname, './src/views'),
      '$state': resolve(__dirname, './src/state')
    }
  }
});
