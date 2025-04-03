import TaskForm from '../components/tasks/TaskForm';
import Layout from '../components/layout/Layout';

const EditTaskPage = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto py-6">
        <TaskForm isEditing={true} />
      </div>
    </Layout>
  );
};

export default EditTaskPage;
