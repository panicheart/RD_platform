# MCP Server: Task Coordinator
# 管理所有Agent任务的中央协调器

import sqlite3
import json
from datetime import datetime
from typing import List, Dict, Optional
import os

class TaskCoordinator:
    """RDP项目任务协调器"""
    
    def __init__(self, db_path: str = "agents/data/tasks.db"):
        self.db_path = db_path
        self.init_database()
    
    def init_database(self):
        """初始化数据库"""
        os.makedirs(os.path.dirname(self.db_path), exist_ok=True)
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # Agent任务表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS agent_tasks (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                task_id TEXT UNIQUE NOT NULL,
                agent_name TEXT NOT NULL,
                phase INTEGER NOT NULL,
                title TEXT NOT NULL,
                status TEXT DEFAULT 'pending',
                priority TEXT DEFAULT 'P1',
                dependencies TEXT,
                input_specs TEXT,
                output_specs TEXT,
                assignee_session TEXT,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                started_at TIMESTAMP,
                completed_at TIMESTAMP,
                git_branch TEXT,
                review_status TEXT
            )
        ''')
        
        # Agent消息表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS agent_messages (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                from_agent TEXT NOT NULL,
                to_agent TEXT,
                message_type TEXT,
                content TEXT NOT NULL,
                context_refs TEXT,
                read_status BOOLEAN DEFAULT 0,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        conn.commit()
        conn.close()
    
    def assign_task(self, task_id: str, agent_name: str, title: str,
                   phase: int = 1, priority: str = "P1",
                   dependencies: Optional[List[str]] = None) -> dict:
        """分配任务给Agent"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        try:
            cursor.execute('''
                INSERT INTO agent_tasks (task_id, agent_name, phase, title, priority, dependencies)
                VALUES (?, ?, ?, ?, ?, ?)
            ''', (task_id, agent_name, phase, title, priority, 
                  json.dumps(dependencies or [])))
            conn.commit()
            
            return {
                "status": "success",
                "task_id": task_id,
                "assigned_to": agent_name,
                "message": f"任务 {task_id} 已分配给 {agent_name}"
            }
        except sqlite3.IntegrityError:
            return {
                "status": "error",
                "message": f"任务 {task_id} 已存在"
            }
        finally:
            conn.close()
    
    def get_agent_tasks(self, agent_name: str) -> List[dict]:
        """获取Agent的所有任务"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            SELECT * FROM agent_tasks 
            WHERE agent_name = ? 
            ORDER BY priority, created_at
        ''', (agent_name,))
        
        tasks = []
        for row in cursor.fetchall():
            tasks.append({
                "task_id": row[1],
                "title": row[4],
                "status": row[5],
                "priority": row[6],
                "dependencies": json.loads(row[7] or "[]"),
                "created_at": row[11]
            })
        
        conn.close()
        return tasks
    
    def update_task_status(self, task_id: str, status: str,
                          session_id: Optional[str] = None) -> dict:
        """更新任务状态"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        timestamp = datetime.now().isoformat()
        
        if status == "in_progress":
            cursor.execute('''
                UPDATE agent_tasks 
                SET status = ?, assignee_session = ?, started_at = ?
                WHERE task_id = ?
            ''', (status, session_id, timestamp, task_id))
        elif status == "completed":
            cursor.execute('''
                UPDATE agent_tasks 
                SET status = ?, completed_at = ?
                WHERE task_id = ?
            ''', (status, timestamp, task_id))
        else:
            cursor.execute('''
                UPDATE agent_tasks SET status = ? WHERE task_id = ?
            ''', (status, task_id))
        
        conn.commit()
        conn.close()
        
        return {
            "status": "success",
            "task_id": task_id,
            "new_status": status
        }
    
    def send_message(self, from_agent: str, to_agent: str,
                    message_type: str, content: str,
                    context_refs: Optional[List[str]] = None) -> dict:
        """发送消息给其他Agent"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            INSERT INTO agent_messages (from_agent, to_agent, message_type, content, context_refs)
            VALUES (?, ?, ?, ?, ?)
        ''', (from_agent, to_agent, message_type, content, 
              json.dumps(context_refs or [])))
        
        conn.commit()
        conn.close()
        
        return {
            "status": "success",
            "message": f"消息已发送给 {to_agent}"
        }
    
    def get_messages(self, agent_name: str, unread_only: bool = False) -> List[dict]:
        """获取Agent的消息"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        if unread_only:
            cursor.execute('''
                SELECT * FROM agent_messages 
                WHERE to_agent = ? AND read_status = 0
                ORDER BY created_at DESC
            ''', (agent_name,))
        else:
            cursor.execute('''
                SELECT * FROM agent_messages 
                WHERE to_agent = ? OR to_agent IS NULL
                ORDER BY created_at DESC
                LIMIT 50
            ''', (agent_name,))
        
        messages = []
        for row in cursor.fetchall():
            messages.append({
                "id": row[0],
                "from": row[1],
                "type": row[3],
                "content": row[4],
                "context": json.loads(row[5] or "[]"),
                "read": row[6],
                "timestamp": row[7]
            })
        
        conn.close()
        return messages
    
    def get_phase_status(self, phase: int) -> dict:
        """获取Phase整体状态"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            SELECT status, COUNT(*) FROM agent_tasks 
            WHERE phase = ?
            GROUP BY status
        ''', (phase,))
        
        status_counts = dict(cursor.fetchall())
        total = sum(status_counts.values())
        completed = status_counts.get('completed', 0)
        progress = (completed / total * 100) if total > 0 else 0
        
        conn.close()
        
        return {
            "phase": phase,
            "total_tasks": total,
            "completed": completed,
            "in_progress": status_counts.get('in_progress', 0),
            "pending": status_counts.get('pending', 0),
            "progress": f"{progress:.1f}%"
        }

# MCP Server入口
def main():
    """MCP Server主函数"""
    import sys
    
    coordinator = TaskCoordinator()
    
    # 读取标准输入的MCP请求
    for line in sys.stdin:
        try:
            request = json.loads(line)
            method = request.get("method")
            params = request.get("params", {})
            
            if method == "assign_task":
                result = coordinator.assign_task(**params)
            elif method == "get_agent_tasks":
                result = coordinator.get_agent_tasks(**params)
            elif method == "update_task_status":
                result = coordinator.update_task_status(**params)
            elif method == "send_message":
                result = coordinator.send_message(**params)
            elif method == "get_messages":
                result = coordinator.get_messages(**params)
            elif method == "get_phase_status":
                result = coordinator.get_phase_status(**params)
            else:
                result = {"status": "error", "message": f"未知方法: {method}"}
            
            # 输出响应
            print(json.dumps(result), flush=True)
            
        except json.JSONDecodeError:
            print(json.dumps({"status": "error", "message": "无效的JSON"}), flush=True)
        except Exception as e:
            print(json.dumps({"status": "error", "message": str(e)}), flush=True)

if __name__ == "__main__":
    main()
