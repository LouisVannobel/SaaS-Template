import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/layout/Layout';

const HomePage = () => {
  const { isAuthenticated } = useAuth();

  return (
    <Layout>
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-4xl font-bold text-gray-900 mb-6">
          Bienvenue sur SaaS Template
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Une application SaaS multi-tiers avec authentification et gestion de tâches
        </p>
        
        <div className="bg-white shadow-md rounded-lg p-8 mb-8">
          <h2 className="text-2xl font-bold mb-4">Fonctionnalités</h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div className="p-4">
              <h3 className="text-lg font-semibold mb-2">Authentification</h3>
              <p className="text-gray-600">
                Inscription, connexion et gestion de profil utilisateur sécurisée avec JWT
              </p>
            </div>
            <div className="p-4">
              <h3 className="text-lg font-semibold mb-2">Gestion de tâches</h3>
              <p className="text-gray-600">
                Création, modification, suppression et suivi de vos tâches personnelles
              </p>
            </div>
            <div className="p-4">
              <h3 className="text-lg font-semibold mb-2">Architecture multi-tiers</h3>
              <p className="text-gray-600">
                Frontend React, API backend Go et base de données PostgreSQL
              </p>
            </div>
          </div>
        </div>
        
        <div className="mt-8">
          {isAuthenticated ? (
            <Link
              to="/tasks"
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg text-lg transition-colors"
            >
              Accéder à mes tâches
            </Link>
          ) : (
            <div className="space-x-4">
              <Link
                to="/login"
                className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg text-lg transition-colors"
              >
                Se connecter
              </Link>
              <Link
                to="/register"
                className="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-3 px-6 rounded-lg text-lg transition-colors"
              >
                S'inscrire
              </Link>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default HomePage;
