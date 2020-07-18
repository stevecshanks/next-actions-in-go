import React from "react";
import MockDate from "mockdate";
import { render } from "@testing-library/react";
import { NextAction } from "./NextAction";
import { buildAction } from "../models/Action.test";

const NOW = new Date(2020, 0, 15, 10, 30, 0);
const ONE_SECOND = 1000;
const TWENTY_FOUR_HOURS = 24 * 60 * 60 * 1000;

beforeEach(() => MockDate.set(NOW));

afterEach(() => MockDate.reset());

test("displays action name", () => {
  const anAction = buildAction({
    name: "An action",
  });

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("An action");

  expect(foundAction).toBeInTheDocument();
});

test("links the action to the specified URL", () => {
  const action = buildAction({
    name: "An action",
    url: "https://example.com/",
  });

  const { getByText } = render(<NextAction action={action} />);

  const link = getByText("An action") as HTMLAnchorElement;

  expect(link.href).toBe("https://example.com/");
});

test("displays action due date", () => {
  const later = new Date(NOW.getTime() + TWENTY_FOUR_HOURS);
  const anAction = buildAction({
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

  const overdueAction = buildAction({
    name: "An overdue action",
    dueBy: overdue,
  });

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/15/2020");

  expect(foundDate.classList).toContain("badge-danger");
});

test("highlights due dates in the near future", () => {
  const soon = new Date(NOW.getTime() + TWENTY_FOUR_HOURS - ONE_SECOND);

  const overdueAction = buildAction({
    name: "An action due soon",
    dueBy: soon,
  });

  const { getByText } = render(<NextAction action={overdueAction} />);

  const foundDate = getByText("1/16/2020");

  expect(foundDate.classList).toContain("badge-warning");
});

test("displays the relevant image for the action", () => {
  const actionWithImage = buildAction({
    name: "An action with image",
    imageUrl: "example.jpg",
  });

  const { getByText } = render(<NextAction action={actionWithImage} />);

  const foundAction = getByText("An action with image");

  expect(foundAction.querySelector(".action-image")).toHaveStyle(
    `background-image: url("example.jpg")`
  );
});

test("displays a placeholder if the action has no image", () => {
  const actionWithImage = buildAction({
    name: "An action with no image",
    imageUrl: null,
  });

  const { getByText } = render(<NextAction action={actionWithImage} />);

  const foundAction = getByText("An action with no image");

  expect(foundAction.querySelector(".action-image")).toHaveClass(
    "action-no-image"
  );
});

test("displays project name", () => {
  const anAction = buildAction({
    projectName: "A project with name",
  });

  const { getByText } = render(<NextAction action={anAction} />);

  const foundAction = getByText("A project with name");

  expect(foundAction).toBeInTheDocument();
});
