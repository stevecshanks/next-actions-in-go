import React, { useState, useEffect } from "react";
import Alert from "react-bootstrap/Alert";
import { Action, NextActionsList } from "./NextActionsList";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

type JsonAction = {
  id: string;
  name: string;
  dueBy?: string;
};

const actionsFromJson = (json: JsonAction[]): Action[] =>
  json.map(
    action =>
      new Action({
        ...action,
        dueBy: action.dueBy ? new Date(action.dueBy) : undefined,
      }),
  );

const App: React.FC = () => {
  const [actions, setActions] = useState<Action[]>([]);
  const [errorMessage, setErrorMessage] = useState<String | null>(null);

  const fetchActions = () => {
    fetch("api/actions")
      .then(response => response.json())
      .then(json => setActions(actionsFromJson(json.data)))
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
