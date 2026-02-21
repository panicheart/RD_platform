# Agent 任务卡片 - Phase 4 优化完善

---

## P4-T1: AnalyticsAgent - 数据仪表盘

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T1 |
| **Agent** | AnalyticsAgent |
| **模块** | 数据分析 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 项目统计 |
| 用户数据 | UserAgent | 人员统计 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 分析模型 | `services/api/models/analytics.go` | .go |
| 分析Handler | `services/api/handlers/analytics.go` | .go |
| 分析服务 | `services/api/services/analytics.go` | .go |
| 分析表迁移 | `database/migrations/016_analytics.sql` | .sql |
| 仪表盘页面 | `apps/web/src/pages/analytics/Dashboard.tsx` | .tsx |
| 统计卡片 | `apps/web/src/components/analytics/StatCard.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T1.1 | 创建分析数据模型 | 模型定义正确 |
| P4-T1.2 | 实现统计数据接口 | 数据正确 |
| P4-T1.3 | 实现仪表盘页面 | 布局正确 |
| P4-T1.4 | 实现统计卡片 | 数据展示正确 |

---

## P4-T2: AnalyticsAgent - 项目统计报表

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T2 |
| **Agent** | AnalyticsAgent |
| **模块** | 数据分析 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 项目统计服务 | `services/api/services/project_stats.go` | .go |
| 进度图表 | `apps/web/src/components/analytics/ProjectProgress.tsx` | .tsx |
| 状态分布图 | `apps/web/src/components/analytics/ProjectStatus.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T2.1 | 实现项目数量统计 | 统计正确 |
| P4-T2.2 | 实现进度统计 | 进度显示正确 |
| P4-T2.3 | 实现状态分布 | 饼图显示正确 |

---

## P4-T3: AnalyticsAgent - 人员绩效统计

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T3 |
| **Agent** | AnalyticsAgent |
| **模块** | 数据分析 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户数据 | UserAgent | 用户信息 |
| 项目数据 | ProjectAgent | 完成情况 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 绩效服务 | `services/api/services/performance.go` | .go |
| 贡献度排名 | `apps/web/src/components/analytics/ContributionRanking.tsx` | .tsx |
| 工作量统计 | `apps/web/src/components/analytics/WorkloadChart.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T3.1 | 实现工作量统计 | 统计正确 |
| P4-T3.2 | 实现贡献度排名 | 排名正确 |
| P4-T3.3 | 实现ECharts图表 | 图表显示正确 |

---

## P4-T4: AnalyticsAgent - 报表生成

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T4 |
| **Agent** | AnalyticsAgent |
| **模块** | 数据分析 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 统计数据 | P4-T1 | 统计结果 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 报表服务 | `services/api/services/report.go` | .go |
| 报表Handler | `services/api/handlers/report.go` | .go |
| PDF导出 | `apps/web/src/utils/pdf_export.ts` | .ts |
| Excel导出 | `apps/web/src/utils/excel_export.ts` | .ts |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T4.1 | 实现PDF报表生成 | 生成正常 |
| P4-T4.2 | 实现Excel导出 | 导出正常 |
| P4-T4.3 | 实现报表模板 | 模板正确 |

---

## P4-T5: MonitorAgent - 系统监控仪表盘

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T5 |
| **Agent** | MonitorAgent |
| **模块** | 运维监控 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 系统指标 | 部署环境 | 监控数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 监控模型 | `services/api/models/monitor.go` | .go |
| 监控Handler | `services/api/handlers/monitor.go` | .go |
| 监控服务 | `services/api/services/monitor.go` | .go |
| 监控表迁移 | `database/migrations/017_monitor.sql` | .sql |
| 监控页面 | `apps/web/src/pages/monitor/SystemMonitor.tsx` | .tsx |
| 资源图表 | `apps/web/src/components/monitor/ResourceCharts.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T5.1 | 创建监控数据模型 | 模型定义正确 |
| P4-T5.2 | 实现系统指标采集 | 采集正常 |
| P4-T5.3 | 实现监控仪表盘 | 显示正确 |
| P4-T5.4 | 实现资源图表 | 图表正确 |

