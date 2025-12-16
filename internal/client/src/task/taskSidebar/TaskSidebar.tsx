import { ReactEventHandler, useEffect, useState } from "react";
import { ReactComponent as Plus } from "./plus.svg";
import TaskList from "../taskList/TaskList";
import "./TaskSidebar.css";
import { Task } from "../../types/types";
import { getLatestTasks } from "../../services/api";

interface TaskSidebarParams {
  openModal: ReactEventHandler;
}

function TaskSidebar({ openModal }: TaskSidebarParams) {
  const [recentTasks, setRecentTasks] = useState<Task[]>([]);

  useEffect(() => {
    const fetchLatestTasks = async () => {
      const serviceResponse = await getLatestTasks();
      if (serviceResponse.success) {
        setRecentTasks(serviceResponse.tasks);
      }
    };
    fetchLatestTasks();
  }, []);

  return (
    <div className="task-sidebar-header">
      <div className="task-sidebar-container">
        <h2 className="task-sidebar-title">Task Manager App</h2>
        <div className="new-task">
          <button type="button" className="new-task-button" onClick={openModal}>
            <div className="new-task-button-text">
              <Plus />
              New Task
            </div>
          </button>
        </div>
      </div>
      <div className="recent-tasks">
        <TaskList tasks={recentTasks} />
      </div>
    </div>
  );
}

export default TaskSidebar;
