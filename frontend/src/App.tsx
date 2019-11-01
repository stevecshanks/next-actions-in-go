import React, { useState, useEffect } from "react";
import "./App.css";

interface Action {
  id: string;
  name: string;
}

const App: React.FC = () => {
  throw Error("uhoh");
  const [actions, setActions] = useState<Action[]>([]);

  const fetchActions = () => {
    fetch("/actions")
      .then(response => response.json())
      .then(json => setActions(json.data))
      .catch(error => console.log("An error occurred:", error));
  };

  useEffect(fetchActions, []);

  return (
    <div>
      <h1>Next Actions</h1>
      <ul>
        {actions.map((action: Action) => (
          <li key={action.id}>{action.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default App;
