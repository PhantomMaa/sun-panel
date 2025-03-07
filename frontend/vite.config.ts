import path from 'path'
import type { PluginOption } from 'vite'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

function setupPlugins(env: ImportMetaEnv): PluginOption[] {
  return [
    vue(),
    env.VITE_GLOB_APP_PWA === 'true' && VitePWA({
      injectRegister: 'auto',
      manifest: {
        name: 'Sun-Panel',
        short_name: 'Sun-Panel',
        icons: [
          { src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png' },
          { src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png' },
        ],
      },
    }),
    createSvgIconsPlugin({
      // 更新 SVG 图标目录路径，指向 public/assets/svg-icons
      iconDirs: [path.resolve(process.cwd(), 'public/assets/svg-icons')],
      symbolId: '[name]',
    }),
  ]
}

export default defineConfig((env) => {
  const viteEnv = loadEnv(env.mode, process.cwd()) as unknown as ImportMetaEnv

  // 定义公共资源目录
  const publicDir = path.resolve(process.cwd(), 'public')
  // 定义资源目录
  const assetsDir = 'assets'
  
  return {
    base: '/', // 确保所有环境下基础路径一致
    publicDir, // 设置公共资源目录
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
        // 添加资源路径别名，指向 public/assets 目录
        '/assets': path.resolve(process.cwd(), 'public/assets'),
      },
    },
    plugins: setupPlugins(viteEnv),
    // 资源处理配置
    css: {
      devSourcemap: true,
    },
    // 静态资源处理
    assetsInclude: ['**/*.svg', '**/*.png', '**/*.jpg', '**/*.jpeg', '**/*.gif', '**/*.webp'],
    server: {
      host: '0.0.0.0',
      port: 1002,
      open: false,
      proxy: {
        '/api': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/api/', '/api/'),
        },
        '/uploads': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/uploads/', '/uploads/'),
        },
      },
    },
    build: {
      outDir: path.resolve(__dirname, '../backend/web'),
      reportCompressedSize: false,
      sourcemap: false,
      // Set assets directory to match development environment
      assetsDir,
      // Copy public directory assets to output directory
      copyPublicDir: true,
      // Use esbuild for minification instead of Terser
      minify: 'esbuild',
      // Increase chunk size warning limit to reduce warnings
      chunkSizeWarningLimit: 600,
      // Empty the output directory before building
      emptyOutDir: true,
      commonjsOptions: {
        ignoreTryCatch: false,
      },
      // Configure drop_console for esbuild
      esbuildOptions: {
        drop: ['console'],
      },
      // Configure rollup options to ensure consistent asset paths
      rollupOptions: {
        output: {
          // Ensure asset filenames have consistent paths
          assetFileNames: `${assetsDir}/[name].[hash].[ext]`,
          chunkFileNames: `${assetsDir}/js/[name].[hash].js`,
          entryFileNames: `${assetsDir}/js/[name].[hash].js`,
          // Improve chunking to reduce large bundle sizes
          manualChunks: {
            'vue-vendor': ['vue', 'vue-router', 'pinia'],
            'ui-vendor': ['naive-ui'],
          },
        },
      },
    },
  }
})
