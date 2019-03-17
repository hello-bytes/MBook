# MBook - 简单方便的电子书

## 1 功能

MBook一个基于本地`markdown文件`的电子书阅读器，目前实现了静态站点阅读功能（不需要预处理以生成html）。

### 1.1 关于电子书结构

为了方便`MBook`解析，生成，渲染电子书，电子书目录具有以下特点：

- 一般目录下包含3个文件，分别是`SUMMARY.md`，`book.json`,`README.md`。其中`SUMMARY.md`必须存在。
- `SUMMARY.md`，用来生成电子书的目录结构，`SUMMARY.md`的格式与其它电子书系统(如GitBook)大体一样。
- `book.json`，用来标识电子书相关的配置，如果全部是默认配置，则此目录为非必须的，具体配置内容可以参考`book.json`的说明。
- `README.md`，用来做为电子书的首页，可以没有，如果没有，则首页右侧为空白。

### 1.2 SUMMARY.md 说明

这个文件相对于是一本书的目录结构。与`gitbook`大体一致。以下是一个简单的标例 :

```
* [介绍](README.md)
* [第一章：世界，你好](chapte1/hello.md)
   * [认识这个世界](chapte1/meet-the-world.md)
   * [试着喜欢他](chapte1/try-to-loveit.md)
* [第二章：适应与改变](chapte2/follow-and-changes.md)
   * [适应](chapte2/follow.md)
   * [改变](chapte2/changes.md)
* [结束](end/README.md)
```

### 1.3 book.json 说明

略


## 2 代码说明

工程文件主要分布在以下目录：`src`,`resource`。

### 2.1 `resource`目录

存放Go代码无关的资源，模板，配置等，其中下有三个目录：`config`，`public`，`template`,其中分工如下：

- `config`用于存放服务器的配置文件，可以配置的内容包括：端口号，主题，电子书路径等，具体的配置请参考第4章。
- `public`用于存放对外暴露的静态资源，此目录下的文件可以建立CDN以加快访问速度。
- `template`用于存放`golang`的模板文件，不对外暴露。

### 2.2 `src`目录

src目录存放所有的golang代码。

### 2.3 依赖其它的开源代码库

 - [GoMd](https://github.com/hello-bytes/GoMD)
 - [goconfig](https://github.com/hello-bytes/goconfig)

## 4 配置文件

配置文件目录为：`/resource/config/env.config`，模块文件参考:`/resource/config/sample.env.config`

配置的字段包括`port`,`theme`,`books`等；配置的节包括：`WebSite`，`CodeBook`及各个电子书的节。

### 4.1 WebSite 节的相关配置

WebSite 节下主要有以下Key的配置：

- `port`：网站监听的端口号。
- `theme`：网站的主题，通过指定主题，可以得到完全不同的电子书呈现样式，目前内置`blog`与`gitbook`这2个样式。
- `cdn`：用于配置cdn的前缀，这样，在`golang html template`中就可以使用`cdnAssets`这个方法了

### 4.2 MBook 节的相关配置

电子书配置的根节点在`CodeBook`下，`CodeBook`节通过`books`这个key来指定所有的电子书。`books`的值为所有电子书的id列表，以`;`进行分隔，示例如下：

```
[MBook]
books=HelloWorld;GolangCookBook;
```

如上所示，HelloWorld与GolangCookBook为2本电子书的`BookId`，具体每一本书的详细信息，则在各自的节中配置，各自的节的节点名为`CodeBook:` + {BookId}，以`HelloWorld`这本书来说，具体的配置节就是`CodeBook:HelloWorld`，示例如下：

```
[MBook:HelloWorld]
name=HelloWorld，程序员的第一堂课
path=/Users/Aaron/Articles/CodeBook/Echo
```

### 4.3 Ini示例

如上所述，一个典型的ini配置如下所示:

```
[WebSite]
port=8080
theme=blog
cdn=https://codingsky.oss-cn-hangzhou.aliyuncs.com/hellotech

[CodeBook]
books=hellobytes;

[CodeBook:hellobytes]
name=blog
path=/Users/Aaron/Project/Docuement/Blog
```

## 5 后台地址

MBook提供了后台管理功能，几个比较重要的页面分别为：

- **/backend/initenv**：初始部署时，初始化的地址，会同时生成`Security Key`。
- **/backend/login**：后台登录页面，需要使用 `Security Key`。
- **/backend/index**：后台管理主页。

> **请注意：**`Security Key`的相对目录为`resource/config/sk.config`。sk.config文件在部署后第一次被访问时生成，此进系统会跳到`/backend/initenv`进行告知，请妥善保管。每次登录后台都需要使用这个`Security Key`。

## 6 开发，运行，调试

如果您要在开发环境下运行代码，可以试着按以下步骤操作：

- **1.**下载本地，假设目录为`/Users/Aaron/go/src/github.com/hello-bytes/MBook`，后面举例也以此目录进行举例
- **2.**安装官方的Golang依赖管理模块：`dep`
- **3.**建立`vender`目录：`cd /Users/Aaron/go/src/github.com/hello-bytes/MBook/src/ && dep init`。
- **5.**初始化配置文件：`cd /Users/Aaron/go/src/github.com/hello-bytes/MBook/resource && cp sample.env.config env.config`
- **5.**运行根目录下的`run.sh`：`cd /Users/Aaron/go/src/github.com/hello-bytes/MBook/ && bash run.sh`

如果一切顺利，现在应该可以在浏览器中访问了：[http://127.0.0.1:8080](http://127.0.0.1:8080)。