#!/usr/bin/env python3
"""
RDP 5-Agent Team 任务协调器
管理5个Agent之间的任务分配、依赖关系和进度同步
"""

import sqlite3
import json
import os
from datetime import datetime
from pathlib import Path

DB_PATH = Path(__file__).parent.parent / "data" / "5agent_tasks.db"

class AgentCoordinator:
    def __init__(self):
        self.db_path = DB_PATH
        self.init_db()
    
    def init_db(self):
        """初始化数据库"""
        os.makedirs(self.db_path.parent, exist_ok=True)
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # 任务表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS tasks (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                task_id TEXT UNIQUE NOT NULL,
                title TEXT NOT NULL,
                description TEXT,
                assignee TEXT NOT NULL,  -- PM-Agent, Architect-Agent, Backend-Agent, Frontend-Agent, DevOps-Agent
                phase INTEGER NOT NULL,
                status TEXT DEFAULT 'pending',  -- pending, in_progress, review, completed, blocked
                priority TEXT DEFAULT 'P1',  -- P0, P1, P2
                dependencies TEXT,  -- JSON array of task_ids
                deliverables TEXT,  -- JSON array
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                started_at TIMESTAMP,
                completed_at TIMESTAMP,
                notes TEXT
            )
        ''')
        
        # Agent状态表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS agent_status (
                agent_name TEXT PRIMARY KEY,
                current_task TEXT,
                status TEXT DEFAULT 'idle',  -- idle, working, blocked
                last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                progress_percent INTEGER DEFAULT 0
            )
        ''')
        
        # 消息表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS messages (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                from_agent TEXT NOT NULL,
                to_agent TEXT,
                type TEXT NOT NULL,  -- task, question, review, blocker, announcement
                content TEXT NOT NULL,
                task_ref TEXT,
                read_status BOOLEAN DEFAULT 0,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        conn.commit()
        conn.close()
        print(f"✅ 数据库初始化完成: {self.db_path}")
    
    def add_task(self, task_id, title, assignee, phase, description="", 
                 priority="P1", dependencies=None, deliverables=None):
        """添加任务"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            INSERT OR REPLACE INTO tasks 
            (task_id, title, description, assignee, phase, priority, dependencies, deliverables)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ''', (task_id, title, description, assignee, phase, priority,
              json.dumps(dependencies or []), json.dumps(deliverables or [])))
        
        conn.commit()
        conn.close()
        print(f"✅ 任务已添加: {task_id} -> {assignee}")
    
    def assign_task(self, task_id, assignee):
        """分配任务给Agent"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            UPDATE tasks SET assignee = ?, status = 'pending' WHERE task_id = ?
        ''', (assignee, task_id))
        
        conn.commit()
        conn.close()
        print(f"✅ 任务 {task_id} 分配给 {assignee}")
    
    def update_task_status(self, task_id, status, notes=""):
        """更新任务状态"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        now = datetime.now().isoformat()
        
        if status == 'in_progress':
            cursor.execute('''
                UPDATE tasks SET status = ?, started_at = ?, notes = ? WHERE task_id = ?
            ''', (status, now, notes, task_id))
        elif status == 'completed':
            cursor.execute('''
                UPDATE tasks SET status = ?, completed_at = ?, notes = ? WHERE task_id = ?
            ''', (status, now, notes, task_id))
        else:
            cursor.execute('''
                UPDATE tasks SET status = ?, notes = ? WHERE task_id = ?
            ''', (status, notes, task_id))
        
        conn.commit()
        conn.close()
        print(f"✅ 任务 {task_id} 状态更新为: {status}")
    
    def list_tasks(self, assignee=None, phase=None, status=None):
        """列出任务"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        query = "SELECT * FROM tasks WHERE 1=1"
        params = []
        
        if assignee:
            query += " AND assignee = ?"
            params.append(assignee)
        if phase:
            query += " AND phase = ?"
            params.append(phase)
        if status:
            query += " AND status = ?"
            params.append(status)
        
        query += " ORDER BY phase, priority, created_at"
        
        cursor.execute(query, params)
        tasks = cursor.fetchall()
        conn.close()
        
        return tasks
    
    def send_message(self, from_agent, to_agent, msg_type, content, task_ref=None):
        """发送消息"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            INSERT INTO messages (from_agent, to_agent, type, content, task_ref)
            VALUES (?, ?, ?, ?, ?)
        ''', (from_agent, to_agent, msg_type, content, task_ref))
        
        conn.commit()
        conn.close()
        print(f"✅ 消息已发送: {from_agent} -> {to_agent or 'all'}")
    
    def get_messages(self, agent_name, unread_only=False):
        """获取消息"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        query = "SELECT * FROM messages WHERE to_agent = ? OR to_agent IS NULL"
        params = [agent_name]
        
        if unread_only:
            query += " AND read_status = 0"
        
        query += " ORDER BY created_at DESC"
        
        cursor.execute(query, params)
        messages = cursor.fetchall()
        conn.close()
        
        return messages
    
    def init_phase1_tasks(self):
        """初始化Phase 1任务"""
        tasks = [
            # Architect-Agent 任务
            ("P1-A1", "数据库Schema设计", "Architect-Agent", 1, 
             "设计users, projects, activities等核心表结构", "P0", [],
             ["database/migrations/001_init_schema.sql", "docs/data_model.md"]),
            
            ("P1-A2", "API接口规范定义", "Architect-Agent", 1,
             "定义RESTful API规范，包括路径、请求/响应格式", "P0", ["P1-A1"],
             ["services/api/docs/api_spec.md"]),
            
            # Backend-Agent 任务
            ("P1-B1", "用户管理API", "Backend-Agent", 1,
             "实现用户CRUD、认证、RBAC权限", "P0", ["P1-A1", "P1-A2"],
             ["services/api/handlers/user.go", "services/api/services/user.go"]),
            
            ("P1-B2", "项目管理API", "Backend-Agent", 1,
             "实现项目CRUD、成员管理", "P0", ["P1-B1"],
             ["services/api/handlers/project.go", "services/api/services/project.go"]),
            
            # Frontend-Agent 任务
            ("P1-F1", "门户界面开发", "Frontend-Agent", 1,
             "部门首页、个人工作台、通知中心", "P0", [],
             ["apps/web/src/pages/portal/", "apps/web/src/pages/workbench/"]),
            
            ("P1-F2", "用户管理界面", "Frontend-Agent", 1,
             "登录/注册、用户列表、组织架构、Profile", "P0", ["P1-B1"],
             ["apps/web/src/pages/users/"]),
            
            ("P1-F3", "项目管理界面", "Frontend-Agent", 1,
             "项目列表、创建向导、详情页、甘特图", "P0", ["P1-B2"],
             ["apps/web/src/pages/projects/"]),
            
            # DevOps-Agent 任务
            ("P1-D1", "数据库初始化脚本", "DevOps-Agent", 1,
             "创建数据库、用户、初始表结构", "P0", ["P1-A1"],
             ["database/init.sql", "deploy/scripts/init-db.sh"]),
            
            ("P1-D2", "systemd服务配置", "DevOps-Agent", 1,
             "创建rdp-api、casdoor等服务配置", "P1", [],
             ["deploy/systemd/rdp-api.service", "deploy/systemd/rdp-casdoor.service"]),
            
            ("P1-D3", "部署脚本", "DevOps-Agent", 1,
             "一键安装脚本install.sh", "P1", ["P1-B1", "P1-B2", "P1-F1"],
             ["deploy/scripts/install.sh", "deploy/scripts/backup.sh"]),
        ]
        
        for task in tasks:
            self.add_task(*task)
        
        print(f"✅ Phase 1 任务初始化完成，共 {len(tasks)} 个任务")
    
    def print_task_board(self):
        """打印任务看板"""
        print("\n" + "="*80)
        print("📋 RDP 5-Agent Team 任务看板")
        print("="*80)
        
        agents = ["PM-Agent", "Architect-Agent", "Backend-Agent", "Frontend-Agent", "DevOps-Agent"]
        
        for agent in agents:
            tasks = self.list_tasks(assignee=agent)
            if tasks:
                print(f"\n👤 {agent}:")
                for task in tasks:
                    task_id, title, _, _, phase, status, priority, _, _, _, _, _, _ = task[:13]
                    status_icon = {
                        'pending': '⏳',
                        'in_progress': '🔄',
                        'review': '👀',
                        'completed': '✅',
                        'blocked': '❌'
                    }.get(status, '⚪')
                    print(f"  {status_icon} [{priority}] {task_id}: {title} ({status})")
        
        print("\n" + "="*80)


def main():
    import sys
    
    coordinator = AgentCoordinator()
    
    if len(sys.argv) < 2:
        print("""
RDP 5-Agent Team 任务协调器

用法:
  python3 coordinator.py init              # 初始化Phase 1任务
  python3 coordinator.py list              # 列出所有任务
  python3 coordinator.py list <agent>      # 列出指定Agent的任务
  python3 coordinator.py status            # 显示任务看板
  python3 coordinator.py update <task_id> <status> [notes]  # 更新任务状态
        """)
        return
    
    cmd = sys.argv[1]
    
    if cmd == "init":
        coordinator.init_phase1_tasks()
    elif cmd == "list":
        if len(sys.argv) > 2:
            tasks = coordinator.list_tasks(assignee=sys.argv[2])
        else:
            tasks = coordinator.list_tasks()
        
        for task in tasks:
            print(f"{task[1]} | {task[4]} | {task[6]} | {task[3]} | {task[5]}")
    
    elif cmd == "status":
        coordinator.print_task_board()
    
    elif cmd == "update" and len(sys.argv) >= 4:
        coordinator.update_task_status(sys.argv[2], sys.argv[3], sys.argv[4] if len(sys.argv) > 4 else "")
    
    else:
        print("未知命令")


if __name__ == "__main__":
    main()
