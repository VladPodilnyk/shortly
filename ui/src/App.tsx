import { Header } from './components/Header';
import { Layout } from './components/Layout';
import { UserForm } from './components/UserForm';

export function App() {
  return (
    <Layout>
      <Header title="Shortly" />
      <UserForm />
    </Layout>
  );
}
