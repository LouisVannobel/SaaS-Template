import { Link } from 'react-router-dom';
import Layout from '../components/layout/Layout';

const NotFoundPage = () => {
  return (
    <Layout>
      <div className="max-w-md mx-auto text-center py-16">
        <h1 className="text-6xl font-bold text-gray-900 mb-4">404</h1>
        <h2 className="text-2xl font-semibold text-gray-700 mb-6">Page non trouvu00e9e</h2>
        <p className="text-gray-600 mb-8">
          La page que vous recherchez n'existe pas ou a u00e9tu00e9 du00e9placu00e9e.
        </p>
        <Link
          to="/"
          className="btn-primary inline-block"
        >
          Retour u00e0 l'accueil
        </Link>
      </div>
    </Layout>
  );
};

export default NotFoundPage;
