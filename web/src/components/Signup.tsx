import { useState } from "react";

const Signup = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confPassword, setConfPassword] = useState("");
  const [error, setError] = useState<Error | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    setError(null);

    try {
      const trimmedEmail = email.trim();
      const trimmedPassword = password.trim();
      const trimmedConfPassword = confPassword.trim();

      if (!trimmedEmail || !trimmedPassword || !trimmedConfPassword) {
        throw new Error("All the fields are necessary");
      }

      if (trimmedPassword !== trimmedConfPassword) {
        throw new Error("Passwords do not match");
      }

      const response = await fetch("http://localhost:3000/api/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      if (response.ok) {
        const data = await response.json();

        console.log(data);
      }
    } catch (error) {
      setError(error as Error);
    }
  }

  return (
    <div style={{border: "1px solid gray", padding: 10,}}>
      <form onSubmit={handleSubmit}>
        <h1>Sign Up</h1>
        <div style={{ display: "flex", flexDirection: "column", gap: "5px" }}>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <input
            type="password"
            placeholder="Confirm Password"
            value={confPassword}
            onChange={(e) => setConfPassword(e.target.value)}
          />
          {error && <p style={{ color: "red" }}>{error.message}</p>}
          <button>Submit</button>
        </div>
      </form>
    </div>
  );
};

export default Signup;
