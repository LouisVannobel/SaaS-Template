import { useState } from 'react';
import Layout from '../components/layout/Layout';
import { useAuth } from '../contexts/AuthContext';

const ProfilePage = () => {
  const { user } = useAuth();
  const [name, setName] = useState(user?.name || '');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  // Cette fonction serait implémentée pour mettre à jour le profil utilisateur
  // mais nous n'avons pas encore créé l'endpoint correspondant dans le backend
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess(false);
    
    try {
      // Simuler une mise à jour réussie
      await new Promise(resolve => setTimeout(resolve, 500));
      setSuccess(true);
    } catch (err: any) {
      setError(err.response?.data || 'Une erreur est survenue');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-2xl mx-auto py-8">
        <h1 className="text-3xl font-bold mb-8">Mon profil</h1>
        
        <div className="card">
          {success && (
            <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
              Profil mis à jour avec succès
            </div>
          )}
          
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
                className="form-input bg-gray-100"
                value={user?.email || ''}
                disabled
              />
              <p className="text-sm text-gray-500 mt-1">L'email ne peut pas être modifié</p>
            </div>
            
            <div className="mb-6">
              <label htmlFor="name" className="block text-gray-700 mb-2">Nom</label>
              <input
                type="text"
                id="name"
                className="form-input"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={loading}
              />
            </div>
            
            <div className="flex justify-end">
              <button
                type="submit"
                className="btn-primary"
                disabled={loading}
              >
                {loading ? 'Chargement...' : 'Mettre à jour'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
};

export default ProfilePage;
