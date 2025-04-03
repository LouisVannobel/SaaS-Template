import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

const ProtectedRoute = () => {
  const { isAuthenticated, loading } = useAuth();

  // Si l'authentification est en cours de vérification, afficher un indicateur de chargement
  if (loading) {
    return <div className="text-center py-8">Chargement...</div>;
  }

  // Si l'utilisateur n'est pas authentifié, rediriger vers la page de connexion
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Si l'utilisateur est authentifié, afficher le contenu de la route
  return <Outlet />;
};

export default ProtectedRoute;
