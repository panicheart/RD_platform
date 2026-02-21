#!/bin/bash
#
# Phase 2 Agent Monitoring Script
# Usage: ./scripts/monitor-phase2-agents.sh

set -e

echo "=================================="
echo "Phase 2 Agent Team Monitor"
echo "=================================="
echo ""

# Agent task IDs
AGENTS=(
  "bg_bb7cde4d:WorkflowAgent:Layer1"
  "bg_c20448c0:ProjectAgent:Layer1"
  "bg_77fc654c:DevAgent:Layer2"
  "bg_450f216e:ShelfAgent:Layer2"
  "bg_66d64f1a:QMAgent:Layer2"
  "bg_d04f788c:DesktopAgent:Layer2"
)

echo "Checking agent statuses..."
echo ""

for agent_info in "${AGENTS[@]}"; do
  IFS=':' read -r task_id agent_name layer <<< "$agent_info"
  
  echo "----------------------------------------"
  echo "Agent: $agent_name ($layer)"
  echo "Task ID: $task_id"
  
  # Check for deliverables
  case $agent_name in
    WorkflowAgent)
      echo "Expected: workflow.go, statemachine.go, activity.go, review.go"
      check_files \
        "services/api/models/workflow.go" \
        "services/api/services/statemachine.go" \
        "services/api/handlers/activity.go" \
        "services/api/handlers/review.go"
      ;;
    ProjectAgent)
      echo "Expected: gitea.go, git.go, GanttChart.tsx"
      check_files \
        "services/api/clients/gitea.go" \
        "services/api/services/git.go" \
        "apps/web/src/components/projects/GanttChart.tsx"
      ;;
    DevAgent)
      echo "Expected: ProcessFlow.tsx, rdp_protocol.go, ActivityExecute.tsx"
      check_files \
        "apps/web/src/components/development/ProcessFlow.tsx" \
        "services/api/services/rdp_protocol.go" \
        "apps/web/src/pages/development/ActivityExecute.tsx"
      ;;
    ShelfAgent)
      echo "Expected: product.go, technology.go, ProductList.tsx"
      check_files \
        "services/api/models/product.go" \
        "services/api/models/technology.go" \
        "apps/web/src/pages/shelf/ProductList.tsx"
      ;;
    QMAgent)
      echo "Expected: requirement.go, defect.go, Requirements.tsx"
      check_files \
        "services/api/models/requirement.go" \
        "services/api/models/defect.go" \
        "apps/web/src/pages/quality/Requirements.tsx"
      ;;
    DesktopAgent)
      echo "Expected: rdp-helper/ directory"
      if [ -d "desktop/rdp-helper" ]; then
        echo "  ✓ desktop/rdp-helper exists"
      else
        echo "  ⏳ Waiting for creation"
      fi
      ;;
  esac
  echo ""
done

echo "=================================="
echo "Dependency Status"
echo "=================================="
echo ""
echo "Layer 1 (Foundation):"
echo "  - WorkflowAgent: Building state machine, activity flow, DCP review"
echo "  - ProjectAgent:  Building Gitea integration, Git version, Gantt chart"
echo ""
echo "Layer 2 (Business Logic):"
echo "  - DevAgent:      Building process view, rdp protocol, activity panel"
echo "  - ShelfAgent:    Building product shelf, shopping cart, tech tree"
echo "  - QMAgent:       Building requirements, changes, defects, quality gates"
echo "  - DesktopAgent:  Building protocol handler, local app integration, Git auto-commit"
echo ""
echo "Key Dependency:"
echo "  DesktopAgent waiting for DevAgent's rdp:// protocol definition"
echo ""

# Function to check file existence
check_files() {
  local found=0
  local total=$#
  
  for file in "$@"; do
    if [ -f "$file" ]; then
      echo "  ✓ $file"
      ((found++))
    else
      echo "  ⏳ $file"
    fi
  done
  
  echo ""
  echo "  Progress: $found/$total files"
}

# Run git status to see uncommitted changes
echo "=================================="
echo "Git Status"
echo "=================================="
git status --short | head -20
if [ ${PIPESTATUS[0]} -eq 0 ] && [ -n "$(git status --short)" ]; then
  echo ""
  echo "Uncommitted changes detected. Run 'git add -A && git commit' to save."
fi

echo ""
echo "=================================="
echo "Next Steps"
echo "=================================="
echo ""
echo "1. Monitor agent progress using background_output commands"
echo "2. Check deliverables as agents complete"
echo "3. Coordinate DesktopAgent <-> DevAgent for protocol definition"
echo "4. Run integration tests after Layer 1 completes"
echo "5. Merge and validate all deliverables"
echo ""
echo "To check a specific agent:"
echo "  background_output task_id=\"<task_id>\""
echo ""
