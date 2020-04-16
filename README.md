# POSTER 小程序分享海报

## GO建纯后端绘制的小程序海报图片

生成带二维码的分享海报图片，常见的解决方式是在前端页面进行canvas绘制并获取快照图片，后端生成的案例寥寥无几。poster就是后端生成海报图片的一种简单实现，其适用于微信小程序海报分享，保存到相册后可以分享到朋友圈。用户点击朋友圈图片，识别二维码后就可以跳转到小程序中。简直是一大利器。

### 特点

封装了小程序码的获取功能，生成海报一气呵成，无需繁琐的调用。

### 示例

生成的海报内容包括主图、二维码图、文字内容和标题

![海报](https://s1.ax1x.com/2020/04/16/JkS2pn.jpg)

***当前分支仅支持命令行调用和http调用***

## 运行环境

- 字体文件（必须，build/resources目录中有个测试字体）

## 调用方式

- 命令行调用，用于快速测试
- JSON-RPC HTTP调用(端口2019，支持http直接请求，examples中也有示例)

build目录中的二进制可执行文件poster为JSON-RPC服务，和resoures目录放在一起即可运行。使用`./poster -h` 可以查看命令行测试时需要的参数。

## RPC参数说明

```bash
# 下面是生成小程序二维码的调用参数
# Scene 小程序页面的场景值，一般为页面的变量ID
Scene
# Page 小程序的页面，详情见小程序文档
Page
# Width 生成小程序二维码的尺寸，可选，默认480
Width

# 下面是生成海报的调用参数
# Title 海报的标题，显示在底部
Title
# Content 海报文字内容，显示会限制在3行以内
Content
# ImageURL 主图资源，可以为http资源，也可以是本地文件资源，图片会拉伸至完全覆盖中心，格式常见的jpg、png都行
ImageURL
# QRCodeURL 二维码图片资源，可以为http资源，也可以为本地文件资源，这个参数是否必要要看调用的哪个RPC方法
QRCodeURL

```

## 运行结果

生成小程序二维码调用成功后会将二维码图片保存在输出目录中（默认`output/`），返回文件名称。
生成海报调用成功后会在输出目录生成一张海报图片，并返回该图片的文件名称。

## 开始

首先准备必要的程序和资源，主要是resources/font.ttc字体文件，如果要使用小程序码，还需要准备好小程序的参数 `appid` 和 `secret`（填写在配置文件poster.json中的MPAppID和MPAppSecret）。

```bash
$ ls
poster resources

# 查看命令行方式运行帮助
# 如果'提示配置文件已生成，请填写完成后再次运行',可以在生成的poster.json中填写小程序的参数MPAppID和MPAppSecret，然后再次运行
$ ./poster -h
配置文件已生成，请填写完成后再次运行
```

配置文件poser.json，OutputDIR为空时默认输出目录为`output/`,自定义输出目录时需要在结尾加上`/`防止报错！

```json
{
    "MPAppID": "wx02************",
    "MPAppSecret": "1ba28*************d061b",
    "OutputDIR": ""
}
```

### CLI测试调用

填写好小程序的参数后，就可以通过命令行方式测试运行了

```bash
./poster -cli \
-imageURL=http://img.mcdsh.com/storage/images/800/O5xEIoExzNk97bRLhaIe0izqo3XbnXKi6j9BWPQb.jpeg \
-title=BestFriendsChina \
-content=我叫小爱，今年一岁啦，领养代替买卖，赶紧扫码带我回家吧！ \
-scene=15625 \
-page=pages/pets/petDetail/petDetail
```

如果参数都没有问题的话，输出结果：

```bash
运行成功
小程序码文件 Rw695X5PjQypzIxO.jpg
海报文件 4oikqQLBVrZHBanX.jpg
```

使用微信扫描海报或者小程序码后，会进入到小程序详情页中，并携带scene参数，达到我们的期望效果。

### JSON-RPC 调用

运行基于HTTP的JSON-RPC服务

```bash
# 通过nohup调用
# nohup ./poster >> log.txt &
# 直接运行
$ ./poster
```

使用HTTP进行调用Poster.Both方法，一次完成小程序码和海报图片的生成。Both方法不需要传递qrcodeURL参数。还有两个Both的子方法，一个是单独生成小程序码，一个是单独生成海报，这里就不再描述了。

```bash
$ curl \
-s \
-H "Content-Type: application/json" \
-X POST \
-d '{"jsonrpc":"2.0","method":"Poster.Both","params":[{"accessToken":"小程序服务接口的access_token","scene":"15625","page":"pages/pets/petDetail/petDetail","title":"BestFriendsChina","content":"我叫小爱，今年一岁啦，领养代替买卖，赶紧扫码带我回家吧！","imageURL":"http://img.mcdsh.com/storage/images/800/O5xEIoExzNk97bRLhaIe0izqo3XbnXKi6j9BWPQb.jpeg"}],"id":1}' \
http://127.0.0.1:2019
```

响应

```bash
# 生成的小程序码和海报图片将会在output目录中
{"id":1,"result":{"QRCodeName":"lWFTlBtg1cjGmjLC.jpg","PosterName":"r9cZ0nudX7SGSj4y.jpg"},"error":null}
```

二维码和海报都生成成功

## 总结

使用poster生成小程序的分享海报非常简单，通过cli测试成功后，就可以使用后端程序进行http调用。通过修改web服务器配置将output目录暴露出来就可以直接访问到海报图片和小程序码图片。