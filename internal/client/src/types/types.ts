export interface Task {
  id: string;
  name: string;
  created_at: Date;
  updated_at: Date;
  priority: number;
  description: string;
  tags: string[];
}

export interface HttpSuccessResponseMultiple {
  success: boolean;
  tasks: Task[];
}

export interface HttpSuccessResponse {
  success: boolean;
  task: Task[];
}
