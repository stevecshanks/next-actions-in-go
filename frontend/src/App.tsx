import React, { useState, useEffect } from "react";
import Alert from "react-bootstrap/Alert";
import ListGroup from "react-bootstrap/ListGroup";
import "bootstrap/dist/css/bootstrap.min.css";
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
    <>
      {errorMessage && <Alert variant="danger">{errorMessage}</Alert>}
      <ListGroup>
        {actions.map((action: Action) => (
          <ListGroup.Item key={action.id}>{action.name}</ListGroup.Item>
        ))}
      </ListGroup>
    </>
  );
};

export default App;
