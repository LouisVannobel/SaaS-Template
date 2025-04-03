import TaskDetail from '../components/tasks/TaskDetail';
import Layout from '../components/layout/Layout';

const TaskDetailPage = () => {
  return (
    <Layout>
      <div className="max-w-4xl mx-auto py-6">
        <TaskDetail />
      </div>
    </Layout>
  );
};

export default TaskDetailPage;
