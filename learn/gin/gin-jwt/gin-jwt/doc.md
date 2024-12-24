```shell
├── client.py  # 客户端测试脚本
├── controllers  # 控制器相关包
│   └── auth.go  # 控制器方法实现
├── gin-jwt.bin  # 编译的二进制文件
├── go.mod  # go 项目文件
├── go.sum  # go 项目文件
├── main.go  # 程序入口文件
├── middlewares  # 中间件相关包
│   └── middlewares.go  # 中间件代码文件
├── models  # 存储层相关包
│   ├── setup.go  # 配置数据库连接
│   └── user.go  # user模块相关数据交互的代码文件
├── README.md  # git repo的描述文件
└── utils  # 工具类包
    └── token  # token相关工具类包
        └── token.go  # token工具的代码文件
```