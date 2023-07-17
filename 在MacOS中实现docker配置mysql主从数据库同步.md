# 在MacOS中实现docker配置mysql主从数据库同步

## 一、安装Docker

1. **下载 Docker Desktop for Mac**

   访问 Docker 的官方网站：https://www.docker.com/products/docker-desktop

   从页面上选择 Docker Desktop for Mac 下载选项。下载完成后你会得到一个 `.dmg` 文件。

2. **安装 Docker Desktop**

   双击下载的 `.dmg` 文件，然后将 Docker.app 拖放到你的 Applications 文件夹中。

3. **启动 Docker Desktop**

   打开 Applications 文件夹，然后双击 Docker.app 开始启动 Docker。第一次启动 Docker 时，它将询问你是否想要发送匿名统计信息给 Docker，你可以选择接受或者拒绝。

4. **验证 Docker 安装**

   打开一个终端窗口，然后输入 `docker version`，然后按回车键。如果 Docker 已经正确安装，你应该会看到相关的 Docker 版本信息。

## 二、下载 MySQL Docker 镜像

首先，确保你已经下载了 MySQL Docker 镜像。如果还没有，可以在终端执行以下命令来下载 MySQL 镜像：

```bash
docker pull mysql:8.0
```

> 在大多数情况下，你应该使用 Docker 官方镜像库提供的最新版本，例如 `docker pull mysql:8.0`。除非你有特别的需求，否则没有必要精确匹配 MySQL 的小版本号。

接下来，根据你的需求，我们可以创建一个主节点和一个从节点。

## 三、创建主从节点

### **运行主数据库容器**

```bash
docker run --name master -p 3306:3306 -e MYSQL_ROOT_PASSWORD=yourpassword -d mysql:8.0
```

这个命令会创建一个名为“master”的 Docker 容器，映射容器的 3306 端口到主机的 3306 端口，并设置 root 用户的密码为你选择的密码。

### **运行从数据库容器**

```bash
docker run --name slave -p 3307:3306 -e MYSQL_ROOT_PASSWORD=yourpassword -d mysql:8.0 --server-id=2
```

这个命令会创建一个名为“slave”的 Docker 容器，映射容器的 3306 端口到主机的 3307 端口，并设置 root 用户的密码为你选择的密码。==一定要添加`--server-id=2`，否则由于两个节点默认服务器id都是1，在设置主从同步时会发生报错。==

### 验证创建结果

现在，如果你运行 `sudo docker ps`，你应该可以看到两个运行中的 MySQL 容器，一个是主节点，另一个是从节点。

运行结果演示：

```shel
CONTAINER ID   IMAGE       COMMAND                   CREATED          STATUS          PORTS                               NAMES
a93004ab26a7   mysql:8.0   "docker-entrypoint.s…"   44 minutes ago   Up 44 minutes   33060/tcp, 0.0.0.0:3307->3306/tcp   slave
db04ab686666   mysql:8.0   "docker-entrypoint.s…"   2 hours ago      Up 2 hours      0.0.0.0:3306->3306/tcp, 33060/tcp   master
```



另外，由于 macOS 系统没有提供 `netstat -antup | grep docker` 命令，你可以用以下命令查看所有被占用的端口（包括docker使用的）：

```bash
sudo lsof -nP -iTCP -sTCP:LISTEN
```

运行结果演示：

```she
COMMAND    PID   USER   FD   TYPE            DEVICE SIZE/OFF NODE NAME
ToDesk_Se  651   root    9u  IPv4 0x3c7e31dda1a556b      0t0  TCP 127.0.0.1:35600 (LISTEN)
ECAgent   1336   root    8u  IPv4 0x3c7e31dda1a6b8b      0t0  TCP 127.0.0.1:54530 (LISTEN)
clashr-da 1992  daiyh   12u  IPv6 0x3c7e32772be4253      0t0  TCP *:4780 (LISTEN)
clashr-da 1992  daiyh   13u  IPv4 0x3c7e31dda780cbb      0t0  TCP 127.0.0.1:4788 (LISTEN)
clashr-da 1992  daiyh   14u  IPv6 0x3c7e32772be4a53      0t0  TCP *:4781 (LISTEN)
mysqld    2880 _mysql   18u  IPv6 0x3c7e32772be5a53      0t0  TCP *:33060 (LISTEN)
mysqld    2880 _mysql   20u  IPv6 0x3c7e32772be2a53      0t0  TCP *:3306 (LISTEN)
```

### 删除容器（容错）

当你发现你的容器创建得有问题，可以删除之后重新创建

```bash
docker rm master
```

## 四、设置主从同步

