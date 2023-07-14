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

