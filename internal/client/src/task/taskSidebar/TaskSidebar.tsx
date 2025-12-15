import { ReactEventHandler, useState } from "react";
import { ReactComponent as Plus } from "./plus.svg";
import TaskList from "../taskList/TaskList";
import "./TaskSidebar.css";

function populateExampleTaskList() {
  return [
    {
      id: "1",
      name: "Task1",
    },
    {
      id: "2",
      name: "Task2",
    },
  ];
}

interface TaskSidebarParams {
  openModal: ReactEventHandler;
}

function TaskSidebar({ openModal }: TaskSidebarParams) {
  const [recentTasks, setRecentTasks] = useState(populateExampleTaskList());

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
