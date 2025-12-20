# 中南民族大学计算机学院智慧党建系统

**开发时间**: 2025年8月——2025年11月  
**开发团队**: 新思路实验室Geek组2023级全体

本项目基于Go语言和Gin框架开发的党校管理系统，提供学员信息管理、培训记录管理、党员发展历程管理等功能。系统采用RESTful API设计，支持完整的CRUD操作和JWT认证。

## 🚀 快速体验

**在线测试**: 打开 `api_test_utf8.html` 文件即可在浏览器中测试所有API接口

- 支持所有认证、学员管理、培训管理等接口
- 实时响应显示和错误处理
- 现代化界面，操作简单直观

## ✨ 主要功能

- 🔐 **用户认证与授权** - JWT Token认证，支持多角色权限控制（管理员/教师/学生）
- 👨‍🎓 **学员信息管理** - 完整的学员CRUD操作，支持分页查询、Excel导入导出、照片上传
- 📚 **培训记录管理** - 培训课程记录，成绩管理和证书管理
- 📈 **党员发展历程** - 从群众到正式党员的完整发展轨迹，支持文本格式时间输入
- 🏆 **奖惩记录管理** - 学员奖励和处分记录管理
- 📊 **数据统计分析** - 培训成绩统计和数据分析
- 🏢 **组织关系管理** - 党支部和党小组管理，支持按组织筛选
- 🌐 **RESTful API** - 标准化的API接口设计
- 🛠️ **在线测试工具** - 内置API测试页面，支持实时测试

## 🛠️ 技术栈

- **后端框架**: Gin (Go Web框架)
- **数据库**: MySQL 5.7+
- **ORM**: GORM (Go ORM库)
- **认证**: JWT (JSON Web Token)
- **编程语言**: Go 1.21+
- **前端测试**: HTML5 + JavaScript (API测试页面)

## 📁 项目结构

```
dangjiantest/
├── api/ # API处理器层
│ ├── auth_handler.go # 认证相关接口
│ ├── student_handler.go # 学员管理接口
│ ├── training_handler.go # 培训管理接口
│ └── stats_handler.go # 统计分析接口
├── config/ # 配置管理
│ └── config.go # 配置结构定义
├── dao/ # 数据访问层
│ ├── base_dao.go # 基础DAO
│ ├── user_dao.go # 用户数据访问
│ ├── student_dao.go # 学员数据访问
│ ├── training_dao.go # 培训数据访问
│ ├── develop_dao.go # 发展历程数据访问
│ └── reward_dao.go # 奖惩记录数据访问
├── entity/ # 数据模型
│ ├── user.go # 用户实体
│ ├── student.go # 学员实体
│ ├── training.go # 培训实体
│ ├── development.go # 发展历程实体
│ └── reward.go # 奖惩记录实体
├── middleware/ # 中间件
│ └── auth.go # JWT认证中间件
├── pkg/ # 公共包
│ └── jwt/
│ └── jwt.go # JWT工具函数
├── service/ # 业务逻辑层
│ ├── auth_service.go # 认证业务逻辑
│ ├── student_service.go # 学员管理业务逻辑
│ ├── training_service.go# 培训管理业务逻辑
│ └── stats_service.go # 统计分析业务逻辑
├── util/ # 工具函数
│ └── password.go # 密码加密工具
├── db/ # 数据库相关
│ ├── db.go # 数据库连接
│ └── restore.go # 数据恢复
├── main.go # 程序入口
├── API_DOCUMENTATION.md # API接口文档
├── config.ini # 环境配置示例
├── go.mod # Go模块文件
└── README.md # 项目说明文档
```

## 📋 环境要求

- **Go**: 1.21 或更高版本
- **MySQL**: 5.7 或更高版本
- **浏览器**: 支持HTML5和JavaScript的现代浏览器（用于API测试）

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd dangjiantest
```

### 2. 安装依赖

```bash
go mod tidy
```



### 3. 配置环境变量

复制 `config.env` 文件并修改数据库配置：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=你的用户名
DB_PASSWORD=你的密码
DB_NAME=party_school

# JWT密钥
JWT_SECRET=your_jwt_secret_key

# 服务器配置
APP_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=development
```



### 4. 创建数据库

