import api from './api';

export interface Task {
  id: number;
  title: string;
  description: string;
  status: string;
  due_date?: string;
  user_id: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTaskRequest {
  title: string;
  description: string;
  status: string;
  due_date?: string;
}

export interface UpdateTaskRequest {
  title: string;
  description: string;
  status: string;
  due_date?: string;
}

const TaskService = {
  getAllTasks: async (): Promise<Task[]> => {
    const response = await api.get('/tasks');
    // Adapter la réponse au format attendu
    return response.data.tasks || [];
  },

  getTaskById: async (id: number): Promise<Task> => {
    const response = await api.get(`/tasks/${id}`);
    // Adapter la réponse au format attendu
    return response.data.task;
  },

  createTask: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await api.post('/tasks', data);
    // Adapter la réponse au format attendu
    return response.data.task;
  },

  updateTask: async (id: number, data: UpdateTaskRequest): Promise<Task> => {
    const response = await api.put(`/tasks/${id}`, data);
    // Adapter la réponse au format attendu
    return response.data.task;
  },

  deleteTask: async (id: number): Promise<void> => {
    await api.delete(`/tasks/${id}`);
  }
};

export default TaskService;
