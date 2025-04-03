import LoginForm from '../components/auth/LoginForm';
import Layout from '../components/layout/Layout';

const LoginPage = () => {
  return (
    <Layout>
      <div className="max-w-md mx-auto py-8">
        <h1 className="text-3xl font-bold text-center mb-8">Connexion</h1>
        <LoginForm />
      </div>
    </Layout>
  );
};

export default LoginPage;