```sql
CREATE DATABASE party_school CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```



### 5. 启动服务器

#### 方法一：使用启动脚本（推荐）

```bash
# Windows
start_server.bat
```

#### 方法二：手动启动

```bash
cd behind
go run main.go
```



### 6. 测试API

服务器启动后：

1. **健康检查**: 访问 `http://localhost:8080/health`
2. **API测试**: 打开 `api_test_utf8.html` 文件进行完整测试
3. **API文档**: 查看 `API_DOCUMENTATION.md` 了解详细接口

### 7. 默认账号

- **管理员**: `admin` / `password`
- **教师**: `teacher1` / `password`（所属党支部：计算机学院2021级本科生党支部）
- **学生**: `student1` / `password`

## 📚 API接口文档

### 核心接口概览

#### 🔑 用户认证

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `GET /api/v1/profile` - 获取用户信息

#### 👨‍🎓 学员管理

- `GET /api/v1/students` - 获取学员列表（分页，支持按党支部筛选）
- `GET /api/v1/students/branches` - 获取所有党支部列表（仅管理员）
- `GET /api/v1/students/{id}` - 获取学员详情
- `POST /api/v1/students` - 创建学员
- `PUT /api/v1/students/{id}` - 更新学员信息
- `DELETE /api/v1/students/{id}` - 删除学员（仅管理员）
- `POST /api/v1/students/{id}/photo` - 上传学员照片

#### 📈 发展历程

- `GET /api/v1/students/{id}/development` - 获取发展历程
- `PUT /api/v1/students/{id}/development` - 更新发展历程（支持文本格式时间输入）

**注意**：党员发展历程的时间字段支持任意文本格式输入（如"2024年10月8日"），不进行日期格式校验。

#### 🏆 奖惩记录

- `GET /api/v1/students/{id}/rewards` - 获取奖惩记录
- `POST /api/v1/students/{id}/rewards` - 创建奖惩记录

#### 📚 培训管理

- `GET /api/v1/trainings` - 获取培训列表（分页）
- `GET /api/v1/trainings/{id}` - 获取培训详情
- `POST /api/v1/trainings` - 创建培训记录
- `PUT /api/v1/trainings/{id}` - 更新培训记录
- `DELETE /api/v1/trainings/{id}` - 删除培训记录

#### 📊 统计分析

- `GET /api/v1/statistics/trainings` - 获取培训统计

#### 🏥 系统状态

- `GET /health` - 健康检查

### 详细文档

完整的API接口文档请查看：[API_DOCUMENTATION.md](behind\API_DOCUMENTATION_NEW.md)

**重要更新**：
- ✅ 党员发展信息的时间字段已改为文本格式，支持任意文本输入（如"2024年10月8日"）
- ✅ 所有接口的权限控制已完善
- ✅ 数据库迁移脚本已提供（见"数据库迁移"章节）

## 🗃️ 数据库结构

### 用户表 (users)

- id: 主键
- username: 用户名（唯一）
- password: 密码（加密）
- role: 角色（admin/teacher/student）
- name: 姓名
- created_at: 创建时间
- updated_at: 更新时间

### 学员表 (students)

- id: 主键
- student_no: 学号（唯一）
- name: 姓名
- gender: 性别（男/女）
- ethnicity: 民族
- birth_date: 出生日期
- education: 学历
- phone: 电话
- id_card: 身份证号
- major_class: 专业班级
- type: 政治面貌（群众/入党积极分子/发展对象/预备党员/正式党员）
- admission_date: 入学日期
- conversion_date: 转正日期
- created_at: 创建时间
- updated_at: 更新时间

### 培训表 (trainings)

- id: 主键
- period: 期次
- student_id: 学员ID（外键）
- unit: 培训单位
- start_date: 开始日期
- end_date: 结束日期
- score: 成绩（优秀/良好/及格/不及格）
- certificate_no: 证书编号
- created_at: 创建时间
- updated_at: 更新时间

### 发展历程表 (developments)

