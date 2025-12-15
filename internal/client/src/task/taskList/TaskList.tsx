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
          <TaskListEl task={task} />
        ))}
      </div>
    </>
  );
}

function TaskListEl({ task }: TaskListElProps) {
  return (
    <div className="task-list-el-container">
      <div className="task-list-circle"></div>
      <div id={"task-" + task.id} className="task-list-el">
        {task.name}
      </div>
    </div>
  );
}

export default TaskList;
