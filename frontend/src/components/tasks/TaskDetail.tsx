import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import TaskService, { Task } from '../../services/task.service';

const TaskDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTask = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const data = await TaskService.getTaskById(parseInt(id));
        setTask(data);
      } catch (err: any) {
        setError(err.response?.data || 'Erreur lors du chargement de la tâche');
      } finally {
        setLoading(false);
      }
    };

    fetchTask();
  }, [id]);

  const handleDelete = async () => {
    if (!id || !task) return;
    
    if (window.confirm('Êtes-vous sûr de vouloir supprimer cette tâche ?')) {
      try {
        await TaskService.deleteTask(parseInt(id));
        navigate('/tasks');
      } catch (err: any) {
        setError(err.response?.data || 'Erreur lors de la suppression de la tâche');
      }
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'in-progress':
        return 'bg-blue-100 text-blue-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'pending':
        return 'À faire';
      case 'in-progress':
        return 'En cours';
      case 'completed':
        return 'Terminé';
      default:
        return status;
    }
  };

  if (loading) {
    return <div className="text-center py-8">Chargement...</div>;
  }

  if (error) {
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p>{error}</p>
        <Link to="/tasks" className="text-red-700 font-bold hover:underline mt-2 inline-block">
          Retour à la liste des tâches
        </Link>
      </div>
    );
  }

  if (!task) {
    return (
      <div className="text-center py-8">
        <p>Tâche non trouvée</p>
        <Link to="/tasks" className="text-blue-600 hover:underline mt-2 inline-block">
          Retour à la liste des tâches
        </Link>
      </div>
    );
  }

  return (
    <div className="card max-w-2xl mx-auto">
      <div className="flex justify-between items-start mb-4">
        <h2 className="text-2xl font-bold">{task.title}</h2>
        <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusBadgeClass(task.status)}`}>
          {getStatusLabel(task.status)}
        </span>
      </div>

      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2">Description</h3>
        <p className="text-gray-700 whitespace-pre-line">{task.description || 'Aucune description'}</p>
      </div>

      {task.due_date && (
        <div className="mb-6">
          <h3 className="text-lg font-semibold mb-2">Date d'échéance</h3>
          <p className="text-gray-700">{new Date(task.due_date).toLocaleDateString()}</p>
        </div>
      )}

      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2">Informations</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-gray-500 text-sm">Créée le</p>
            <p className="text-gray-700">{new Date(task.created_at).toLocaleDateString()}</p>
          </div>
          <div>
            <p className="text-gray-500 text-sm">Dernière mise à jour</p>
            <p className="text-gray-700">{new Date(task.updated_at).toLocaleDateString()}</p>
          </div>
        </div>
      </div>

      <div className="flex justify-between pt-4 border-t border-gray-200">
        <Link to="/tasks" className="btn-secondary">
          Retour
        </Link>
        <div className="space-x-2">
          <Link to={`/tasks/${task.id}/edit`} className="btn-primary">
            Modifier
          </Link>
          <button onClick={handleDelete} className="btn-danger">
            Supprimer
          </button>
        </div>
      </div>
    </div>
  );
};

export default TaskDetail;
