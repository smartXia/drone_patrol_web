import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      port: 3000,
      open: true,
      proxy: {
        '/api': {
          target: 'http://180.102.24.183:18080',
          changeOrigin: true,
          secure: false
        },
        '/ws': {
          target: 'ws:/180.102.24.183:18080',
          ws: true,
          changeOrigin: true
        }
      }
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
      assetsDir: 'assets',
      base: './',
      rollupOptions: {
        output: {
          assetFileNames: 'assets/[name]-[hash][extname]',
          chunkFileNames: 'assets/[name]-[hash].js',
          entryFileNames: 'assets/[name]-[hash].js'
        }
      }
    },
    define: {
      // 确保环境变量在构建时被正确替换
      'import.meta.env.VITE_API_BASE_URL': JSON.stringify(env.VITE_API_BASE_URL || 'http://180.102.24.183:18080'),
      'import.meta.env.VITE_WS_BASE_URL': JSON.stringify(env.VITE_WS_BASE_URL || 'ws://180.102.24.183:18080'),
      'import.meta.env.VITE_ENV': JSON.stringify(env.VITE_ENV || 'production'),
      'import.meta.env.VITE_DEV_MODE': JSON.stringify(env.VITE_DEV_MODE || 'false'),
      'import.meta.env.VITE_DEVICE_SN': JSON.stringify(env.VITE_DEVICE_SN || '8UUXN2B00A00ST'),
      'import.meta.env.VITE_GATEWAY_SN': JSON.stringify(env.VITE_GATEWAY_SN || '8UUXN2B00A00ST')
    }
  }
})
