import React, { useState, useEffect } from "react";
import "./App.css";

interface Action {
  id: string;
  name: string;
}

const App: React.FC = () => {
  const [actions, setActions] = useState<Action[]>([]);
  const [errorMessage, setErrorMessage] = useState<String | null>(null);

  const fetchActions = () => {
    fetch("api/actions")
      .then(response => response.json())
      .then(json => setActions(json.data))
      .catch(() => setErrorMessage("An error occurred"));
  };

  useEffect(fetchActions, []);

  return (
    <div>
      <h1>Next Actions</h1>
      {errorMessage && <h2>{errorMessage}</h2>}
      <ul>
        {actions.map((action: Action) => (
          <li key={action.id}>{action.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default App;
