import { Layout } from "./components/Layout";
import { Header } from "./components/Header";
import { UserForm } from "./components/UserForm";

// const Centered = styled("div", {
//   display: "flex",
//   justifyContent: "center",
//   alignItems: "center",
//   height: "100%",
// });


function App() {
  return (
    <Layout>
      <Header title="Shortly" />
      <UserForm />
    </Layout>
  );
}

export default App
