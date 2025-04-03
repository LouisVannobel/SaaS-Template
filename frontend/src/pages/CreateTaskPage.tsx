import TaskForm from '../components/tasks/TaskForm';
import Layout from '../components/layout/Layout';

const CreateTaskPage = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto py-6">
        <TaskForm isEditing={false} />
      </div>
    </Layout>
  );
};

export default CreateTaskPage;
