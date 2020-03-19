import React from "react";
import MockDate from "mockdate";
import { render } from "@testing-library/react";
import { Action } from "./NextActionsList";
import { NextAction } from "./NextAction";

const now = new Date(2020, 0, 15, 10, 30, 0);
const oneSecond = 1000;
const twentyFourHours = 24 * 60 * 60 * 1000;

beforeEach(() => MockDate.set(now));

afterEach(() => MockDate.reset());

test("displays action name", () => {
  const anAction: Action = { id: "1", name: "An action" };

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("An action");

  expect(foundAction).toBeInTheDocument();
});

test("displays action due date", () => {
  const later = new Date(now.getTime() + twentyFourHours);
  const anAction: Action = {
    id: "1",
    name: "An action",
    dueBy: later,
  };

  const { getByText } = render(<NextAction action={anAction} />);

  const foundDate = getByText("1/16/2020");

  expect(foundDate).toBeInTheDocument();
  expect(foundDate.classList).toContain("badge-primary");
});

test("highlights overdue due dates", () => {
  const overdue = new Date(now.getTime() - oneSecond);

  const overdueAction: Action = {
    id: "1",
    name: "An overdue action",
    dueBy: overdue,
  };

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/15/2020");

  expect(foundDate.classList).toContain("badge-danger");
});

test("highlights due dates in the near future", () => {
  const soon = new Date(now.getTime() + twentyFourHours - oneSecond);

  const overdueAction: Action = {
    id: "1",
    name: "An action due soon",
    dueBy: soon,
  };

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/16/2020");

  expect(foundDate.classList).toContain("badge-warning");
});
