import "./App.css";
import Login from "./components/Login";
import Signup from "./components/Signup";

function App() {
  return (
    <div style={{ display: "flex", gap: "10px" }}>
      <Signup />
      <Login />
    </div>
  );
}

export default App;
