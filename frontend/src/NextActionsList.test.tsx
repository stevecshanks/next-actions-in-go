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

test("sorts actions by due date", () => {
  const actions: Action[] = [
    { id: "1", name: "An action with no due date" },
    { id: "2", name: "An action due later", dueBy: new Date(2020, 9, 9) },
    { id: "3", name: "An action due soon", dueBy: new Date(2020, 1, 1) },
  ];

  const { getAllByText } = render(<NextActionsList actions={actions} />);

  const foundActions = getAllByText("An action", { exact: false });

  expect(foundActions).toHaveLength(3);
  expect(foundActions[0].textContent).toMatch(/^An action due soon.*/);
  expect(foundActions[1].textContent).toMatch(/^An action due later.*/);
  expect(foundActions[2].textContent).toMatch(/^An action with no due date.*/);
});
