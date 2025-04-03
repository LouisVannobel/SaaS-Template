import TaskList from '../components/tasks/TaskList';
import Layout from '../components/layout/Layout';

const TasksPage = () => {
  return (
    <Layout>
      <div className="max-w-6xl mx-auto py-6">
        <TaskList />
      </div>
    </Layout>
  );
};

export default TasksPage;
