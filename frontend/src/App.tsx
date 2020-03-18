import React, { useState, useEffect } from "react";
import Alert from "react-bootstrap/Alert";
import { Action, NextActionsList } from "./NextActionsList";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

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
    <>
      {errorMessage && <Alert variant="danger">{errorMessage}</Alert>}
      <NextActionsList actions={actions} />
    </>
  );
};

export default App;
