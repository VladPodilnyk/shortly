import { HeadingLarge } from 'baseui/typography';
import { Layout } from './components/Layout';
import { UserForm } from './components/UserForm';

export function App() {
  return (
    <Layout>
      <HeadingLarge marginBottom="10rem">Shortly</HeadingLarge>
      <UserForm />
    </Layout>
  );
}
