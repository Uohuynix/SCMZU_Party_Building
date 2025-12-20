## 👥 学员、教师、管理员之间的关系说明

本系统围绕“学员（student）—教师（teacher）—管理员（admin）”三类角色构建权限模型，用于支撑学院党支部、党小组场景下的智慧党建业务。

### 一、角色定位与层级关系

- **学生（student）**
  - 对应具体的在校学生，是党建工作的服务对象。
  - 每名学生在系统中有唯一的学员档案（`student_no` 通常对应学号）。
  - 学生只能管理和查看**自己的**信息。

- **教师（teacher）**
  - 通常是党支部书记、组织委员或党小组负责人等管理角色。
  - 教师与一个具体的党支部（`branch`）和党小组（`group_name`）绑定。
  - 教师负责**本支部、本小组学生**的档案维护与培训管理。

- **管理员（admin）**
  - 学院/学校层面的系统管理员。
  - 负责全局配置、账号管理以及所有支部、所有学员数据的统一维护。
  - 对系统内所有数据（用户、学员、培训、奖惩等）拥有最高权限。

层级关系可以概括为：

- **管理员**：跨支部、跨小组的全局管理者；
- **教师**：在本支部、本小组内的“局部管理员”；
- **学生**：仅能管理与查看自身相关的数据。

### 二、教师与学员（学生）之间的组织关系

学员实体中与组织关系相关的关键字段如下（见 `entity/Student`）：

- `branch`：所属党支部名称（例如“计算机学院学生第一党支部”）。
- `group_name`：所属党小组名称（例如“第一党小组”）。

教师用户自身也包含 `branch`、`group_name` 字段，用于表示该教师所负责的党支部和党小组。  
在后端访问控制逻辑中，教师只能操作**与自己 `branch` 和 `group_name` 完全匹配**的学生档案。

可以简单理解为：

- 教师 A（支部 = X，小组 = Y）只能管理“支部 X + 小组 Y” 下的学生；
- 教师 B（支部 = X，小组 = Z）与教师 A 管理的是同一支部但不同小组的学生；
- 不同支部的教师互相之间看不到对方支部的学员信息。

### 三、教师可以查看的学生档案信息范围

当教师访问“本支部、本小组”学生档案时，可以查看到的关键信息包括但不限于：

- **基础身份信息**
  - `id`：系统内部主键 ID；
  - `student_no`：学号；
  - `name`：姓名；
  - `gender`：性别（`male`/`female`）；
  - `ethnicity`：民族；
  - `birth_date`：出生日期；
  - `education`：学历；
  - `phone`：联系电话；
  - `id_card`：身份证号；
  - `major_class`：专业与班级。

- **组织与党员发展信息**
  - `type`：政治面貌 / 党员发展阶段（`masses` 群众、`activist` 入党积极分子、`candidate` 发展对象、`probationary` 预备党员、`formal` 正式党员等）；
  - `branch`：所属党支部；
  - `group_name`：所属党小组；
  - `admission_date`：入学日期；
  - `conversion_date`：转正日期（如适用）。

- **其他辅助信息**
  - `photo_url`：学员照片访问地址；
  - `created_at` / `updated_at`：档案创建与最后更新时间。

> 注意：虽然从接口设计上教师可以获取到以上字段，但在实际业务中可以通过前端或接口网关进一步限制敏感字段（如身份证号）的展示范围，以符合隐私保护要求。

### 四、教师视角下的典型使用示例

以下示例以“计算机学院学生第一党支部 - 第一党小组”的教师账号为例，说明教师在系统中的典型操作与可见数据范围。

#### 示例 1：查看本支部、本小组学生列表

- 教师账号信息：
  - 用户名：`teacher001`
  - 角色：`teacher`
  - `branch`：`计算机学院学生第一党支部`
  - `group_name`：`第一党小组`

- 访问接口：
  - `GET /api/v1/students?page=1&pageSize=10`

- 教师在登录后携带 Token 调用该接口时，后端会根据教师的 `branch` 和 `group_name` 自动过滤，只返回：

