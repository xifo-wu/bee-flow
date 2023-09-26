
```
______     ______     ______     ______   __         ______     __     __
/\  == \   /\  ___\   /\  ___\   /\  ___\ /\ \       /\  __ \   /\ \  _ \ \
\ \  __<   \ \  __\   \ \  __\   \ \  __\ \ \ \____  \ \ \/\ \  \ \ \/ ".\ \
 \ \_____\  \ \_____\  \ \_____\  \ \_\    \ \_____\  \ \_____\  \ \__/".~\_\
  \/_____/   \/_____/   \/_____/   \/_/     \/_____/   \/_____/   \/_/   \/_/
```


## TODO

- [ ] 支持 rename 命令。重命名文件夹内文件
- [ ] 只修改集数重命名方式

## 分类
程序通过 qbittorrent 分类来判断需要执行什么操作

1. BeeFlow
不包含 BeeFlow 的分类程序将不会执行任何操作

## RSS 订阅
根据路径进行管理 RSS 订阅

### 重命名模式
重命名方式安装 Emby 建议格式重命名，[官方文档地址](https://emby.media/support/articles/TV-Naming.html)

- ```0```: 不进行重命名
- 其他模式（懒得起名字）：
  - 模式 ```1``` : 日常 S01E01.mp4
  - 模式 ```2``` : 日常 - S01E01 - Baha.1080P.繁体内嵌.mp4

### 使用方法

```bash
beeflow rss add "https://mikanani.me/RSS/Bangumi?bangumiId=3079&subgroupid=21" \
--mode 0 \
--name "文豪野犬" \
--year 2016 \
--season 5 \
--savePath "文豪野犬 (2016)/Season 05" \
--path "/root/backup/CloudDrive/阿里云盘Open/Temp/文豪野犬 (2016)/Season 05" \
--resolution "1080P" \
--subtitle "简繁内封" \
--group "LoliHouse" \
--offset 0

```
