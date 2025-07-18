import { defineConfig } from 'vite';
import nodePolyfills from 'rollup-plugin-node-polyfills';

export default defineConfig({
  define: {
    global: 'window'
  },
  optimizeDeps: {
    esbuildOptions: {
      define: {
        global: 'window'
      }
    }
  },
  build: {
    rollupOptions: {
      plugins: [
        nodePolyfills()
      ]
    }
  }
}); 