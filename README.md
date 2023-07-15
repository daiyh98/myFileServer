## 服务架构

![image](https://github.com/daiyh98/myFileServer/assets/43029102/d713ecab-7deb-46e4-83e5-c7266a4439b3)

### 用户与server之间的交互

1. 上传/下载
2. 查看/删除

### 接口列表

| 接口描述 | 接口url            |
| -------- | ------------------ |
| 文件上传 | POST /file/upload  |
| 文件下载 | GET /file/download |
| 文件查询 | GET /file/query    |
| 文件删除 | POST /file/update  |
|          |                    |

### 文件上传功能原理：

![go_server_structure_annotated_interface](https://github.com/daiyh98/myFileServer/assets/43029102/f52b0344-430b-4fe5-8ebe-0c0d9e3c5153)

1. 获取上传页面
2. 选取本地文件，以form形式上传文件
3. 云端接收文件流，写入本地存储
4. 云端更新文件元信息集合

### 文件元信息记录、更新与查询

首先创建了一个metaInfo包，包内部定义了一个`FileMeta`结构体，用来存储每个文件的元信息，包括：

- 文件哈希值
- 文件名
- 文件位置
- 更新时间
- 文件大小

然后为了存储所有元信息，采用`var`语句声明出一个map变量，key是文件哈希值，value是该文件的元信息结构体。在该包的`init()`函数中为map变量进行了`make`函数初始化，这样就不会发生对nil map进行更新而发生的panic。
