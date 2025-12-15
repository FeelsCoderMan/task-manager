import TaskForm from "./task/taskForm/TaskForm";
import TaskGrid from "./task/taskGrid/TaskGrid";
import TaskSidebar from "./task/taskSidebar/TaskSidebar";
import "./TaskApp.css";
import { useState } from "react";

function TaskApp() {
  const [isCreateTaskButtonClicked, setIsCreateTaskButtonClicked] =
    useState(false);
  const openTaskForm = () => {
    setIsCreateTaskButtonClicked(true);
  };
  const closeTaskForm = () => {
    setIsCreateTaskButtonClicked(false);
  };

  return (
    <div className="task-manager-app">
      <div className="task-column task-sidebar">
        <TaskSidebar openModal={openTaskForm} />
      </div>
      <div className="task-column task-grid">
        <TaskGrid />
      </div>
      <TaskForm
        showModal={isCreateTaskButtonClicked}
        closeModal={closeTaskForm}
      />
    </div>
  );
}

export default TaskApp;