[3-2 MySQL主从数据同步演示_哔哩哔哩_bilibili](https://www.bilibili.com/video/BV1Uh411u7rg?p=10&vd_source=bc5ee05972cc0c66277362a57b9e054c)

### 具体操作

```bash
mysql -u root -h 127.0.0.1 -P 3307 -p
```

```mysql
mysql> CHANGE MASTER TO MASTER_HOST='192.168.2.238',MASTER_USER='reader',MASTER_LOG_FILE='binlog.000002',MASTER_LOG_POS=0,MASTER_PASSWORD='reader';
mysql> start slave;
```

### 细节

#### 如何获取`MASTER_HOST`

1. 打开 Terminal 应用。

2. 输入 `ifconfig` 命令并回车。

   ```bash
   ifconfig
   ```

3. 在输出结果中，找到 `en0`（通常这是有线连接的网卡）或者 `en1`（通常这是无线连接的网卡）部分，`inet` 后面的就是你的 IP 地址。下面是一个例子：

   ```bash
   en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
       inet 192.168.1.2 netmask 0xffffff00 broadcast 192.168.1.255
       nd6 options=201<PERFORMNUD,DAD>
       media: autoselect
       status: active
   ```

在这个例子中，IP 地址就是 `192.168.1.2`。

注意，这是你在本地网络中的 IP 地址，不是你在互联网中的公共 IP 地址。如果你想要获取公共 IP 地址，你可以访问一些提供该服务的网站，如 https://www.whatismyip.com/。

#### 如何获取`MASTER_LOG_FILE`

```shell
mysql -u root -h 127.0.0.1 -p
```

```mysql
mysql> show master status;
```

#### 如何获取`MASTER_USER`以及`MASTER_PASSWORD`

1. 首先，获取你 Docker 容器的 ID。你可以通过 `docker ps` 命令来获取正在运行的 Docker 容器的列表，其中也包括它们的 ID。

2. 然后，使用 `docker exec` 命令连接到你的 Docker 容器。例如，如果你的 Docker 容器 ID 是 `db04ab686666`，你可以使用以下命令连接到该容器：

```bash
docker exec -it db04ab686666 bash
```

`-it` 参数表示你希望在容器内部启动一个交互式的 shell。

3. 一旦你连接到了 Docker 容器，你就可以在容器内部运行 MySQL 命令了。你可以通过以下命令登录到 MySQL：

```bash
mysql -u root -p
```

然后按提示输入 root 用户的密码。

4. 接下来你就可以按照之前的指示来创建新的用户和授权了。

注意，每次你使用 `docker exec` 命令时，你都是在新的 shell 会话中运行命令，而这个会话是独立于你的主机系统的。当你退出这个 shell 会话（例如通过 `exit` 命令）时，你将回到你的主机系统的 shell。

你需要确保在 MySQL 主服务器上已经创建了这个用户，并且赋予了它复制的权限。创建新用户并赋予复制权限的命令如下：

```mysql
mysql> CREATE USER 'myrepluser'@'%' IDENTIFIED BY 'password';（错误）
mysql> CREATE USER 'myrepluser'@'%' IDENTIFIED WITH mysql_native_password BY 'password';（正确）
mysql> GRANT REPLICATION SLAVE ON *.* TO 'myrepluser'@'%';
```

不要使用第一种创建用户方式，这会导致设置主从同步之后出现如下错误：

```shell
Error connecting to source 'reader@192.168.2.238:3306'. This was attempt 1/86400, with a delay of 60 seconds between attempts. Message: Authentication plugin 'caching_sha2_password' reported error: Authentication requires secure connection.
```

可能是由于MySQL 8.0使用的新的默认身份验证插件"caching_sha2_password"。这个新的插件要求通过SSL/TLS加密的连接进行身份验证。

赋予该用户访问和操作数据库的权限。例如，如果你希望该用户能够访问名为 'mydatabase' 的数据库，并对其进行全部操作，你可以使用以下命令：

```mysql
mysql> GRANT ALL PRIVILEGES ON mydatabase.* TO 'myuser'@'localhost';
mysql> FLUSH PRIVILEGES;
```

#### TIPS

如果你创建的用户名是 `'reader'@'localhost'`，那么这个用户只能在本地连接MySQL，不能从远程连接。

`localhost`是一个特殊的主机名，在MySQL中，它指的是本地使用套接字文件进行的连接，而不是使用网络进行的连接。当你从同一台机器上的应用程序（例如，运行在同一台机器上的web服务器）连接到MySQL时，使用 `localhost` 会更有效率。但如果你需要从远程主机连接，你需要将用户的主机部分设置为 `'%'`（代表任何主机）或者具体的IP地址。

你需要修改你的用户来允许从远程主机连接。你可以使用 `ALTER USER` 命令来做到这一点。以下是如何更改 `localhost` 到 `'%'` 的例子：

```sql
mysql> ALTER USER 'reader'@'localhost' RENAME TO 'reader'@'%';
mysql> FLUSH PRIVILEGES;
```

请注意，对于网络连接，使用 `'%'` 允许任何主机都可以连接，这可能存在安全风险。如果可能的话，最好指定一个具体的IP地址或者地址范围。

注意：如果你有防火墙或者网络策略限制了MySQL服务器的访问，那么你可能还需要更新这些设置来允许从远程主机的连接。

## 五、验证主从同步

[3-2 MySQL主从数据同步演示_哔哩哔哩_bilibili](https://www.bilibili.com/video/BV1Uh411u7rg?p=10&vd_source=bc5ee05972cc0c66277362a57b9e054c)

```mysql
start slave;
show slave status\G;
```

检查：

```shell
Slave_IO_Running: Yes
Slave_SQL_Running: Yes
```