```json
{
  "data": [
    {
      "id": 1,
      "student_no": "20230001",
      "name": "张三",
      "gender": "male",
      "ethnicity": "汉族",
      "birth_date": "2005-03-12",
      "education": "本科",
      "phone": "13800000001",
      "id_card": "4201**********1234",
      "major_class": "计算机科学与技术2201班",
      "type": "activist",
      "branch": "计算机学院学生第一党支部",
      "group_name": "第一党小组",
      "photo_url": "/uploads/photos/1_1700000000.jpg",
      "admission_date": "2023-09-01",
      "conversion_date": null,
      "created_at": "2025-01-10T09:00:00Z",
      "updated_at": "2025-02-01T15:30:00Z"
    },
    {
      "id": 2,
      "student_no": "20230002",
      "name": "李四",
      "gender": "female",
      "ethnicity": "土家族",
      "birth_date": "2005-06-20",
      "education": "本科",
      "phone": "13800000002",
      "id_card": "4201**********5678",
      "major_class": "计算机科学与技术2201班",
      "type": "candidate",
      "branch": "计算机学院学生第一党支部",
      "group_name": "第一党小组",
      "photo_url": "/uploads/photos/2_1700000500.jpg",
      "admission_date": "2023-09-01",
      "conversion_date": null,
      "created_at": "2025-01-12T10:20:00Z",
      "updated_at": "2025-02-02T08:15:00Z"
    }
  ],
  "total": 2,
  "page": 1
}
```

> 可以看到，返回的所有学生都隶属于 `计算机学院学生第一党支部` 且 `第一党小组`，其他支部/小组学生不会出现在列表中。

#### 示例 2：查看单个学生的详细档案

- 访问接口：
  - `GET /api/v1/students/1`

- 当前教师若与学生 `id=1` 的 `branch` 和 `group_name` 匹配，则会返回类似如下内容：

```json
{
  "id": 1,
  "student_no": "20230001",
  "name": "张三",
  "gender": "male",
  "ethnicity": "汉族",
  "birth_date": "2005-03-12",
  "education": "本科",
  "phone": "13800000001",
  "id_card": "4201**********1234",
  "major_class": "计算机科学与技术2201班",
  "type": "activist",
  "branch": "计算机学院学生第一党支部",
  "group_name": "第一党小组",
  "photo_url": "/uploads/photos/1_1700000000.jpg",
  "admission_date": "2023-09-01",
  "conversion_date": null,
  "created_at": "2025-01-10T09:00:00Z",
  "updated_at": "2025-02-01T15:30:00Z"
}
```

- 若教师尝试访问非本支部、本小组学生，如：
  - 学生 `id=10` 的 `branch = 计算机学院学生第二党支部`，则后端会返回类似：

```json
{
  "error": "no permission to access this student"
}
```

#### 示例 3：教师更新本小组学生的部分信息

- 访问接口：
  - `PUT /api/v1/students/1`

- 请求体示例（仅修改电话与政治面貌）：

```json
{
  "phone": "13800009999",
  "type": "candidate"
}
```

- 在教师与该学生 `branch + group_name` 匹配的前提下，后端会：
  - 校验教师是否有权限访问该学生；
  - 在原有学生记录基础上，仅更新 `phone` 与 `type` 字段；
  - 保持 `student_no`、`branch`、`group_name` 等关键归属字段不被意外清空或改错。

> 通过上述控制，系统在保证教师能高效管理本支部、本小组学员的同时，也能防止跨支部误操作和数据泄露。

---

## 五、党员发展信息管理

### 权限说明

- **学生**：只能查看和修改自己的党员发展信息
- **教师**：可以查看和修改本支部所有学生的党员发展信息
- **管理员**：可以查看和修改所有学生的党员发展信息

### 时间字段格式

党员发展信息中的所有时间字段均采用**文本格式**（VARCHAR(50)），支持任意文本输入：
- 支持中文格式：如"2024年10月8日"
- 支持标准格式：如"2024-10-08"
- 支持其他自定义格式
- **不进行日期格式校验**，直接保存用户输入的文本内容

### 字段说明

党员发展信息包含以下时间节点：
- `join_league_date`: 入团时间
- `apply_date`: 递交入党申请书时间
- `party_talk_date`: 党组织派人谈话时间
- `league_recommend_date`: 团组织推优时间
- `activist_date`: 确立为入党积极分子时间
- `candidate_date`: 确定发展对象时间
- `superior_report_date`: 提交上级党组织备案时间
- `committee_review_date`: 支委会审查时间
- `party_committee_pre_review_date`: 党委预审通过时间
- `branch_meeting_date`: 支部大会讨论时间（发展）
- `conversion_apply_date`: 递交转正申请书时间
- `conversion_branch_meeting_date`: 支部大会讨论时间（转正）
- `probation_date`: 预备党员时间
- `conversion_date`: 转正时间
- `transfer_date`: 党组织关系转接时间
- `introduction_date`: 党组织关系介绍信时间

---

本说明文档可作为理解系统角色关系及教师权限范围的参考，具体的技术实现细节可结合 `README.md` 中的"系统整体说明"和"角色与权限说明"章节以及 `api/student_handler.go` 中的权限控制代码一并阅读。









