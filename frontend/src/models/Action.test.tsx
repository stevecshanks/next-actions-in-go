import MockDate from "mockdate";
import { Action, ActionArgs } from "../models/Action";

const ACTION_DEFAULTS = {
  id: "anId123",
  name: "An Action",
  url: "https://trello.com/c/abc123/a-card",
  imageUrl: "https://example.com/image.jpg",
};

export const buildAction = (args?: Partial<ActionArgs>) =>
  new Action({
    ...ACTION_DEFAULTS,
    ...args,
  });

const NOW = new Date(2020, 0, 15, 10, 30, 0);
const ONE_SECOND = 1000;
const TWENTY_FOUR_HOURS = 24 * 60 * 60 * 1000;

beforeEach(() => MockDate.set(NOW));

afterEach(() => MockDate.reset());

test("actions with no due date are not overdue", () => {
  const action = buildAction();

  expect(action.isOverdue()).toBeFalsy();
});

test("actions due in the past are overdue", () => {
  const action = buildAction({
    dueBy: new Date(NOW.getTime() - ONE_SECOND),
  });

  expect(action.isOverdue()).toBeTruthy();
});

test("actions due now are not overdue", () => {
  const action = buildAction({
    dueBy: NOW,
  });

  expect(action.isOverdue()).toBeFalsy();
});

test("actions with no due date are not due soon", () => {
  const action = buildAction();

  expect(action.isDueSoon()).toBeFalsy();
});

test("actions due in less than a day are due soon", () => {
  const action = buildAction({
    dueBy: new Date(NOW.getTime() + TWENTY_FOUR_HOURS - ONE_SECOND),
  });

  expect(action.isDueSoon()).toBeTruthy();
});

test("actions due in a day are not due soon", () => {
  const action = buildAction({
    dueBy: new Date(NOW.getTime() + TWENTY_FOUR_HOURS),
  });

  expect(action.isDueSoon()).toBeFalsy();
});
