import React from "react";
import { render } from "@testing-library/react";
import { Action, NextActionsList } from "./NextActionsList";

test("renders the list of actions", () => {
  const actions: Action[] = [
    { id: "1", name: "An action" },
    { id: "2", name: "Another action" },
  ];

  const { getByText } = render(<NextActionsList actions={actions} />);

  const action = getByText("An action");
  const anotherAction = getByText("Another action");
  expect(action).toBeInTheDocument();
  expect(anotherAction).toBeInTheDocument();
});
