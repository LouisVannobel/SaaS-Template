import { useState } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const LoginForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login, loading } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validation simple
    if (!email || !password) {
      setError('Email et mot de passe sont requis');
      return;
    }

    try {
      await login(email, password);
      navigate('/tasks'); // Rediriger vers la page des tâches après connexion
    } catch (err: any) {
      console.error('Erreur de connexion:', err);
      const errorMessage = err.response?.data?.error || 
                         (typeof err.response?.data === 'object' ? JSON.stringify(err.response.data) : err.response?.data) || 
                         err.message || 
                         'Email ou mot de passe incorrect';
      setError(errorMessage);
    }
  };

  return (
    <div className="card max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-4">Connexion</h2>
      
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="email" className="block text-gray-700 mb-2">Email</label>
          <input
            type="email"
            id="email"
            className="form-input"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={loading}
          />
        </div>
        
        <div className="mb-6">
          <label htmlFor="password" className="block text-gray-700 mb-2">Mot de passe</label>
          <input
            type="password"
            id="password"
            className="form-input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading}
          />
        </div>
        
        <div className="flex items-center justify-between">
          <button
            type="submit"
            className="btn-primary w-full"
            disabled={loading}
          >
            {loading ? 'Chargement...' : 'Se connecter'}
          </button>
        </div>
        
        <div className="mt-4 text-center">
          <p>Pas encore inscrit ? <a href="/register" className="text-blue-600 hover:underline">S'inscrire</a></p>
        </div>
      </form>
    </div>
  );
};

export default LoginForm;
