const http = require('http');
const fs = require('fs');
const path = require('path');
const { exec } = require('child_process');
const os = require('os');

const PORT = 3000;

// MIME类型映射
const mimeTypes = {
  '.html': 'text/html',
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.png': 'image/png',
  '.jpg': 'image/jpg',
  '.gif': 'image/gif',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.eot': 'application/vnd.ms-fontobject'
};

// 创建HTTP服务器
const server = http.createServer((req, res) => {
  let filePath = path.join(__dirname, 'dist', req.url === '/' ? 'index.html' : req.url);
  
  // 处理Vue Router的history模式
  if (!fs.existsSync(filePath) || fs.statSync(filePath).isDirectory()) {
    filePath = path.join(__dirname, 'dist', 'index.html');
  }
  
  const extname = path.extname(filePath).toLowerCase();
  const contentType = mimeTypes[extname] || 'application/octet-stream';
  
  fs.readFile(filePath, (error, content) => {
    if (error) {
      if (error.code === 'ENOENT') {
        // 文件不存在，返回index.html（Vue Router处理）
        fs.readFile(path.join(__dirname, 'dist', 'index.html'), (err, data) => {
          if (err) {
            res.writeHead(500);
            res.end('服务器内部错误');
          } else {
            res.writeHead(200, { 'Content-Type': 'text/html' });
            res.end(data);
          }
        });
      } else {
        res.writeHead(500);
        res.end('服务器内部错误');
      }
    } else {
      res.writeHead(200, { 'Content-Type': contentType });
      res.end(content);
    }
  });
});

// 启动服务器
server.listen(PORT, () => {
  console.log(`🚀 大疆无人机监控系统已启动`);
  console.log(`📱 访问地址: http://localhost:${PORT}`);
  console.log(`💻 系统信息: ${os.platform()} ${os.arch()}`);
  console.log(`⏰ 启动时间: ${new Date().toLocaleString()}`);
  console.log(`\n按 Ctrl+C 停止服务器\n`);
  
  // 自动打开浏览器
  const url = `http://localhost:${PORT}`;
  let command;
  
  switch (os.platform()) {
    case 'win32':
      command = `start ${url}`;
      break;
    case 'darwin':
      command = `open ${url}`;
      break;
    default:
      command = `xdg-open ${url}`;
  }
  
  exec(command, (error) => {
    if (error) {
      console.log(`请手动打开浏览器访问: ${url}`);
    }
  });
});

// 优雅关闭
process.on('SIGINT', () => {
  console.log('\n👋 正在关闭服务器...');
  server.close(() => {
    process.exit(0);
  });
});

process.on('SIGTERM', () => {
  console.log('\n👋 正在关闭服务器...');
  server.close(() => {
    process.exit(0);
  });
});
