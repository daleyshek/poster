# POSTER 小程序分享海报

## GO建纯后端绘制的小程序海报图片

生成带二维码的分享海报图片，常见的解决方式是在前端页面进行canvas绘制并获取快照图片，后端生成的案例寥寥无几。  
`poster`就是后端生成海报图片的一种简单实现，其适用于微信小程序海报分享，提示用户长按保存到系统相册后可以转发给朋友、分享到朋友圈。用户点击朋友圈图片，长按识别二维码后就可以跳转到小程序中。

### 海报案例

生成的海报内容包括主图、二维码图、文字内容和标题。  


![海报](https://s1.ax1x.com/2020/04/16/JkS2pn.jpg)

## 开始

### 准备

下载build里面的可执行文件和相关字体  
目录文件,除了poster二进制程序外还需要font字体文件  

```bash
- poster
- resources/font.ttc
```

### 运行

```bash
# 首次需要增加执行权 chmod +x ./poster
# 首次运行会提示填写配置，填写完后再次运行
./poster
# 调试ok后在后台运行
# nohup ./poster &
```

## 配置文件

```json
{
    "mapp": {
        "app_id": "小程序app_id，可选，除非命令行测试运行，否则请用自己的业务系统维护的access_token传值替代",
        "app_secret": "小程序app_secret，可选，除非命令行测试运行，否则请用自己的业务系统维护的access_token传值替代"
    },
    "output_dir": "海报保存目录，可选，默认：./output/",
    "listen_port": "http服务监听端口，可选，默认：2020",
    "fontfile_path": "字体路径，可选，默认: ./resources/font.ttc"
}
```

### HTTP调用

```http
POST / HTTP/1.1
Host: 127.0.0.1:2020
Content-Type: application/x-www-form-urlencoded

access_token=32_fsmZvI40T32ZFTagB4h264eeK1Xgk-XU3Z7DeWTXj4KB24eJkw4Tq3olPVI0Y-s0rBBb8WIVKFVCN5_t7_NXaj5FBq9SDvzHTbrj4_xczKGq3QOByxankwjTbmskDsSG8BPJkReAC0csBu7RKTIaAAAUDS&
scene=10&
page=pages/index/index&
width=640&
title=标题文字&
content=内容文字&
image_url=http://img.mcdsh.com/storage/images/800/O5xEIoExzNk97bRLhaIe0izqo3XbnXKi6j9BWPQb.jpeg&
border_color=#ffffff
```
**参数**  

`access_token` 是微信请求微信接口的access_token，需要在自己的业务系统中保持唯一，避免单独使用appid 和 appsecret生成，生成小程序码需要使用到这个参数。  
`scene`、`page`是小程序获取二维码的参数
`width` 是海报图片的宽度，string类型，太大会失真
`title` 是海报底部主标题，不易过长
`content` 是文字描述，最多可显示3行文字  
`image_url` 是海报主图的URL链接，也可以是主机本地文件绝对路径  
`border_color` 边框和背景色，换其他颜色二维码会有些锯齿

**正常返回**  



```json
{
  "error": 0,
  "result": {
    "poster": "DSGSrZHwzOwFrr6V.jpg",
    "qrcode": "d7VWwZQRprO56Zun.jpg"
  }
}
```
输出目录中会生成海报图片和小程序码图片，将目录暴露到web目录中即可访问到（比如 `ln -s /your/output/dir /var/www/html/poster`）。

HTTP请求也支持JSON格式  

```http
POST / HTTP/1.1
Host: 127.0.0.1:2020
Content-Type: application/json

{
	"access_token":"32_lpIVQGG-qAdXpaUn__X3W47MmWwm-8fIKDmreYnm77ON0E76BmtZXsoZpoIL0_HOre_D6Gr53s4Wv_n41n7mrmUAA7hm-SROYREM0xhX6oCtSds2QoYwv5wtqGh4uRIzCJPJ-ka8BWaJr_l6IYSdADAQGA",
	"scene":"10",
	"page":"pages/index/index",
	"width":"640",
	"title":"标题文字",
	"content":"内容文字",
	"image_url":"http://img.mcdsh.com/storage/images/800/O5xEIoExzNk97bRLhaIe0izqo3XbnXKi6j9BWPQb.jpeg",
    "border_color":"#ffffff"
}
```


### CLI命令行测试运行

可以直接使用命令行方式跑一遍看下效果，此方式需要填写配置文件中的mapp.app_id 和 mapp.app_secret，会挤掉自己业务系统中的access_token，仅供测试。

```bash
./poster -cli \
-imageURL=http://img.mcdsh.com/storage/images/800/O5xEIoExzNk97bRLhaIe0izqo3XbnXKi6j9BWPQb.jpeg \
-title=BestFriendsChina \
-content=我叫小爱，今年一岁啦，领养代替买卖，赶紧扫码带我回家吧！ \
-scene=15625 \
-page=pages/pets/petDetail/petDetail
```

**输出结果**  

```bash
运行成功
小程序码文件 Rw695X5PjQypzIxO.jpg
海报文件 4oikqQLBVrZHBanX.jpg
```


## 总结

这次更新取消强制填写appid和appsecret参数，改为由用户自己实现access_token的获取并传值，解决小程序生成二维码时，access_token容易被自己的业务系统挤掉导致失败的问题。  
移除了复杂的调用方式，使用最简单的http表单数据格式和json格式发送请求，简化到开箱即用。  
