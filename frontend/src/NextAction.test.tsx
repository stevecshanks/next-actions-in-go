import React from "react";
import { render } from "@testing-library/react";
import { Action } from "./NextActionsList";
import { NextAction } from "./NextAction";

test("displays action name", () => {
  const anAction: Action = { id: "1", name: "An action" };

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("An action");

  expect(foundAction).toBeInTheDocument();
});

test("displays action due date", () => {
  const anAction: Action = {
    id: "1",
    name: "An action",
    dueBy: new Date(2020, 0, 15),
  };

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("1/15/2020");

  expect(foundAction).toBeInTheDocument();
});