---

## P4-T6: MonitorAgent - APM性能监控

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T6 |
| **Agent** | MonitorAgent |
| **模块** | 运维监控 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| API性能 | 部署环境 | 请求数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| APM服务 | `services/api/services/apm.go` | .go |
| 慢查询分析 | `apps/web/src/components/monitor/SlowQueries.tsx` | .tsx |
| 端点性能 | `apps/web/src/components/monitor/EndpointPerformance.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T6.1 | 实现性能采集 | 采集正常 |
| P4-T6.2 | 实现慢查询分析 | 分析正确 |
| P4-T6.3 | 实现端点性能 | 性能显示正确 |

---

## P4-T7: MonitorAgent - 日志集中管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T7 |
| **Agent** | MonitorAgent |
| **模块** | 运维监控 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 应用日志 | 各服务 | 日志数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 日志服务 | `services/api/services/logging.go` | .go |
| 日志Handler | `services/api/handlers/log.go` | .go |
| 日志查询页面 | `apps/web/src/pages/monitor/LogQuery.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T7.1 | 实现日志采集 | 采集正常 |
| P4-T7.2 | 实现日志存储 | 存储正常 |
| P4-T7.3 | 实现日志查询 | 查询正常 |

---

## P4-T8: MonitorAgent - 告警机制

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T8 |
| **Agent** | MonitorAgent |
| **模块** | 运维监控 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 监控指标 | P4-T5 | 指标数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 告警服务 | `services/api/services/alert.go` | .go |
| 告警Handler | `services/api/handlers/alert.go` | .go |
| 告警规则配置 | `config/alert_rules.yaml` | .yaml |
| 告警页面 | `apps/web/src/pages/monitor/Alerts.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T8.1 | 实现告警规则配置 | 配置正确 |
| P4-T8.2 | 实现告警检测 | 检测正确 |
| P4-T8.3 | 实现告警通知 | 通知发送 |
| P4-T8.4 | 实现告警页面 | 显示正确 |

---

## P4-T9: MonitorAgent - Prometheus集成

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T9 |
| **Agent** | MonitorAgent |
| **模块** | 运维监控 |
| **优先级** | P2 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Prometheus | 部署环境 | 监控系统 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Prometheus配置 | `config/prometheus.yml` | .yml |
| 告警规则 | `config/alert-rules.yml` | .yml |
| Grafana面板 | `deploy/grafana/dashboards/` | .json |
| 指标端点 | `/metrics` | 端点 |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T9.1 | 实现Prometheus配置 | 配置正确 |
| P4-T9.2 | 实现告警规则 | 规则正确 |
| P4-T9.3 | 实现Grafana面板 | 面板正确 |
| P4-T9.4 | 实现指标端点 | 暴露正确 |

---

## P4-T10: 全局性能优化

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P4-T10 |
| **Agent** | 所有Feature Agent |
| **模块** | 优化 |
| **优先级** | P1 |
| **阶段** | Phase 4 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 性能数据 | MonitorAgent | 性能指标 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 缓存优化 | - | - |
| SQL优化 | - | - |
| 前端优化 | - | - |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P4-T10.1 | 实现Redis缓存 | 缓存生效 |
| P4-T10.2 | 优化SQL查询 | 响应提升 |
| P4-T10.3 | 优化前端加载 | 首屏<3s |

---

## Phase 4 集成测试清单

| 序号 | 测试项 | 验收条件 |
|------|--------|----------|
| IT-50 | 数据仪表盘 | 数据展示正确 |
| IT-51 | 项目统计报表 | 报表生成正常 |
| IT-52 | 人员绩效统计 | 排名正确 |
| IT-53 | 报表导出 | PDF/Excel正常 |
| IT-54 | 系统监控 | 指标采集正常 |
| IT-55 | APM性能监控 | 慢请求记录 |
| IT-56 | 日志查询 | 检索正常 |
| IT-57 | 告警机制 | 告警触发 |
| IT-58 | Prometheus集成 | 指标采集 |
| IT-59 | 性能优化 | 响应时间达标 |
| IT-60 | 最终验收 | 全功能稳定 |
