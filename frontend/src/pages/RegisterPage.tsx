import RegisterForm from '../components/auth/RegisterForm';
import Layout from '../components/layout/Layout';

const RegisterPage = () => {
  return (
    <Layout>
      <div className="max-w-md mx-auto py-8">
        <h1 className="text-3xl font-bold text-center mb-8">Créer un compte</h1>
        <RegisterForm />
      </div>
    </Layout>
  );
};

export default RegisterPage;
