import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Go Web 编程',
  description: '一本用 Go 语言编写 Web 应用的开源书籍',
  lang: 'zh-CN',
  lastUpdated: true,

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],
  ],

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'Go Web 编程',

    nav: [
      { text: '首页', link: '/' },
      { text: '开始阅读', link: '/01/0' },
      { text: '参考资料', link: '/ref' },
    ],

    sidebar: [
      {
        text: '第1章 Go环境配置',
        collapsed: true,
        items: [
          { text: '概述', link: '/01/0' },
          { text: '1.1 Go安装', link: '/01/1' },
          { text: '1.2 GOPATH 与工作空间', link: '/01/2' },
          { text: '1.3 Go 命令', link: '/01/3' },
          { text: '1.4 Go开发工具', link: '/01/4' },
          { text: '1.5 小结', link: '/01/5' },
        ]
      },
      {
        text: '第2章 Go语言基础',
        collapsed: true,
        items: [
          { text: '概述', link: '/02/0' },
          { text: '2.1 你好，Go', link: '/02/1' },
          { text: '2.2 Go基础', link: '/02/2' },
          { text: '2.3 流程和函数', link: '/02/3' },
          { text: '2.4 struct', link: '/02/4' },
          { text: '2.5 面向对象', link: '/02/5' },
          { text: '2.6 interface', link: '/02/6' },
          { text: '2.7 并发', link: '/02/7' },
          { text: '2.8 小结', link: '/02/8' },
        ]
      },
      {
        text: '第3章 Web基础',
        collapsed: true,
        items: [
          { text: '概述', link: '/03/0' },
          { text: '3.1 Web工作方式', link: '/03/1' },
          { text: '3.2 Go搭建一个简单的Web服务', link: '/03/2' },
          { text: '3.3 Go如何使得Web工作', link: '/03/3' },
          { text: '3.4 Go的http包详解', link: '/03/4' },
          { text: '3.5 小结', link: '/03/5' },
        ]
      },
      {
        text: '第4章 表单',
        collapsed: true,
        items: [
          { text: '概述', link: '/04/0' },
          { text: '4.1 处理表单的输入', link: '/04/1' },
          { text: '4.2 验证表单的输入', link: '/04/2' },
          { text: '4.3 预防跨站脚本', link: '/04/3' },
          { text: '4.4 防止多次递交表单', link: '/04/4' },
          { text: '4.5 处理文件上传', link: '/04/5' },
          { text: '4.6 小结', link: '/04/6' },
        ]
      },
      {
        text: '第5章 访问数据库',
        collapsed: true,
        items: [
          { text: '概述', link: '/05/0' },
          { text: '5.1 database/sql接口', link: '/05/1' },
          { text: '5.2 使用MySQL数据库', link: '/05/2' },
          { text: '5.3 使用SQLite数据库', link: '/05/3' },
          { text: '5.4 使用PostgreSQL数据库', link: '/05/4' },
          { text: '5.5 使用beedb库进行ORM开发', link: '/05/5' },
          { text: '5.6 NOSQL数据库操作', link: '/05/6' },
          { text: '5.7 小结', link: '/05/7' },
        ]
      },
      {
        text: '第6章 Session和数据存储',
        collapsed: true,
        items: [
          { text: '概述', link: '/06/0' },
          { text: '6.1 Session和Cookie', link: '/06/1' },
          { text: '6.2 Go如何使用Session', link: '/06/2' },
          { text: '6.3 Session存储', link: '/06/3' },
          { text: '6.4 预防Session劫持', link: '/06/4' },
          { text: '6.5 小结', link: '/06/5' },
        ]
      },
      {
        text: '第7章 文本文件处理',
        collapsed: true,
        items: [
          { text: '概述', link: '/07/0' },
          { text: '7.1 XML处理', link: '/07/1' },
          { text: '7.2 JSON处理', link: '/07/2' },
          { text: '7.3 正则处理', link: '/07/3' },
          { text: '7.4 模板处理', link: '/07/4' },
          { text: '7.5 文件操作', link: '/07/5' },
          { text: '7.6 字符串处理', link: '/07/6' },
          { text: '7.7 小结', link: '/07/7' },
        ]
      },
      {
        text: '第8章 Web服务',
        collapsed: true,
        items: [
          { text: '概述', link: '/08/0' },
          { text: '8.1 Socket编程', link: '/08/1' },
          { text: '8.2 WebSocket', link: '/08/2' },
          { text: '8.3 REST', link: '/08/3' },
          { text: '8.4 RPC', link: '/08/4' },
          { text: '8.5 小结', link: '/08/5' },
        ]
      },
      {
        text: '第9章 安全与加密',
        collapsed: true,
        items: [
          { text: '概述', link: '/09/0' },
          { text: '9.1 预防CSRF攻击', link: '/09/1' },
          { text: '9.2 确保输入过滤', link: '/09/2' },
          { text: '9.3 避免XSS攻击', link: '/09/3' },
          { text: '9.4 避免SQL注入', link: '/09/4' },
          { text: '9.5 存储密码', link: '/09/5' },
          { text: '9.6 加密和解密数据', link: '/09/6' },
          { text: '9.7 小结', link: '/09/7' },
        ]
      },
      {
        text: '第10章 国际化和本地化',
        collapsed: true,
        items: [
          { text: '概述', link: '/10/0' },
          { text: '10.1 设置默认地区', link: '/10/1' },
          { text: '10.2 本地化资源', link: '/10/2' },
          { text: '10.3 国际化站点', link: '/10/3' },
          { text: '10.4 小结', link: '/10/4' },
        ]
      },
      {
        text: '第11章 错误处理、调试和测试',
        collapsed: true,
        items: [
          { text: '概述', link: '/11/0' },
          { text: '11.1 错误处理', link: '/11/1' },
          { text: '11.2 使用GDB调试', link: '/11/2' },
          { text: '11.3 Go怎么写测试用例', link: '/11/3' },
          { text: '11.4 小结', link: '/11/4' },
        ]
      },
      {
        text: '第12章 部署与维护',
        collapsed: true,
        items: [
          { text: '概述', link: '/12/0' },
          { text: '12.1 应用日志', link: '/12/1' },
          { text: '12.2 网站错误处理', link: '/12/2' },
          { text: '12.3 应用部署', link: '/12/3' },
          { text: '12.4 备份和恢复', link: '/12/4' },
          { text: '12.5 小结', link: '/12/5' },
        ]
      },
      {
        text: '第13章 如何设计一个Web框架',
        collapsed: true,
        items: [
          { text: '概述', link: '/13/0' },
          { text: '13.1 项目规划', link: '/13/1' },
          { text: '13.2 自定义路由器设计', link: '/13/2' },
          { text: '13.3 Controller设计', link: '/13/3' },
          { text: '13.4 日志和配置设计', link: '/13/4' },
          { text: '13.5 实现博客的增删改', link: '/13/5' },
          { text: '13.6 小结', link: '/13/6' },
        ]
      },
      {
        text: '第14章 扩展Web框架',
        collapsed: true,
        items: [
          { text: '概述', link: '/14/0' },
          { text: '14.1 静态文件支持', link: '/14/1' },
          { text: '14.2 Session支持', link: '/14/2' },
          { text: '14.3 表单支持', link: '/14/3' },
          { text: '14.4 用户认证', link: '/14/4' },
          { text: '14.5 多语言支持', link: '/14/5' },
          { text: '14.6 pprof支持', link: '/14/6' },
          { text: '14.7 小结', link: '/14/7' },
        ]
      },
    ],

    outline: {
      level: [2, 3],
      label: '页面导航',
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/astaxie/build-web-application-with-golang' }
    ],

    footer: {
      message: '基于 BSD 许可发布',
      copyright: '原作者 astaxie（谢孟军） | VitePress 版本优化'
    },

    search: {
      provider: 'local',
      options: {
        translations: {
          button: { buttonText: '搜索文档', buttonAriaLabel: '搜索文档' },
          modal: {
            noResultsText: '无法找到相关结果',
            resetButtonTitle: '清除查询条件',
            footer: { selectText: '选择', navigateText: '切换', closeText: '关闭' }
          }
        }
      }
    },

    docFooter: {
      prev: '上一页',
      next: '下一页'
    },

    lastUpdated: {
      text: '最后更新于',
    },

    returnToTopLabel: '回到顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式',
  }
})
