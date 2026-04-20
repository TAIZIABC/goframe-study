---
layout: home

hero:
  name: "Go Web 编程"
  text: "用 Go 构建 Web 应用"
  tagline: 从入门到精通，由 astaxie（谢孟军）编写的经典 Go Web 开发教程
  image:
    src: /logo.svg
    alt: Go Web 编程
  actions:
    - theme: brand
      text: 开始阅读 →
      link: /01/0
    - theme: alt
      text: GitHub
      link: https://github.com/astaxie/build-web-application-with-golang

features:
  - icon: 🔧
    title: Go 环境配置
    details: 从零开始搭建 Go 开发环境，掌握 GOPATH、Go Module 和常用命令。
    link: /01/0
  - icon: 📝
    title: Go 语言基础
    details: 深入理解 Go 语法、struct、interface、并发编程等核心概念。
    link: /02/0
  - icon: 🌐
    title: Web 基础
    details: 学习 HTTP 协议原理，用 Go 搭建 Web 服务器，理解 net/http 包。
    link: /03/0
  - icon: 📋
    title: 表单处理
    details: 处理表单输入、验证、文件上传，防止跨站脚本和重复提交。
    link: /04/0
  - icon: 🗄️
    title: 数据库操作
    details: 使用 MySQL、PostgreSQL、SQLite、Redis、MongoDB 等数据库。
    link: /05/0
  - icon: 🔐
    title: 安全与加密
    details: 防范 CSRF、XSS、SQL 注入等攻击，安全存储密码和加密数据。
    link: /09/0
  - icon: 🏗️
    title: Web 框架设计
    details: 从零设计一个 Web 框架，理解路由、控制器、中间件的设计思想。
    link: /13/0
  - icon: 🚀
    title: 部署与维护
    details: 应用日志、错误处理、部署方案、备份恢复等生产环境实践。
    link: /12/0
---

<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: -webkit-linear-gradient(120deg, #00ADD8 30%, #5DC9E2);
  --vp-home-hero-image-background-image: linear-gradient(-45deg, #00ADD8aa 50%, #5DC9E266 50%);
  --vp-home-hero-image-filter: blur(44px);
}

@media (min-width: 640px) {
  :root {
    --vp-home-hero-image-filter: blur(56px);
  }
}

@media (min-width: 960px) {
  :root {
    --vp-home-hero-image-filter: blur(68px);
  }
}
</style>
