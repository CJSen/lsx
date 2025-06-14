
# lsx

### Installation

#### 1) 使用 `go get` 安装

```shell
$ go get -u github.com/CJSen/lsx
```

#### 2) 使用编译好的二进制版本

https://github.com/CJSen/lsx/releases

### Usages

```shell
~ 🐶 lsx --help
Impressive Linux commands cheat sheet cli.

Usage:
  lsx [command]

Available Commands:
  help        Help about any command
  search      Search command by keywords
  show        Show the specified command usage.
  upcommands  Update the embedded linux-command.json to the latest version.
  upgrade     Upgrade all commands from remote.
  version     Prints the version of lsx

Flags:
  -h, --help   help for lsx

Use "lsx [command] --help" for more information about a command.
```

建议第一次使用的时候先初始化所有命令
```shell
$ 🐶 lsx upgrade
```

数据文件默认位于 `~/.lsx/xxx` 可以更改环境变量LSXCONFIG的值来调整配置文件，设置数据存放路径

可以将输出结果传入到 less 管道
```shell
$ 🐶 lsx show curl | less
```

效果图

![](https://user-images.githubusercontent.com/19553554/122259619-f1e3f780-cf04-11eb-949e-763d82a4e3b9.png)
![](https://user-images.githubusercontent.com/19553554/122258451-a0873880-cf03-11eb-865f-067416787cb7.png)


### LICENSE

MIT [©chenjiandongx](https://github.com/chenjiandongx)
