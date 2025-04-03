import { useState, useEffect } from 'react';
import TaskService, { Task } from '../../services/task.service';
import { Link } from 'react-router-dom';

const TaskList = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        setLoading(true);
        console.log('Récupération des tâches...');
        const data = await TaskService.getAllTasks();
        console.log('Tâches récupérées:', data);
        setTasks(data);
        setError('');
      } catch (err: any) {
        console.error('Erreur lors de la récupération des tâches:', err);
        setError(
          err.response?.data?.error || 
          err.message || 
          'Erreur lors du chargement des tâches'
        );
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, []);

  const handleDelete = async (id: number) => {
    if (window.confirm('Êtes-vous sûr de vouloir supprimer cette tâche ?')) {
      try {
        await TaskService.deleteTask(id);
        setTasks(tasks.filter(task => task.id !== id));
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

  if (loading) {
    return <div className="text-center py-8">Chargement...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Mes tâches</h2>
        <Link to="/tasks/new" className="btn-primary">
          Nouvelle tâche
        </Link>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {tasks.length === 0 ? (
        <div className="text-center py-8 bg-gray-50 rounded-lg">
          <p className="text-gray-500">Aucune tâche trouvée</p>
          <Link to="/tasks/new" className="text-blue-600 hover:underline mt-2 inline-block">
            Créer votre première tâche
          </Link>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {tasks.map((task) => (
            <div key={task.id} className="card hover:shadow-lg transition-shadow">
              <div className="flex justify-between items-start">
                <h3 className="text-xl font-semibold mb-2">{task.title}</h3>
                <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusBadgeClass(task.status)}`}>
                  {task.status}
                </span>
              </div>
              <p className="text-gray-600 mb-4 line-clamp-2">{task.description}</p>
              
              {task.due_date && (
                <p className="text-sm text-gray-500 mb-4">
                  Échéance: {new Date(task.due_date).toLocaleDateString()}
                </p>
              )}
              
              <div className="flex justify-end space-x-2 mt-4">
                <Link 
                  to={`/tasks/${task.id}`} 
                  className="text-blue-600 hover:text-blue-800"
                >
                  Voir
                </Link>
                <Link 
                  to={`/tasks/${task.id}/edit`} 
                  className="text-green-600 hover:text-green-800"
                >
                  Modifier
                </Link>
                <button 
                  onClick={() => handleDelete(task.id)} 
                  className="text-red-600 hover:text-red-800"
                >
                  Supprimer
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default TaskList;
