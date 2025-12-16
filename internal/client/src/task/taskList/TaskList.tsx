import { Task } from "../../types/types";
import "./TaskList.css";

interface TaskListProps {
  tasks: Task[];
}

interface TaskListElProps {
  task: Task;
}

function TaskList({ tasks }: TaskListProps) {
  return (
    <>
      <h2 className="task-list-title">Recent Tasks</h2>
      <div className="task-list">
        {tasks.map((task) => (
          <TaskListEl key={task.id} task={task} />
        ))}
      </div>
    </>
  );
}

function TaskListEl({ task }: TaskListElProps) {
  return (
    <div className="task-list-el-container">
      <button
        type="button"
        className="task-list-el-button"
        id={"task-" + task.id}
      >
        <div className="task-list-circle"></div>
        <div className="task-name">{task.name}</div>
      </button>
    </div>
  );
}

export default TaskList;
