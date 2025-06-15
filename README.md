# lsx

> 该项目是基于 [pls](https://github.com/chenjiandongx/pls) 项目进行二次开发的版本，在此感谢原作者 chenjiandongx 的贡献。
由于该项目已经多年没有更新，linux-command项目已经更新了许多命令，每次都需要拉取最新命令列表，再打包比较麻烦。故在大佬的基础上，开发lsx项目，优化一些逻辑和数据结构，以适配最新的linux-command项目。

## 特性

- 支持通过关键字搜索 Linux 命令
- 提供详细的命令使用说明
- 支持更新内置的 Linux 命令数据库
- 支持多平台适配（包括 ARM 和 AMD 架构）
- 支持输出结果管道传输（如 `| less`）
- 自动化构建和发布流程
- 新增 useShow 参数，默认 false，查看命令为：如 `lsx ls`。当 true 时：`lsx show ls`
- 新增 useLess 参数，默认 false。开启后自动以 less 分页方式输出结果
- 支持通过配置文件自定义数据目录、远程数据源等
- 支持命令自动补全（completion）

## 安装

### 1) 使用编译好的二进制版本

https://github.com/CJSen/lsx/releases

### 2) 自行构建

请参考下方“构建流程”部分。

## 使用方法
1) 下载并解压到任意目录下，赋予可执行文件权限，重命名为 `lsx`

2) 复制`lsx`到`/usr/local/bin`等可直接查找使用的环境目录下，方便全局使用

3) 相关命令使用方法

```shell
~ lsx -h
Impressive Linux commands cheat sheet cli.

Usage:
  lsx [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  search      Search command by keywords
  show        Show the specified command usage.
  upcommands  Update the embedded linux-command.json to the latest version.
  upgrade     Upgrade all commands from remote.
  version     Prints the version about lsx

Flags:
  -h, --help   help for lsx

Use "lsx [command] --help" for more information about a command.
```

建议第一次使用的时候先初始化所有命令(需联网)
```shell
lsx upcommands && lsx upgrade
```

### 命令用法模式说明

- 默认（useShow: false）：
  - 直接输入命令名即可查看用法，如：
    ```shell
    lsx ls
    lsx grep
    ```
- 兼容模式（useShow: true）：
  - 需通过 show 子命令查看，如：
    ```shell
    lsx show ls
    lsx show grep
    ```

### 配置文件说明

lsx 支持通过环境变量 `LSXCONFIG` 指定 YAML 配置文件，覆盖默认配置。

**配置项说明：**
- `dataDir`：本地数据存储目录，默认为 `~/.lsx`（自定义时请写完整路径，不要使用 ~，结尾不要有 /）
- `remoteBaseUrl`：远程命令数据源地址，默认为 `https://unpkg.com/linux-command@latest`
- `useShow`：是否启用 show 子命令模式，布尔值，默认为 false
- `useLess`：是否自动以 less 分页方式输出结果，布尔值，默认为 false

**配置文件示例：**
```yaml
dataDir: "/Users/yourname/.lsx"
remoteBaseUrl: "https://unpkg.com/linux-command@latest"
useShow: false
useLess: true
```

**指定配置文件方法：**
```shell
export LSXCONFIG=/path/to/lsx.yaml
```

### 数据目录说明

- 默认数据目录为 `~/.lsx`，可通过配置文件自定义。
- 所有命令数据和缓存均存放于该目录。

### 命令自动补全

lsx 支持生成自动补全脚本：
```shell
lsx completion bash   # 生成 bash 补全脚本
lsx completion zsh    # 生成 zsh 补全脚本
lsx completion fish   # 生成 fish 补全脚本
```

### 支持 less 分页

- 若 `useLess: true`，则命令输出自动分页。
- 也可手动管道：
  ```shell
  lsx show curl | less
  ```

### 常见问题与错误提示

- 若提示 `[sorry]: could not found command <xxx>`，请确认命令拼写或先执行 `lsx upcommands && lsx upgrade`。
- 若提示 `[error]: failed to parse command index`，请检查数据文件是否损坏，可尝试重新初始化。
- 若自定义数据目录无效，请确保配置文件路径正确，且目录有写权限。

### 多平台二进制文件说明

- 构建产物命名规则：`lsx_{os}_{arch}[.exe]`
- 其中 `{os}` 为操作系统（如 linux、darwin、windows），`{arch}` 为架构（如 amd64、arm64）

## 构建流程

项目使用自动化构建脚本 `build.sh`，支持以下特性：

- 多平台构建（Linux、macOS、Windows）
- 多架构支持（amd64 和 arm64）
- 自动下载最新的 linux-command 数据
- 自动化版本号管理

构建命令：
```shell
./build.sh
```

构建产物会生成在当前目录下，包含以下文件：
- lsx_{os}_{arch}[.exe]
- 其中 {os} 是操作系统，{arch} 是架构

## 开发贡献

该项目欢迎贡献和改进。主要改进方向包括：

- 扩展更多 Linux 命令
- 改进用户界面和交互体验
- 优化性能和稳定性
- 支持更多平台和架构

### 开发环境依赖

- Go 1.18 及以上
- 依赖包：
  - github.com/spf13/cobra
  - gopkg.in/yaml.v3
  - github.com/olekukonko/tablewriter
  - github.com/fatih/color
  - github.com/MichaelMure/go-term-markdown

可通过 `go mod tidy` 自动安装依赖。

对于重大改动，建议先创建 issue 讨论。对于小的改进可以直接提交 PR。

## 致谢

特别感谢jaywcjlove开发的[linux-command](https://github.com/jaywcjlove/linux-command) 为本项目的开发提供了基础。
特别感谢原作者 chenjiandongx 开发了 [pls](https://github.com/chenjiandongx/pls) 项目。

## LICENSE

MIT [@CJSen](https://github.com/CJSen)