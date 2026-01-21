# glance-bilibil 需求文档

## 项目背景

glance https://github.com/glanceapp/glance 是一个自托管面板服务，里面可以显示不同的组件；但是官方的videos组件只适用于youtube视频；
我准备开发一个bilibili视频组件，采用extension插件的方式

## 项目目标

按照glance的extension插件开发文档，开发一个bilibili视频组件
输出的extension插件可以被glance加载
插件显示样式参考glance的videos组件，但是视频来源是bilibili
通过Get请求的传参控制显示样式

## 技术栈

go

## 参考资料：

### glance extension插件开发文档

https://github.com/glanceapp/glance/blob/main/docs/extension.md

### glance videos widget文档

./doc/widget-videos.md

### glance videos 插件相关源码 

./doc/widget-videos.go
./doc/videos.html
./doc/videos-grid.html
./doc/videos-vertical-list.html
./doc/video-card-conntents.html

### bilibili api文档

./doc/api.md
./doc/wbi.md
