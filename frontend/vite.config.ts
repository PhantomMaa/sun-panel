import path from 'path'
import type { PluginOption } from 'vite'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

// 资源路径处理插件 - 优化版本
function createAssetsPlugin(): PluginOption {
  return {
    name: 'assets-plugin',
    configureServer(server) {
      // 添加中间件处理资源路径
      server.middlewares.use((req, res, next) => {
        if (!req.url) {
          return next();
        }
        
        let modified = false;
        let newUrl = req.url;
        
        // 处理 /@fs/assets/ 路径 - 主要用于开发环境
        if (newUrl.includes('/@fs/assets/')) {
          newUrl = newUrl.replace('/@fs/assets/', '/assets/');
          modified = true;
        }
        
        // 保留对 /src/assets/ 路径的处理作为安全网
        // 虽然我们已经删除了 src/assets 目录，但某些代码可能仍然引用它
        if (newUrl.includes('/src/assets/')) {
          newUrl = newUrl.replace('/src/assets/', '/assets/');
          modified = true;
        }
        
        if (modified) {
          console.log(`Rewriting URL from ${req.url} to ${newUrl}`);
          req.url = newUrl;
        }
        
        next();
      });
    },
    transform(code, id) {
      if (!code) return code;
      
      // 只有当代码中包含需要替换的路径时才进行处理
      if (!code.includes('/assets/') && 
          !code.includes('/@fs/assets/') && 
          !code.includes('/src/assets/') && 
          !code.includes('http://127.0.0.1:') && 
          !code.includes('http://localhost:')) {
        return code;
      }
      
      let newCode = code;
      
      // 统一替换所有资源路径格式为标准的 /assets/ 格式
      const replacements = [
        // 基本路径替换
        { pattern: /\/src\/assets\//g, replacement: '/assets/' },
        { pattern: /\/@fs\/assets\//g, replacement: '/assets/' },
        
        // 带有特定端口的路径替换
        { pattern: /http:\/\/127\.0\.0\.1:\d+\/@fs\/assets\//g, replacement: '/assets/' },
        { pattern: /http:\/\/localhost:\d+\/@fs\/assets\//g, replacement: '/assets/' }
      ];
      
      // 应用所有替换
      for (const { pattern, replacement } of replacements) {
        newCode = newCode.replace(pattern, replacement);
      }
      
      return newCode !== code ? newCode : code;
    }
  };
}

function setupPlugins(env: ImportMetaEnv): PluginOption[] {
  return [
    vue(),
    // 添加全面的资源路径处理插件，支持开发和生产环境
    createAssetsPlugin(),
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
        '/@fs/assets': path.resolve(process.cwd(), 'public/assets'),
        // 兼容旧的资源路径
        '/src/assets': path.resolve(process.cwd(), 'public/assets'),
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
      // 开发服务器静态资源处理
      fs: {
        strict: false,
        allow: [
          // 允许访问项目根目录
          path.resolve(process.cwd()),
          // 显式允许访问 src/assets 目录
          path.resolve(process.cwd(), 'src/assets'),
        ],
      },
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
