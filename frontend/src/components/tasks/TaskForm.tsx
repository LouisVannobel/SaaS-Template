
import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import TaskService, { CreateTaskRequest, UpdateTaskRequest } from '../../services/task.service';

interface TaskFormProps {
  isEditing?: boolean;
}

const TaskForm = ({ isEditing = false }: TaskFormProps) => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('pending');
  const [dueDate, setDueDate] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    // Si on est en mode édition, charger les données de la tâche
    if (isEditing && id) {
      const fetchTask = async () => {
        try {
          setLoading(true);
          const task = await TaskService.getTaskById(parseInt(id));
          setTitle(task.title);
          setDescription(task.description);
          setStatus(task.status);
          if (task.due_date) {
            // Formater la date pour l'input date
            const date = new Date(task.due_date);
            setDueDate(date.toISOString().split('T')[0]);
          }
        } catch (err: any) {
          const errorMessage = err.response?.data?.error || 
                             (typeof err.response?.data === 'object' ? JSON.stringify(err.response.data) : err.response?.data) || 
                             err.message || 
                             'Erreur lors du chargement de la tâche';
          setError(errorMessage);
        } finally {
          setLoading(false);
        }
      };

      fetchTask();
    }
  }, [isEditing, id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validation simple
    if (!title) {
      setError('Le titre est requis');
      return;
    }

    try {
      setLoading(true);
      
      const taskData: CreateTaskRequest | UpdateTaskRequest = {
        title,
        description,
        status,
        due_date: dueDate || undefined
      };

      if (isEditing && id) {
        await TaskService.updateTask(parseInt(id), taskData);
      } else {
        await TaskService.createTask(taskData);
      }

      navigate('/tasks');
    } catch (err: any) {
      console.error('Erreur de soumission:', err);
      const errorMessage = err.response?.data?.error || 
                         (typeof err.response?.data === 'object' ? JSON.stringify(err.response.data) : err.response?.data) || 
                         err.message || 
                         `Erreur lors de ${isEditing ? 'la modification' : 'la création'} de la tâche`;
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card max-w-2xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">
        {isEditing ? 'Modifier la tâche' : 'Créer une nouvelle tâche'}
      </h2>
      
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="title" className="block text-gray-700 mb-2">Titre *</label>
          <input
            type="text"
            id="title"
            className="form-input"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={loading}
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="description" className="block text-gray-700 mb-2">Description</label>
          <textarea
            id="description"
            className="form-input min-h-[100px]"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            disabled={loading}
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="status" className="block text-gray-700 mb-2">Statut</label>
          <select
            id="status"
            className="form-input"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            disabled={loading}
          >
            <option value="pending">À faire</option>
            <option value="in-progress">En cours</option>
            <option value="completed">Terminé</option>
          </select>
        </div>
        
        <div className="mb-6">
          <label htmlFor="dueDate" className="block text-gray-700 mb-2">Date d'échéance</label>
          <input
            type="date"
            id="dueDate"
            className="form-input"
            value={dueDate}
            onChange={(e) => setDueDate(e.target.value)}
            disabled={loading}
          />
        </div>
        
        <div className="flex items-center justify-between">
          <button
            type="button"
            className="btn-secondary"
            onClick={() => navigate('/tasks')}
            disabled={loading}
          >
            Annuler
          </button>
          <button
            type="submit"
            className="btn-primary"
            disabled={loading}
          >
            {loading ? 'Chargement...' : (isEditing ? 'Mettre à jour' : 'Créer')}
          </button>
        </div>
      </form>
    </div>
  );
};

export default TaskForm;
