import { ReactEventHandler, useEffect, useRef } from "react";
import "./TaskForm.css";
import { createTask } from "../../services/api";

interface TaskParams {
  showModal: boolean;
  closeModal: ReactEventHandler;
}

function TaskForm({ showModal, closeModal }: TaskParams) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  function handleTaskSubmit(event: React.FormEvent<HTMLFormElement>) {
    const handleCreateTask = async () => {
      event.preventDefault();
      const form = event.currentTarget;
      const formData = new FormData(form);
      const serviceResponse = await createTask(formData);

      if (serviceResponse.success) {
      }
    };
    handleCreateTask();
  }

  useEffect(() => {
    if (showModal) {
      dialogRef.current?.showModal();
    } else {
      dialogRef.current?.close();
    }
  }, [showModal]);

  return (
    <div className="task-dialog-container">
      <dialog className="task-dialog" ref={dialogRef} onCancel={closeModal}>
        <h1>New Task</h1>
        <button className="task-dialog-close-button" onClick={closeModal}>
          X
        </button>
        <form
          id="create-task"
          action=""
          method="POST"
          onSubmit={handleTaskSubmit}
        >
          <label htmlFor="name">Name</label>
          <input
            type="text"
            name="name"
            id="name"
            defaultValue=""
            placeholder="Enter Task Name"
            maxLength={15}
            required
          ></input>
          <label htmlFor="description">Description</label>
          <input
            type="text"
            name="description"
            id="description"
            defaultValue=""
            placeholder="Enter Task Description"
            required
          ></input>
          <label htmlFor="priority">Priority</label>
          <input
            type="number"
            name="priority"
            id="priority"
            defaultValue={0}
            min={0}
            max={5}
            placeholder="Enter Task Priority (0-5)"
          ></input>
          <label htmlFor="tags">Tags</label>
          <select name="tags" id="tags">
            <option value="tag1">Tag1</option>
            <option value="tag2">Tag2</option>
          </select>
          <button className="task-dialog-submit-button" type="submit">
            Create
          </button>
        </form>
      </dialog>
    </div>
  );
}

export default TaskForm;
