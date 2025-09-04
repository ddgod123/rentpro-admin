# 项目初始化里程碑

**日期：** 2024年（待填写具体日期）
**状态：** 已完成

## 完成内容

### 1. 项目结构搭建
- [x] Go 项目基础架构
- [x] 数据库模型定义
- [x] API 路由结构
- [x] 配置文件管理

### 2. 核心模块
- [x] 用户认证系统
- [x] 权限管理系统
- [x] 基础 CRUD 操作

### 3. 数据库设计
- [x] 系统管理表结构
  - sys_user (用户表)
  - sys_role (角色表)
  - sys_menu (菜单表)
  - sys_dept (部门表)
  - sys_post (岗位表)
  
- [x] 租赁管理表结构
  - sys_building (建筑表)
  - sys_house (房屋表)
  - sys_landlord (房东表)
  - sys_tenant (租户表)
  - sys_agent (中介表)
  - sys_contract (合同表)

### 4. 技术栈确定
- **后端框架：** Go + Gin
- **数据库：** MySQL
- **ORM：** GORM
- **认证：** JWT
- **配置管理：** Viper

## 下一步计划
1. 完善 API 接口
2. 前后端联调
3. 业务逻辑完善
4. 单元测试编写
