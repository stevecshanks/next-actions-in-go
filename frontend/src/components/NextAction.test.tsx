import React from "react";
import MockDate from "mockdate";
import { render } from "@testing-library/react";
import { Action } from "../models/Action";
import { NextAction } from "./NextAction";

const NOW = new Date(2020, 0, 15, 10, 30, 0);
const ONE_SECOND = 1000;
const TWENTY_FOUR_HOURS = 24 * 60 * 60 * 1000;

beforeEach(() => MockDate.set(NOW));

afterEach(() => MockDate.reset());

test("displays action name", () => {
  const anAction: Action = new Action({ id: "1", name: "An action" });

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("An action");

  expect(foundAction).toBeInTheDocument();
});

test("displays action due date", () => {
  const later = new Date(NOW.getTime() + TWENTY_FOUR_HOURS);
  const anAction: Action = new Action({
    id: "1",
    name: "An action",
    dueBy: later,
  });

  const { getByText } = render(<NextAction action={anAction} />);

  const foundDate = getByText("1/16/2020");

  expect(foundDate).toBeInTheDocument();
  expect(foundDate.classList).toContain("badge-primary");
});

test("highlights overdue due dates", () => {
  const overdue = new Date(NOW.getTime() - ONE_SECOND);

  const overdueAction: Action = new Action({
    id: "1",
    name: "An overdue action",
    dueBy: overdue,
  });

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/15/2020");

  expect(foundDate.classList).toContain("badge-danger");
});

test("highlights due dates in the near future", () => {
  const soon = new Date(NOW.getTime() + TWENTY_FOUR_HOURS - ONE_SECOND);

  const overdueAction: Action = new Action({
    id: "1",
    name: "An action due soon",
    dueBy: soon,
  });

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/16/2020");

  expect(foundDate.classList).toContain("badge-warning");
});