- id: 主键
- student_id: 学员ID（外键）
- join_league_date: 入团时间（VARCHAR(50)）
- apply_date: 递交入党申请书时间（VARCHAR(50)）
- party_talk_date: 党组织派人谈话时间（VARCHAR(50))
- league_recommend_date: 团组织推优时间（VARCHAR(50)）
- activist_date: 确立为入党积极分子时间（VARCHAR(50)）
- candidate_date: 确定发展对象时间（VARCHAR(50)）
- superior_report_date: 提交上级党组织备案时间（VARCHAR(50)）
- committee_review_date: 支委会审查时间（VARCHAR(50)）
- party_committee_pre_review_date: 党委预审通过时间（VARCHAR(50)）
- branch_meeting_date: 支部大会讨论时间（发展）（VARCHAR(50)）
- conversion_apply_date: 递交转正申请书时间（VARCHAR(50)）
- conversion_branch_meeting_date: 支部大会讨论时间（转正）（VARCHAR(50)）
- probation_date: 预备党员时间（VARCHAR(50)）
- conversion_date: 转正时间（VARCHAR(50)）
- transfer_date: 党组织关系转接时间（VARCHAR(50)）
- introduction_date: 党组织关系介绍信时间（VARCHAR(50)）
- status: 当前状态（masses/activist/candidate/probationary/formal）
- created_at: 创建时间
- updated_at: 更新时间

**注意**：所有时间字段均为文本格式（VARCHAR(50)），支持任意文本输入，如"2024年10月8日"、"2024-10-08"等格式，不进行日期格式校验。

### 奖惩表 (rewards)

- id: 主键
- student_id: 学员ID（外键）
- type: 类型（奖励/处分）
- content: 内容
- date: 日期
- created_at: 创建时间
- updated_at: 更新时间

## 💻 开发指南

### 添加新功能

1. 在 `entity/` 中定义数据模型
2. 在 `dao/` 中实现数据访问逻辑
3. 在 `service/` 中实现业务逻辑
4. 在 `api/` 中实现HTTP接口
5. 在 `main.go` 中注册路由

### 代码规范

- 遵循Go语言编码规范
- 使用 `go fmt` 格式化代码
- 添加必要的注释
- 编写单元测试

## 🐛 故障排除

### 常见问题

1. **Go工具链问题**
   - 如果遇到 "package xxx is not in std" 错误，请检查Go安装
   - 尝试重新安装Go或使用系统PATH中的Go
2. **数据库连接问题**
   - 检查MySQL服务是否启动
   - 验证数据库配置信息
   - 确认数据库用户权限
3. **端口占用问题**
   - 修改 `config.ini` 中的 `SERVER_PORT` 配置
   - 或使用 `netstat -ano | findstr :8080` 查找占用进程

### 日志查看

程序运行时会输出详细的日志信息，包括：

- 数据库连接状态
- 服务器启动信息
- API请求日志
- 错误信息

## 📖 相关文档

- **[API接口文档](API_DOCUMENTATION.md)** - 完整的API接口说明和使用示例
- **[角色关系说明](ROLE_RELATIONSHIP.md)** - 管理员/教师/学生角色关系和权限说明

## 💻 开发指南

### 添加新功能

1. 在 `entity/` 中定义数据模型
2. 在 `dao/` 中实现数据访问逻辑
3. 在 `service/` 中实现业务逻辑
4. 在 `api/` 中实现HTTP接口
5. 在 `main.go` 中注册路由
6. 更新 `api_test_utf8.html` 测试页面

### 代码规范

- 遵循Go语言编码规范
- 使用 `go fmt` 格式化代码
- 添加必要的注释
- 编写单元测试

## 🐛 故障排除

### 常见问题

1. **Go工具链问题**
   - 检查Go安装和版本
   - 尝试重新安装Go或使用系统PATH中的Go
2. **数据库连接问题**
   - 检查MySQL服务是否启动
   - 验证数据库配置信息
   - 确认数据库用户权限
3. **端口占用问题**
   - 修改 `config.env` 中的 `APP_PORT` 配置
   - 使用 `netstat -ano | findstr :8080` 查找占用进程
4. **测试页面乱码问题**
   - 确保使用 `api_test_utf8.html` 文件
   - 检查浏览器编码设置

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue或Pull Request来改进项目

## 📞 联系方式

- **开发团队**: 新思路实验室Geek组2023级全体
- **项目地址**: [GitHub Repository]
- **问题反馈**: 请通过GitHub Issues提交
